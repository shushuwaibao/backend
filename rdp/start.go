package rdp

import (
	"bytes"
	"context"
	"fmt"
	"gin-template/model"
	"gin-template/rdp/guac"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

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
		basearg := new(guac.ReqArgBaseInfo)
		err := c.BindQuery(basearg)
		if err != nil {
			c.JSON(202, err.Error())
			return
		}

		userid, exists := c.Get("id")
		if !exists {
			// 如果不存在，可能是因为用户未认证
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		arg, err := model.GetGuacdInfo(*basearg, userid.(int))

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
		pipeTunnel, err := guac.NewGuacamoleTunnel(&arg, uid)
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
				continue
			}
			if _, err = buf.Write(ins); err != nil {
				return err
			}
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
