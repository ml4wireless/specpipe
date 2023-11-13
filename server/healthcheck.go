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
			checkFmDevices(ctx, t, conn, store, heartbeatTimoutSec, logger)
		}
	}
}

func checkFmDevices(ctx context.Context, t time.Time, conn *nats.Conn, store *Store, heartbeatTimoutSec int64, logger common.ServerLogrus) {
	devices, err := store.GetDevices(ctx, common.FM)
	if err != nil {
		return
	}
	for _, device := range devices {
		fmDevice, ok := device.(*common.FMDevice)
		if !ok {
			logger.Error(fmt.Errorf("casting fm device type error: %w", err))
			continue
		}
		_, err := conn.Request(common.ClusterSubject(common.FM, fmDevice.Name, common.HealthCheckCmd), nil, time.Duration(heartbeatTimoutSec)*time.Second)
		if err != nil {
			logger.Error(fmt.Errorf("send heartbeat request to device %s error: %w", fmDevice.Name, err))
			if err = store.DeleteDevice(ctx, common.FM, fmDevice.Name); err != nil {
				logger.Error(fmt.Errorf("remove device %s error: %w", fmDevice.Name, err))
				continue
			}
			logger.Infof("unhealthy device removed: type=%s name=%s at %d-%d-%d %d:%d:%d", common.FM, fmDevice.Name, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
			continue
		}
		logger.Infof("device healthy: type=%s name=%s at %d-%d-%d %d:%d:%d", common.FM, fmDevice.Name, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
}
