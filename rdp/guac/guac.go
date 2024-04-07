package guac

import (
	"net"

	"github.com/sirupsen/logrus"
)

func NewGuacamoleTunnel(arg *ReqArg, uuid string) (s *SimpleTunnel, err error) {
	config := NewGuacamoleConfiguration()
	config.ConnectionID = uuid
	config.Protocol = arg.AssetProtocol
	config.OptimalScreenHeight = arg.ScreenHeight
	config.OptimalScreenWidth = arg.ScreenWidth
	config.OptimalResolution = arg.ScreenDpi
	config.AudioMimetypes = []string{"audio/L16", "rate=44100", "channels=2"}
	config.Parameters = map[string]string{
		"scheme":      arg.AssetProtocol,
		"hostname":    arg.AssetHost,
		"port":        arg.AssetPort,
		"ignore-cert": "true",
		"security":    "",
		"username":    arg.AssetUser,
		"password":    arg.AssetPassword,
	}
	addr, err := net.ResolveTCPAddr("tcp", arg.GuacadAddr)
	if err != nil {
		logrus.Errorln("error while connecting to guacd", err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logrus.Errorln("error while connecting to guacd", err)
		return nil, err
	}
	stream := NewStream(conn, SocketTimeout)
	// init rdp guacd asset
	err = stream.Handshake(config)
	if err != nil {
		return nil, err
	}
	tunnel := NewSimpleTunnel(stream)
	return tunnel, nil
}
