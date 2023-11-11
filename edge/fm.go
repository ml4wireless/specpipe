package edge

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ml4wireless/specpipe/common"
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
		if err = publisher.Close(); err != nil {
			return err
		}
		return nil
	}

	audio := make([]byte, 2*8192)
	logger.Info("start caputuring FM audio")
	for {
		n, err := stdout.Read(audio)
		for n < 16384 {
			bytesRead, err := stdout.Read(audio[n:])
			if err != nil {
				if err == io.EOF {
					break
				}
				logger.Error(err)
			}
			n += bytesRead
			if n < 16384 {
				time.Sleep(10 * time.Millisecond)
			}
			select {
			case <-ctx.Done():
				if err = cleanup(); err != nil {
					return err
				}
				return nil
			default:
			}
		}
		if err == io.EOF {
			logger.Info("read EOF")
			break
		}
		if err != nil {
			return err
		}
		payload := make([]byte, 2*8192)
		copy(payload, audio)
		msg := message.NewMessage(strconv.FormatInt(time.Now().UTC().Unix(), 10), payload)
		if err := publisher.Publish(config.Pub.Subject, msg); err != nil {
			logger.Error(err)
		}
		select {
		case <-ctx.Done():
			if err = cleanup(); err != nil {
				return err
			}
			return nil
		default:
		}
	}
	return nil
}
