package edge

import (
	"context"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ml4wireless/specpipe/common"
)

func CaptureIQ(ctx context.Context, config *Config, publisher message.Publisher, logger common.EdgeLogrus) error {
	if config.Rtlsdr.Iq.Freq == "" {
		return ErrEmptyFreq
	}
	if config.Rtlsdr.Iq.SampleRate == "" {
		return ErrEmptySampleRate
	}
	cmd := exec.Command("rtl_sdr", "-s", config.Rtlsdr.Iq.SampleRate, "-f", config.Rtlsdr.Iq.Freq, "-b", "262144", "-")
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

	chunk := make([]byte, 16*16384)
	logger.Info("start IQ capturer")
	for {
		n, err := stdout.Read(chunk)
		if err != nil {
			if err == io.EOF {
				logger.Info("read EOF")
			}
			goto CLEANUP
		}
		for n < 262144 {
			bytesRead, err := stdout.Read(chunk[n:])
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
		payload := make([]byte, 16*16384)
		copy(payload, chunk)
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
