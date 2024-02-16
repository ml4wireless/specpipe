package edge

import (
	"context"
	"io"
	"os/exec"
	"time"

	"github.com/ml4wireless/specpipe/common"
	"github.com/ml4wireless/specpipe/proto/edge"
	"google.golang.org/grpc"
)

func ForwardIQ(ctx context.Context, config *ForwardConfig, logger common.EdgeLogrus, conn *grpc.ClientConn) error {
	client := edge.NewForwardServiceClient(conn)
	stream, err := client.Forward(ctx)
	if err != nil {
		return err
	}

	if config.Rtlsdr.Freq == "" {
		return ErrEmptyFreq
	}
	if config.Rtlsdr.SampleRate == "" {
		return ErrEmptySampleRate
	}
	cmd := exec.Command("rtl_sdr", "-s", config.Rtlsdr.SampleRate, "-f", config.Rtlsdr.Freq, "-b", "262144", "-")
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

		if _, err = stream.CloseAndRecv(); err != nil {
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

		forwardPayload := &edge.ForwardPayload{
			Payload: payload,
		}
		if err := stream.Send(forwardPayload); err != nil {
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
