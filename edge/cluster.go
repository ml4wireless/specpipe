package edge

import (
	"context"
	"encoding/json"

	"github.com/ml4wireless/specpipe/common"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func RegisterDevice(ctx context.Context, conn *nats.Conn, kv jetstream.KeyValue, sdrType common.SDRType, deviceName string, deviceInfo []byte) (*nats.Subscription, error) {
	_, err := kv.Put(ctx, common.KVStoreKey(sdrType, deviceName), deviceInfo)
	if err != nil {
		return nil, err
	}
	// healthcheck routine
	return conn.Subscribe(common.ClusterSubject(sdrType, deviceName, common.HealthCheckCmd), func(msg *nats.Msg) {
		msg.Respond([]byte(common.OkMsg))
	})
}

func DeregisterDevice(ctx context.Context, kv jetstream.KeyValue, sdrType common.SDRType, deviceName string) error {
	return kv.Delete(ctx, common.KVStoreKey(sdrType, deviceName))
}

func WatchDevice(ctx context.Context, conn *nats.Conn, kv jetstream.KeyValue, sdrType common.SDRType, deviceName string, logger common.EdgeLogrus) (*nats.Subscription, chan common.Device, error) {
	deviceChan := make(chan common.Device)
	watchSub, err := conn.Subscribe(common.ClusterSubject(sdrType, deviceName, common.WatchConfigCmd), func(msg *nats.Msg) {
		entry, err := kv.Get(ctx, common.KVStoreKey(sdrType, deviceName))
		if err != nil {
			logger.Error(err)
			return
		}
		switch sdrType {
		case common.FM:
			var device common.FMDevice
			if err = json.Unmarshal(entry.Value(), &device); err != nil {
				logger.Error(err)
			}
			deviceChan <- &device
		case common.IQ:
			var device common.IQDevice
			if err = json.Unmarshal(entry.Value(), &device); err != nil {
				logger.Error(err)
			}
			deviceChan <- &device
		}
		msg.Respond([]byte(common.OkMsg))
	})

	if err != nil {
		return nil, nil, err
	}

	return watchSub, deviceChan, nil
}
