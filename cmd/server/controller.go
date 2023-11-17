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

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "specpipe controller",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := server.NewControllerConfig()
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
		svr := server.NewSpecpipeServer(store)
		httpSvr := server.NewHttpServer(svr, logger, config.Http.ServerPort)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Infof("start specpipe http server at :%s", config.Http.ServerPort)
			httpSvr.ListenAndServe()
			log.Info("specpipe http server stopped")
		}()
		<-sig
		httpSvr.Close()

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)

	controllerCmd.Flags().StringP("http-server-port", "", "", "http server port")
	viper.BindPFlag("controller.http.serverPort", controllerCmd.Flags().Lookup("http-server-port"))

	controllerCmd.Flags().StringP("nats-url", "", "", "NATS URL")
	viper.BindPFlag("controller.nats.url", controllerCmd.Flags().Lookup("nats-url"))

	controllerCmd.Flags().StringP("log-level", "", "", "log level")
	viper.BindPFlag("controller.logging.level", controllerCmd.Flags().Lookup("log-level"))
}
