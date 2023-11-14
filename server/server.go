package server

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ml4wireless/specpipe/common"
)

func NewHttpServer(specpipeServer *SpecpipeServer, logger common.ServerLogrus, port string) *http.Server {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CorsMiddleware())
	r.Use(LoggingMiddleware(logger))

	RegisterHandlersWithOptions(r, specpipeServer, GinServerOptions{
		BaseURL:      "/v0",
		ErrorHandler: errorHandler,
	})

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
	return s
}
