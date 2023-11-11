package common

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ServerLogrus struct {
	*log.Entry
}
type EdgeLogrus struct {
	*log.Entry
}

func NewServerLogrus(level string) (ServerLogrus, error) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Writer(os.Stderr)

	if err := initLogging(level); err != nil {
		return ServerLogrus{}, err
	}

	return ServerLogrus{log.WithField("agent", "sp-server")}, nil
}

func NewEdgeLogrus(level string) (EdgeLogrus, error) {
	if err := initLogging(level); err != nil {
		return EdgeLogrus{}, err
	}

	return EdgeLogrus{log.WithField("agent", "sp-edge")}, nil
}

func initLogging(level string) error {
	logrusLevel, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetOutput(os.Stderr)
	log.SetLevel(logrusLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	return nil
}
