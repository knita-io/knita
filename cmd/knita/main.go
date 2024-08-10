package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	brokerv1 "github.com/knita-io/knita/api/broker/v1"
	directorv1 "github.com/knita-io/knita/api/director/v1"
	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/cmd/knita/ui"
	"github.com/knita-io/knita/internal/broker/fixed"
	"github.com/knita-io/knita/internal/director"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor"
	"github.com/knita-io/knita/internal/file"
	"github.com/knita-io/knita/internal/server"
	"github.com/knita-io/knita/internal/version"
)

var rootCmd = &cobra.Command{
	Use: "knita",
}

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Prints the Knita version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(os.Stdout, version.Version)
		return nil
	},
}

var buildCMD = &cobra.Command{
	Use:   "build [pattern command]",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
	Short: "Executes the specified build pattern",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Silence usage on error once we're inside the RunE function, as
		// we know this must be a valid command invocation at this point.
		cmd.SilenceUsage = true

		now := time.Now()
		buildID := uuid.New().String()

		verbose, _ := cmd.Flags().GetBool("verbose")
		if !isatty.IsTerminal(os.Stdout.Fd()) {
			verbose = true
		}

		var (
			err        error
			cliLogPath string
		)
		cliLogPathOverride := os.Getenv("KNITA_CLI_LOG_PATH")
		if cliLogPathOverride != "" {
			cliLogPath = cliLogPathOverride
		} else {
			cliLogPath, err = makeLogFile("knita-", now)
			if err != nil {
				return fmt.Errorf("error making cli log file: %w", err)
			}
		}
		syslog, err := makeLogger(cliLogPath)
		if err != nil {
			return nil
		}
		defer syslog.Sync()

		configFilePath, _ := cmd.Flags().GetString("config")
		config, err := getConfig(syslog, configFilePath)
		if err != nil {
			return err
		}

		work, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %w", err)
		}
		syslog.Infof("Working directory: %v", work)

		var (
			buildOut     io.Writer
			buildLogPath string
		)
		if !verbose {
			buildLogPath, err = makeLogFile("knita-build-", now)
			if err != nil {
				return fmt.Errorf("error making log file: %w", err)
			}
			file, err := os.OpenFile(buildLogPath, os.O_WRONLY, 0)
			if err != nil {
				return fmt.Errorf("error opening build log for writing: %w", err)
			}
			defer file.Close()
			buildOut = file
		} else {
			buildOut = os.Stdout
		}

		socket, err := getSocketPath()
		if err != nil {
			return err
		}
		defer os.Remove(socket)

		// NOTE: This will work on Linux and macOS going way back, but only on Windows 10+.
		listener, err := net.Listen("unix", socket)
		if err != nil {
			return fmt.Errorf("error listening on unix socket %s: %w", socket, err)
		}
		defer listener.Close()

		dialer := func(addr string, t time.Duration) (net.Conn, error) {
			return net.Dial("unix", addr)
		}
		conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithDialer(dialer))
		if err != nil {
			return fmt.Errorf("error dialing local knita socket %s: %w", socket, err)
		}
		brokerClient := brokerv1.NewRuntimeBrokerClient(conn)
		buildLog := director.NewLog(event.NewBroker(syslog), buildID)
		defer buildLog.Close()
		build := director.NewBuild(syslog, buildLog, buildID, brokerClient, file.WriteDirFS(work))
		directorServer := director.NewServer(syslog, build)
		embeddedExecutorName, _ := os.Hostname()
		if embeddedExecutorName == "" {
			embeddedExecutorName = "knita-exec-local"
		}

		var (
			executorSrv *executor.Server
			executors   []*fixed.ExecutorConfig
		)
		if !config.Executors.Local.Disabled {
			executorSrv = executor.NewServer(syslog, executor.Config{Name: embeddedExecutorName, Labels: config.Executors.Local.Labels})
			defer executorSrv.Stop()
			executors = append(executors, &fixed.ExecutorConfig{
				Connection: &brokerv1.RuntimeConnectionInfo{
					Transport: &brokerv1.RuntimeConnectionInfo_Unix{
						Unix: &brokerv1.RuntimeTransportUnix{SocketPath: socket},
					},
				},
			})
		} else {
			syslog.Warnf("Local builds are disabled")
		}
		for _, execConfig := range config.Executors.Remote {
			if execConfig.Disabled {
				continue
			}
			executors = append(executors, &fixed.ExecutorConfig{
				Connection: &brokerv1.RuntimeConnectionInfo{
					Transport: &brokerv1.RuntimeConnectionInfo_Tcp{
						Tcp: &brokerv1.RuntimeTransportTCP{Address: execConfig.Address},
					},
				},
			})
		}
		broker := fixed.NewServer(syslog, fixed.Config{Executors: executors})

		srv := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				recovery.UnaryServerInterceptor(),
				server.MakeUnaryServerLogInterceptor(syslog.Named("grpc"))),
			grpc.ChainStreamInterceptor(
				recovery.StreamServerInterceptor(),
				server.MakeStreamServerLogInterceptor(syslog.Named("grpc"))))
		if executorSrv != nil {
			executorv1.RegisterExecutorServer(srv, executorSrv)
		}
		brokerv1.RegisterRuntimeBrokerServer(srv, broker)
		directorv1.RegisterDirectorServer(srv, directorServer)

		go func() {
			err := srv.Serve(listener)
			if err != nil {
				log.Fatal(err)
			}
		}()
		defer srv.Stop()

		var uiManager *ui.Manager
		if !verbose {
			uiManager = ui.NewManager(buildLog.Stream())
			uiManager.Start()
			defer uiManager.Stop()
		}

		buildLog.Stream().Subscribe(func(event *event.Event) {
			switch p := event.Payload.(type) {
			case *builtinv1.StdoutEvent:
				buildOut.Write([]byte(fmt.Sprintf("%s", string(p.Data))))
			case *builtinv1.StderrEvent:
				buildOut.Write([]byte(fmt.Sprintf("%s", string(p.Data))))
			}
		})

		env := append([]string{}, os.Environ()...)
		env = append(env, []string{
			fmt.Sprintf("KNITA_SOCKET=%s", socket),
			fmt.Sprintf("KNITA_BUILD_ID=%s", buildID),
		}...)

		err = build.Run(func() error {
			execCmd := exec.Command(args[0], args[:1]...)
			execCmd.Env = env
			execCmd.Stdout = buildLog.Stdout()
			execCmd.Stderr = buildLog.Stderr()
			err := execCmd.Run()
			if err != nil {
				return fmt.Errorf("error running command: %w", err)
			}
			return nil
		})
		if !verbose {
			uiManager.Stop()
			fmt.Fprintf(os.Stdout, "\nBuild log available at: %s\n", buildLogPath)
		}
		if err != nil {
			fmt.Fprintf(os.Stdout, "\n")
			return err
		}
		return nil
	},
}

