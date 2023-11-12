package edge

import (
	"github.com/ml4wireless/specpipe/common"
	"github.com/nats-io/nats.go"
)

func RegisterDevice(conn *nats.Conn, sdrType common.SDRType, deviceName string, deviceInfo []byte) (*nats.Subscription, error) {
	// _, err := conn.Request(common.ClusterSubject(sdrType, deviceName, common.RegisterCmd), deviceInfo, time.Second)
	// if err != nil {
	// 	return nil, err
	// }
	// healthcheck routine
	return conn.Subscribe(common.ClusterSubject(sdrType, deviceName, common.HealthCheckCmd), func(msg *nats.Msg) {
		msg.Respond([]byte("ok"))
	})
}

func DeregisterDevice(conn *nats.Conn, sdrType common.SDRType, deviceName string, sub *nats.Subscription) error {
	// _, err := conn.Request(common.ClusterSubject(sdrType, deviceName, common.DeregisterCmd), nil, time.Second)
	// if err != nil {
	// 	return err
	// }
	return sub.Unsubscribe()
}
