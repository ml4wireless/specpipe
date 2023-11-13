package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"

	"github.com/ml4wireless/specpipe/common"
	"github.com/ml4wireless/specpipe/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string
	cfgFile string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Long:  `Print the current version of specpipe server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

var rootCmd = &cobra.Command{
	Use:   "sp-server",
	Short: "specpipe server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := server.NewConfig()
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
		kvConn, err := common.NewNATSConn(config.Nats.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer kvConn.Close()
		js, err := jetstream.New(kvConn)
		if err != nil {
			log.Fatal(err)
		}
		kvCtx, kvCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer kvCancel()
		kv, err := js.KeyValue(kvCtx, common.KVStoreBucket)
		if err != nil {
			log.Fatal(err)
		}
		store := server.NewStore(kv)
		svr := server.NewSpecpipeServer(store)
		httpSvr := server.NewHttpServer(svr, logger, config.Http.ServerPort)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			log.Infof("start specpipe http server at :%s", config.Http.ServerPort)
			httpSvr.ListenAndServe()
			log.Info("specpipe http server stopped")
		}()

		go func() {
			defer wg.Done()
			log.Info("start specpipe heartbeat routine")
			server.Healthcheck(ctx, clusterConn, store, config.Heartbeat.TimeoutSec, logger)
			log.Info("specpipe heartbeat routine stopped")
		}()
		<-sig
		httpSvr.Close()
		cancel()

		wg.Wait()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().StringP("http-server-port", "", "", "http server port")
	viper.BindPFlag("http.serverPort", rootCmd.Flags().Lookup("http-server-port"))

	rootCmd.Flags().StringP("nats-url", "", "", "NATS URL")
	viper.BindPFlag("nats.url", rootCmd.Flags().Lookup("nats-url"))

	rootCmd.Flags().StringP("heartbeat-timeout-sec", "", "", "heartbeat timeout seconds for device monitoring")
	viper.BindPFlag("heartbeat.timeoutSec", rootCmd.Flags().Lookup("heartbeat-timeout-sec"))

	rootCmd.Flags().StringP("log-level", "", "", "log level")
	viper.BindPFlag("logging.level", rootCmd.Flags().Lookup("log-level"))
}

func initConfig() {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func main() {
	Execute()
}
