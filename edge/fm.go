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

func CaptureAudio(ctx context.Context, config *FmConfig, publisher message.Publisher, logger common.EdgeLogrus) error {
	if config.Rtlsdr.Freq == "" {
		return ErrEmptyFreq
	}
	if config.Rtlsdr.SampleRate == "" {
		return ErrEmptySampleRate
	}
	if config.Rtlsdr.ResampleRate == "" {
		return ErrEmptyReampleRate
	}
	cmd := exec.Command("rtl_fm", "-M", "fm", "-s", config.Rtlsdr.SampleRate, "-o", "4", "-A", "fast", "-r", config.Rtlsdr.ResampleRate, "-l", "0", "-E", "deemp", "-f", config.Rtlsdr.Freq)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "RTLSDR_RPC_IS_ENABLED=1", "RTLSDR_RPC_SERV_ADDR="+config.Rtlsdr.Rpc.ServerAddr, "RTLSDR_RPC_SERV_PORT="+config.Rtlsdr.Rpc.ServerPort)
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
