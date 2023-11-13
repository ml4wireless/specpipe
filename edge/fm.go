package edge

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ml4wireless/specpipe/common"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	ErrEmptyFreq        = errors.New("frequency cannot be empty")
	ErrEmptySampleRate  = errors.New("samping rate cannot be empty")
	ErrEmptyReampleRate = errors.New("resamping rate cannot be empty")
)

func CaptureAudio(ctx context.Context, config *Config, publisher message.Publisher, logger common.EdgeLogrus) error {
	if config.Rtlsdr.Fm.Freq == "" {
		return ErrEmptyFreq
	}
	if config.Rtlsdr.Fm.SampleRate == "" {
		return ErrEmptySampleRate
	}
	if config.Rtlsdr.Fm.ResampleRate == "" {
		return ErrEmptyReampleRate
	}
	cmd := exec.Command("rtl_fm", "-M", "fm", "-s", config.Rtlsdr.Fm.SampleRate, "-o", "4", "-A", "fast", "-r", config.Rtlsdr.Fm.ResampleRate, "-l", "0", "-E", "deemp", "-f", config.Rtlsdr.Fm.Freq)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "RTLSDR_RPC_IS_ENABLED=1", "RTLSDR_RPC_SERV_ADDR="+config.Rtlsdr.RpcServerAddr, "RTLSDR_RPC_SERV_PORT="+config.Rtlsdr.RpcServerPort)
	if err := cmd.Start(); err != nil {
		return err
	}

	cleanup := func() error {
		if err = cmd.Process.Kill(); err != nil {
			return err
		}
		return nil
	}

	audio := make([]byte, 2*8192)
	logger.Info("start FM audio capturer")
	for {
		n, err := stdout.Read(audio)
		if err != nil {
			if err == io.EOF {
				logger.Info("read EOF")
			}
			goto CLEANUP
		}
		for n < 16384 {
			bytesRead, err := stdout.Read(audio[n:])
			if err != nil {
				if err == io.EOF {
					logger.Info("read EOF")
				}
				goto CLEANUP
			}
			n += bytesRead
			if n < 16384 {
				time.Sleep(10 * time.Millisecond)
			}
			select {
			case <-ctx.Done():
				goto CLEANUP
			default:
			}
		}
		payload := make([]byte, 2*8192)
		copy(payload, audio)
		msg := message.NewMessage(watermill.NewShortUUID(), payload)
		msg.Metadata.Set(common.TimestampHeader, strconv.FormatInt(time.Now().UTC().Unix(), 10))
		if err := publisher.Publish(config.Nats.Subject, msg); err != nil {
			logger.Error(err)
		}
		select {
		case <-ctx.Done():
			goto CLEANUP
		default:
		}
	}
CLEANUP:
	if err = cleanup(); err != nil {
		return err
	}
	return nil
}

func WatchFmDeviceConfig(ctx context.Context, conn *nats.Conn, kv jetstream.KeyValue, deviceName string, logger common.EdgeLogrus) (*nats.Subscription, chan common.FMDevice, error) {
	fmDeviceConfigChan := make(chan common.FMDevice)
	watchSub, err := conn.Subscribe(common.ClusterSubject(common.FM, deviceName, common.WatchConfigCmd), func(msg *nats.Msg) {
		entry, err := kv.Get(ctx, common.KVStoreKey(common.FM, deviceName))
		if err != nil {
			logger.Error(err)
		}

		var device common.FMDevice
		if err = json.Unmarshal(entry.Value(), &device); err != nil {
			logger.Error(err)
		}

		fmDeviceConfigChan <- device

		msg.Respond([]byte(common.OkMsg))
	})

	if err != nil {
		return nil, nil, err
	}

	return watchSub, fmDeviceConfigChan, nil
}
