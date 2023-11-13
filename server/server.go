package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ml4wireless/specpipe/common"
)

type SpecpipeServer struct {
	store *Store
}

func NewSpecpipeServer(store *Store) *SpecpipeServer {
	return &SpecpipeServer{store}
}

func (s *SpecpipeServer) GetFmDevices(c *gin.Context) {
	fmDevices, err := s.store.GetFmDevices(c.Request.Context())
	if err != nil {
		errorHandler(c, fmt.Errorf("get fm devices error: %w", err), http.StatusInternalServerError)
		return
	}
	fmDevicesPresenter := []FmDevice{}
	for _, fmDevice := range fmDevices {
		fmDevicesPresenter = append(fmDevicesPresenter, FmDevice{
			Freq:       fmDevice.Freq,
			Latitude:   fmDevice.Latitude,
			Longitude:  fmDevice.Longitude,
			Name:       fmDevice.Name,
			SampleRate: fmDevice.SampleRate,
		})
	}
	c.JSON(http.StatusOK, &FmDevicesResponse{
		Devices: fmDevicesPresenter,
	})
}
func (s *SpecpipeServer) GetFmDevicesDevicename(c *gin.Context, deviceName string) {
	exist, fmDevice, err := s.store.GetFmDevice(c.Request.Context(), deviceName)
	if err != nil {
		errorHandler(c, fmt.Errorf("get fm device error: %w", err), http.StatusInternalServerError)
		return
	}
	if !exist {
		errorHandler(c, fmt.Errorf("fm device %s not found", deviceName), http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, &FmDeviceResponse{
		Device: FmDevice{
			Freq:       fmDevice.Freq,
			Latitude:   fmDevice.Latitude,
			Longitude:  fmDevice.Longitude,
			Name:       fmDevice.Name,
			SampleRate: fmDevice.SampleRate,
		},
	})
}
func (s *SpecpipeServer) PutFmDevicesDevicename(c *gin.Context, deviceName string) {
	var updateFmDeviceRequest UpdateFmDeviceRequest
	if err := c.ShouldBindJSON(&updateFmDeviceRequest); err != nil {
		errorHandler(c, fmt.Errorf("parse request error: %w", err), http.StatusBadRequest)
		return
	}
	exist, fmDevice, err := s.store.GetFmDevice(c.Request.Context(), deviceName)
	if err != nil {
		errorHandler(c, fmt.Errorf("get fm device error: %w", err), http.StatusInternalServerError)
		return
	}
	if !exist {
		errorHandler(c, fmt.Errorf("fm device %s not found", deviceName), http.StatusNotFound)
		return
	}

	fmDevice.Freq = updateFmDeviceRequest.Freq
	if updateFmDeviceRequest.SampleRate != nil {
		fmDevice.SampleRate = *updateFmDeviceRequest.SampleRate
	}
	if err = s.store.UpdateFmDevice(c.Request.Context(), fmDevice); err != nil {
		errorHandler(c, fmt.Errorf("update fm device error: %w", err), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &FmDeviceResponse{
		Device: FmDevice{
			Freq:       fmDevice.Freq,
			Latitude:   fmDevice.Latitude,
			Longitude:  fmDevice.Longitude,
			Name:       fmDevice.Name,
			SampleRate: fmDevice.SampleRate,
		},
	})
}

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

func errorHandler(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, ErrorResponse{
		Title:  "specpipe server error",
		Detail: err.Error(),
	})
}
