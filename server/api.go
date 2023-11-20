package server

import (
	"errors"
	"fmt"
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
	devices, err := s.store.GetDevices(c.Request.Context(), common.FM)
	if err != nil {
		errorHandler(c, fmt.Errorf("get fm devices error: %w", err), http.StatusInternalServerError)
		return
	}
	fmDevicesPresenter := []FmDevice{}
	for _, device := range devices {
		fmDevice, ok := device.(*common.FMDevice)
		if !ok {
			errorHandler(c, errors.New("casting fm device type error"), http.StatusInternalServerError)
			return
		}
		fmDevicesPresenter = append(fmDevicesPresenter, FmDevice{
			Freq:         fmDevice.Freq,
			Latitude:     fmDevice.Latitude,
			Longitude:    fmDevice.Longitude,
			Name:         fmDevice.Name,
			SampleRate:   fmDevice.SampleRate,
			ResampleRate: fmDevice.ResampleRate,
		})
	}
	c.JSON(http.StatusOK, &FmDevicesResponse{
		Devices: fmDevicesPresenter,
	})
}
func (s *SpecpipeServer) GetFmDevicesDevicename(c *gin.Context, deviceName string) {
	exist, device, err := s.store.GetDevice(c.Request.Context(), common.FM, deviceName)
	if err != nil {
		errorHandler(c, fmt.Errorf("get fm device error: %w", err), http.StatusInternalServerError)
		return
	}
	if !exist {
		errorHandler(c, fmt.Errorf("fm device %s not found", deviceName), http.StatusNotFound)
		return
	}
	fmDevice, ok := device.(*common.FMDevice)
	if !ok {
		errorHandler(c, fmt.Errorf("casting fm device type error: %w", err), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &FmDeviceResponse{
		Device: FmDevice{
			Freq:         fmDevice.Freq,
			Latitude:     fmDevice.Latitude,
			Longitude:    fmDevice.Longitude,
			Name:         fmDevice.Name,
			SampleRate:   fmDevice.SampleRate,
			ResampleRate: fmDevice.ResampleRate,
		},
	})
}
func (s *SpecpipeServer) PutFmDevicesDevicename(c *gin.Context, deviceName string) {
	var updateFmDeviceRequest UpdateFmDeviceRequest
	if err := c.ShouldBindJSON(&updateFmDeviceRequest); err != nil {
		errorHandler(c, fmt.Errorf("parse request error: %w", err), http.StatusBadRequest)
		return
	}
	exist, device, err := s.store.GetDevice(c.Request.Context(), common.FM, deviceName)
	if err != nil {
		errorHandler(c, fmt.Errorf("get fm device error: %w", err), http.StatusInternalServerError)
		return
	}
	if !exist {
		errorHandler(c, fmt.Errorf("fm device %s not found", deviceName), http.StatusNotFound)
		return
	}
	fmDevice, ok := device.(*common.FMDevice)
	if !ok {
		errorHandler(c, fmt.Errorf("casting fm device type error: %w", err), http.StatusInternalServerError)
		return
	}

	fmDevice.Freq = updateFmDeviceRequest.Freq
	if updateFmDeviceRequest.SampleRate != nil {
		fmDevice.SampleRate = *updateFmDeviceRequest.SampleRate
	}
	if updateFmDeviceRequest.ResampleRate != nil {
		fmDevice.ResampleRate = *updateFmDeviceRequest.ResampleRate
	}
	if err = s.store.UpdateDevice(c.Request.Context(), common.FM, deviceName, fmDevice); err != nil {
		errorHandler(c, fmt.Errorf("update fm device error: %w", err), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &FmDeviceResponse{
		Device: FmDevice{
			Freq:         fmDevice.Freq,
			Latitude:     fmDevice.Latitude,
			Longitude:    fmDevice.Longitude,
			Name:         fmDevice.Name,
			SampleRate:   fmDevice.SampleRate,
			ResampleRate: fmDevice.ResampleRate,
		},
	})
}

func (s *SpecpipeServer) GetIqDevices(c *gin.Context) {
	devices, err := s.store.GetDevices(c.Request.Context(), common.IQ)
	if err != nil {
		errorHandler(c, fmt.Errorf("get iq devices error: %w", err), http.StatusInternalServerError)
		return
	}
	iqDevicesPresenter := []IqDevice{}
	for _, device := range devices {
		iqDevice, ok := device.(*common.IQDevice)
		if !ok {
			errorHandler(c, errors.New("casting iq device type error"), http.StatusInternalServerError)
			return
		}
		iqDevicesPresenter = append(iqDevicesPresenter, IqDevice{
			Freq:       iqDevice.Freq,
			Latitude:   iqDevice.Latitude,
			Longitude:  iqDevice.Longitude,
			Name:       iqDevice.Name,
			SampleRate: iqDevice.SampleRate,
		})
	}
	c.JSON(http.StatusOK, &IqDevicesResponse{
		Devices: iqDevicesPresenter,
	})
}

func (s *SpecpipeServer) GetIqDevicesDevicename(c *gin.Context, deviceName string) {
	exist, device, err := s.store.GetDevice(c.Request.Context(), common.IQ, deviceName)
	if err != nil {
		errorHandler(c, fmt.Errorf("get iq device error: %w", err), http.StatusInternalServerError)
		return
	}
	if !exist {
		errorHandler(c, fmt.Errorf("iq device %s not found", deviceName), http.StatusNotFound)
		return
	}
	iqDevice, ok := device.(*common.IQDevice)
	if !ok {
		errorHandler(c, fmt.Errorf("casting iq device type error: %w", err), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &IqDeviceResponse{
		Device: IqDevice{
			Freq:       iqDevice.Freq,
			Latitude:   iqDevice.Latitude,
			Longitude:  iqDevice.Longitude,
			Name:       iqDevice.Name,
			SampleRate: iqDevice.SampleRate,
		},
	})
}

func (s *SpecpipeServer) PutIqDevicesDevicename(c *gin.Context, deviceName string) {
	var updateIqDeviceRequest UpdateIqDeviceRequest
	if err := c.ShouldBindJSON(&updateIqDeviceRequest); err != nil {
		errorHandler(c, fmt.Errorf("parse request error: %w", err), http.StatusBadRequest)
		return
	}
	exist, device, err := s.store.GetDevice(c.Request.Context(), common.IQ, deviceName)
	if err != nil {
		errorHandler(c, fmt.Errorf("get iq device error: %w", err), http.StatusInternalServerError)
		return
	}
	if !exist {
		errorHandler(c, fmt.Errorf("iq device %s not found", deviceName), http.StatusNotFound)
		return
	}
	iqDevice, ok := device.(*common.IQDevice)
	if !ok {
		errorHandler(c, fmt.Errorf("casting iq device type error: %w", err), http.StatusInternalServerError)
		return
	}

	iqDevice.Freq = updateIqDeviceRequest.Freq
	if updateIqDeviceRequest.SampleRate != nil {
		iqDevice.SampleRate = *updateIqDeviceRequest.SampleRate
	}
	if err = s.store.UpdateDevice(c.Request.Context(), common.IQ, deviceName, iqDevice); err != nil {
		errorHandler(c, fmt.Errorf("update iq device error: %w", err), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &IqDeviceResponse{
		Device: IqDevice{
			Freq:       iqDevice.Freq,
			Latitude:   iqDevice.Latitude,
			Longitude:  iqDevice.Longitude,
			Name:       iqDevice.Name,
			SampleRate: iqDevice.SampleRate,
		},
	})
}

func errorHandler(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, ErrorResponse{
		Title:  "specpipe server error",
		Detail: err.Error(),
	})
}
