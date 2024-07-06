package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/executor"
	"github.com/knita-io/knita/internal/server"
	"github.com/knita-io/knita/internal/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use:   "knita-executor",
	Short: "Starts the Knita Executor server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Silence usage on error once we're inside the RunE function, as
		// we know this must be a valid command invocation at this point.
		cmd.SilenceUsage = true
		syslog, err := makeLogger()
		if err != nil {
			return nil
		}
		configFilePath, _ := cmd.Flags().GetString("config")
		config, err := getConfig(syslog, configFilePath)
		if err != nil {
			return err
		}

		listener, err := net.Listen("tcp", config.BindAddress)
		if err != nil {
			return fmt.Errorf("error listening on tcp socket %s: %w", config.BindAddress, err)
		}
		defer listener.Close()

		executor := executor.NewServer(syslog, executor.Config{Name: config.Name, Labels: config.Labels})
		defer executor.Stop()

		srv := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				recovery.UnaryServerInterceptor(),
				server.MakeUnaryServerLogInterceptor(syslog.Named("grpc"))),
			grpc.ChainStreamInterceptor(
				recovery.StreamServerInterceptor(),
				server.MakeStreamServerLogInterceptor(syslog.Named("grpc"))))
		executorv1.RegisterExecutorServer(srv, executor)
		go func() {
			err := srv.Serve(listener)
			if err != nil {
				log.Fatal(err)
			}
		}()
		defer srv.Stop()

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		return nil
	},
}

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Prints the Knita version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(os.Stdout, version.Version)
		return nil
	},
}

func makeLogger() (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	zLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}
	return zLogger.Sugar(), nil
}

func main() {
	rootCmd.PersistentFlags().StringP("config", "c", defaultConfigFilePath, "Specify a custom path to the executor config file")
	rootCmd.AddCommand(versionCMD)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
