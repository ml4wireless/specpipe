package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: iqstreamer [iq_data_file_path]")
		return
	}
	iqDataFilePath := os.Args[1]

	chunkSize := 16 * 16384

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	natsURL := getenv("NATS_URL", "nats://127.0.0.1:4222")
	deviceName := getenv("DEVICE", "")
	if deviceName == "" {
		fmt.Println("device name cannot be empty")
		return
	}
	natsSubject := fmt.Sprintf("specpipe-iq.data.iq.%s", deviceName)

	marshaler := &nats.NATSMarshaler{}
	logger := watermill.NewStdLogger(false, false)
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}
	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         natsURL,
			NatsOptions: options,
			Marshaler:   marshaler,
			JetStream: nats.JetStreamConfig{
				Disabled:       false,
				AutoProvision:  false,
				PublishOptions: nil,
				TrackMsgId:     false,
				AckAsync:       false,
			},
		},
		logger,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open(iqDataFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Printf("streaming raw IQ data to NATS subject %s\n", natsSubject)
	chunk := make([]byte, chunkSize)
	for {
		n, err := f.Read(chunk)
		if err != nil {
			if err == io.EOF {
				f.Seek(0, 0)
				continue
			}
			fmt.Println(err)
			return
		}

		for n < chunkSize {
			bytesRead, err := f.Read(chunk[n:])
			if err != nil {
				if err == io.EOF {
					f.Seek(0, 0)
					continue
				}
				fmt.Println(err)
				return
			}
			n += bytesRead
			if n < chunkSize {
				time.Sleep(10 * time.Millisecond)
			}
			select {
			case <-sig:
				if err = publisher.Close(); err != nil {
					fmt.Println(err)
				}
				return
			default:
			}
		}

		payload := make([]byte, chunkSize)
		copy(payload, chunk)
		msg := message.NewMessage(watermill.NewShortUUID(), payload)
		msg.Metadata.Set("ts", strconv.FormatInt(time.Now().UTC().Unix(), 10))
		if err := publisher.Publish(natsSubject, msg); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(20 * time.Millisecond)

		select {
		case <-sig:
			if err = publisher.Close(); err != nil {
				fmt.Println(err)
			}
			return
		default:
		}
	}
}
