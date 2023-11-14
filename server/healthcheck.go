package server

import (
	"context"
	"fmt"
	"time"

	"github.com/ml4wireless/specpipe/common"
	"github.com/nats-io/nats.go"
)

func Healthcheck(ctx context.Context, conn *nats.Conn, store *Store, heartbeatTimoutSec int64, logger common.ServerLogrus) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			checkDevices(ctx, t, conn, store, common.FM, heartbeatTimoutSec, logger)
			checkDevices(ctx, t, conn, store, common.IQ, heartbeatTimoutSec, logger)
		}
	}
}

func checkDevices(ctx context.Context, t time.Time, conn *nats.Conn, store *Store, sdrType common.SDRType, heartbeatTimoutSec int64, logger common.ServerLogrus) {
	devices, err := store.GetDevices(ctx, sdrType)
	if err != nil {
		logger.Error(fmt.Errorf("get devices error: %w", err))
		return
	}
OUTER:
	for _, device := range devices {
		var deviceName string
		switch sdrType {
		case common.FM:
			fmDevice, ok := device.(*common.FMDevice)
			if !ok {
				logger.Error(fmt.Errorf("casting fm device type error: %w", err))
				goto OUTER
			}
			deviceName = fmDevice.Name
		case common.IQ:
			iqDevice, ok := device.(*common.IQDevice)
			if !ok {
				logger.Error(fmt.Errorf("casting iq device type error: %w", err))
				goto OUTER
			}
			deviceName = iqDevice.Name
		default:
			goto OUTER
		}
		_, err := conn.Request(common.ClusterSubject(sdrType, deviceName, common.HealthCheckCmd), nil, time.Duration(heartbeatTimoutSec)*time.Second)
		if err != nil {
			logger.Error(fmt.Errorf("send heartbeat request to device %s error: %w", deviceName, err))
			if err = store.DeleteDevice(ctx, sdrType, deviceName); err != nil {
				logger.Error(fmt.Errorf("remove device %s error: %w", deviceName, err))
				continue
			}
			logger.Infof("unhealthy device removed: type=%s name=%s at %d-%d-%d %d:%d:%d", sdrType, deviceName, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
			continue
		}
		logger.Infof("device healthy: type=%s name=%s at %d-%d-%d %d:%d:%d", sdrType, deviceName, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
}
