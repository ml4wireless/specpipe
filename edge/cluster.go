package edge

import (
	"context"

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