func makeLogger(logPath string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	cfg.OutputPaths = []string{logPath}
	zLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}
	return zLogger.Sugar(), nil
}

func makeLogFile(prefix string, ts time.Time) (string, error) {
	temp := os.TempDir()
	logDirectory := filepath.Join(temp, "knita")
	err := os.MkdirAll(logDirectory, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating log directory: %w", err)
	}
	tsStr := strings.Replace(strings.Replace(ts.UTC().Format(time.RFC3339), ":", "", -1), "-", "", -1)
	logPath := filepath.Join(logDirectory, fmt.Sprintf("%s%s.log", prefix, tsStr))
	file, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("error creating log file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("error closing log file: %w", err)
	}
	return logPath, nil
}

func getSocketPath() (string, error) {
	temp, err := os.CreateTemp(os.TempDir(), "knita-cli-*.socket")
	if err != nil {
		return "", fmt.Errorf("error creating socket: %w", err)
	}
	socket := temp.Name()
	err = temp.Close()
	if err != nil {
		return "", fmt.Errorf("error closing temp socket: %w", err)
	}
	err = os.Remove(socket)
	if err != nil {
		return "", fmt.Errorf("error removing temp socket: %w", err)
	}
	return socket, nil
}

func main() {
	buildCMD.PersistentFlags().BoolP("verbose", "v", false, "Set to true to disable the pretty build UI and send the build log directly to stdout")
	buildCMD.PersistentFlags().StringP("config", "c", defaultConfigFilePath, "Specify a custom path to the Knita config file")
	rootCmd.AddCommand(buildCMD)
	rootCmd.AddCommand(versionCMD)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
