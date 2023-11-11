package common

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
)

func NewNatsPublisher(natsUrl string) (message.Publisher, error) {
	marshaler := &nats.NATSMarshaler{}
	logger := watermill.NewStdLogger(false, false)
	return nats.NewPublisher(
		nats.PublisherConfig{
			URL: natsUrl,
			NatsOptions: []nc.Option{
				nc.RetryOnFailedConnect(true),
				nc.Timeout(3 * time.Second),
				nc.ReconnectWait(1 * time.Second),
			},
			Marshaler: marshaler,
			JetStream: nats.JetStreamConfig{
				Disabled:       false,
				AutoProvision:  false,
				PublishOptions: nil,
				TrackMsgId:     false,
			},
		},
		logger,
	)
}
