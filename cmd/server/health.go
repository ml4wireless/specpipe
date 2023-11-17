package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ml4wireless/specpipe/common"
	"github.com/ml4wireless/specpipe/server"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "specpipe health checker",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := server.NewHealthConfig()
		if err != nil {
			log.Fatal(err)
		}
		logger, err := common.NewServerLogrus(config.Logging.Level)
		if err != nil {
			log.Fatal(err)
		}
		clusterConn, err := common.NewNATSConn(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer clusterConn.Close()
		js, err := jetstream.New(clusterConn)
		if err != nil {
			log.Fatal(err)
		}
		kvCtx, kvCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer kvCancel()
		kv, err := js.KeyValue(kvCtx, common.KVStoreBucket)
		if err != nil {
			log.Fatal(err)
		}
		store := server.NewStore(clusterConn, kv)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			log.Info("start specpipe heartbeat routine")
			server.Healthcheck(ctx, clusterConn, store, config.TimeoutSec, logger)
			log.Info("specpipe heartbeat routine stopped")
		}()
		<-sig
		cancel()

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)

	healthCmd.Flags().StringP("heartbeat-timeout-sec", "", "", "heartbeat timeout seconds for device monitoring")
	viper.BindPFlag("health.timeoutSec", healthCmd.Flags().Lookup("heartbeat-timeout-sec"))

	healthCmd.Flags().StringP("nats-url", "", "", "NATS URL")
	viper.BindPFlag("health.nats.url", healthCmd.Flags().Lookup("nats-url"))

	healthCmd.Flags().StringP("log-level", "", "", "log level")
	viper.BindPFlag("health.logging.level", healthCmd.Flags().Lookup("log-level"))
}
