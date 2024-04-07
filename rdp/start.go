package rdp

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"gin-template/rdp/guac"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	fsBase   = "frontend/dist"
	mainPage = "index.html"
)

//go:embed frontend/dist/*
var fs embed.FS

func StartRDPService(r *gin.Engine) {
	logrus.SetReportCaller(true)
	r.Use(ParseUrl("/"))
	r.GET("/ws", MakeConnection())
	r.Run(":9528")
}

func ParseUrl(urlPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := strings.TrimSpace(c.Request.URL.Path)
		if urlPath == urlPrefix {
			urlPath = path.Join(urlPrefix, mainPage)
		}
		urlPath = filepath.Join(fsBase, urlPath)

		if strings.HasSuffix(urlPath, ".html") {
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "text/html; charset=utf-8")
		}
		if strings.HasSuffix(urlPath, ".js") {
			c.Header("Content-Type", "text/javascript; charset=utf-8")
		}
		if strings.HasSuffix(urlPath, ".css") {
			c.Header("Content-Type", "text/css; charset=utf-8")
		}

		f, err := fs.Open(urlPath)
		if err != nil {
			return
		}

		fi, err := f.Stat()
		if err != nil || !fi.IsDir() {
			bs, err := fs.ReadFile(urlPath)
			if err != nil {
				logrus.WithError(err).Error("embed fs")
				return
			}
			c.Status(200)
			c.Writer.Write(bs)
			c.Abort()
		}
	}
}

func MakeConnection() gin.HandlerFunc {
	websocketReadBufferSize := guac.MaxGuacMessage
	websocketWriteBufferSize := guac.MaxGuacMessage * 2
	upgrade := websocket.Upgrader{
		ReadBufferSize:  websocketReadBufferSize,
		WriteBufferSize: websocketWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return func(c *gin.Context) {
		logrus.Println("1. Parse argument")
		arg := new(guac.ReqArg)
		err := c.BindQuery(arg)
		if err != nil {
			c.JSON(202, err.Error())
			return
		}

		logrus.Println("2. Upgrade websocket")
		protocol := c.Request.Header.Get("Sec-Websocket-Protocol")
		ws, err := upgrade.Upgrade(c.Writer, c.Request, http.Header{
			"Sec-Websocket-Protocol": {protocol},
		})
		if err != nil {
			logrus.WithError(err).Error("Upgrade websocket failed")
			return
		}
		defer func() {
			if err = ws.Close(); err != nil {
				logrus.Traceln("Error closing websocket", err)
			}
		}()

		logrus.Println("3. Connect Asset")
		uid := ""
		pipeTunnel, err := guac.NewGuacamoleTunnel(arg, uid)
		if err != nil {
			logrus.Error("Failed to upgrade websocket", err)
			return
		}
		defer func() {
			if err = pipeTunnel.Close(); err != nil {
				logrus.Traceln("Error closing pipeTunnel", err)
			}
		}()

		logrus.Println("4. Handle guacd-tunnel IO")
		CreateIOTunnel(ws, pipeTunnel)
		logrus.Info("websocket session end")
	}
}

func CreateIOTunnel(ws *websocket.Conn, tunnl *guac.SimpleTunnel) {
	writer := tunnl.AcquireWriter()
	reader := tunnl.AcquireReader()

	defer tunnl.ReleaseWriter()
	defer tunnl.ReleaseReader()

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		buf := bytes.NewBuffer(make([]byte, 0, guac.MaxGuacMessage*2))

		for {
			ins, err := reader.ReadSome()
			if err != nil {
				return err
			}
			if bytes.HasPrefix(ins, guac.InternalOpcodeIns) {
				// messages starting with the InternalDataOpcode are never sent to the websocket
				continue
			}
			if _, err = buf.Write(ins); err != nil {
				return err
			}
			// if the buffer has more data in it or we've reached the max buffer size, send the data and reset
			if !reader.Available() || buf.Len() >= guac.MaxGuacMessage {
				if err = ws.WriteMessage(1, buf.Bytes()); err != nil {
					if err == websocket.ErrCloseSent {
						return fmt.Errorf("websocket:%v", err)
					}
					logrus.Traceln("Failed sending message to ws", err)
					return err
				}
				buf.Reset()
			}
		}
	})

	eg.Go(func() error {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				logrus.Traceln("Error reading message from ws", err)
				return err
			}
			if bytes.HasPrefix(data, guac.InternalOpcodeIns) {
				// messages starting with the InternalDataOpcode are never sent to guacd
				continue
			}
			if _, err = writer.Write(data); err != nil {
				logrus.Traceln("Failed writing to guacd", err)
				return err
			}
		}

	})

	if err := eg.Wait(); err != nil {
		logrus.WithError(err).Error("session-err")
	}
}
