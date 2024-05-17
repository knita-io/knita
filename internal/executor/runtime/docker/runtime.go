package docker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/moby/moby/client"
	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/file"
)

type runtimeImageConfig struct {
	OS runtime.OS
}

type runtimeContainerConfig struct {
	Name              string
	GuestWorkspaceDir string
	PID0Command       []string
	Binds             []string
}

// Runtime executes jobs inside a Docker container.
type Runtime struct {
	file.WriteFS
	baseDir          string
	runtimeID        string
	opts             *executorv1.DockerOpts
	containerManager *ContainerManager
	syslog           *zap.SugaredLogger
	log              *runtime.Log
	state            struct {
		started         bool
		containerID     string
		imageConfig     runtimeImageConfig
		containerConfig runtimeContainerConfig
	}
}

func NewRuntime(log *zap.SugaredLogger, buildID string, runtimeID string, stream event.Stream, opts *executorv1.DockerOpts, client *client.Client) (*Runtime, error) {
	baseDir, err := os.MkdirTemp("", "knita-docker-*")
	if err != nil {
		return nil, fmt.Errorf("error creating runtime base dir: %w", err)
	}
	return &Runtime{
		syslog:           log.Named("docker_runtime"),
		runtimeID:        runtimeID,
		log:              runtime.NewLog(stream, buildID, runtimeID),
		baseDir:          baseDir,
		WriteFS:          file.WriteDirFS(baseDir),
		opts:             opts,
		containerManager: NewContainerManager(log, client),
	}, nil
}

// Start initializes the runtime and prepares it to have commands Exec'd inside it.
func (r *Runtime) Start(ctx context.Context) error {
	if r.state.started {
		return fmt.Errorf("error starting docker runtime: already started")
	}
	r.state.started = true

	pullLog := r.Log().Named("docker_pull")
	pullLog.Printf("Pulling Docker image...")
	err := r.containerManager.PullDockerImage(ctx, pullLog, r.opts.Image)
	if err != nil {
		return fmt.Errorf("error pulling Docker image: %w", err)
	}
	imageOS, err := r.containerManager.GetDockerImageOS(ctx, r.opts.Image.ImageUri)
	if err != nil {
		return fmt.Errorf("error discovering image OS: %w", err)
	}
	r.state.imageConfig.OS = imageOS
	config, err := r.prepareJobContainerConfig(ctx)
	if err != nil {
		return err
	}
	r.state.containerConfig = *config
	r.syslog.Infof("Guest OS: %s", r.state.imageConfig.OS)
	r.syslog.Infof("Guest Working dir: %s", config.GuestWorkspaceDir)
	r.syslog.Infof("Binds: %#v", config.Binds)
	cConfig := ContainerConfig{
		Name:       fmt.Sprintf("knita-%s", r.runtimeID),
		ImageURI:   r.opts.Image.ImageUri,
		Entrypoint: config.PID0Command,
		WorkingDir: config.GuestWorkspaceDir,
		Binds:      config.Binds,
		Networks:   []string{},
		// TODO stderr and stdout
	}
	containerID, err := r.containerManager.StartContainer(ctx, cConfig)
	if err != nil {
		return err
	}
	r.state.containerID = containerID
	return nil
}

func (r *Runtime) ID() string {
	return r.runtimeID
}

func (r *Runtime) Log() *runtime.Log {
	return r.log
}

// Close tears down the runtime.
func (r *Runtime) Close() error {
	if !r.state.started {
		return fmt.Errorf("error stopping docker runtime: not started")
	}
	var results *multierror.Error
	if r.state.containerID != "" {
		err := r.containerManager.StopContainer(context.TODO(), r.state.containerID)
		if err != nil {
			results = multierror.Append(results, fmt.Errorf("error stopping job container: %w", err))
		}
	}
	r.log.Close()
	r.state.started = false
	return results.ErrorOrNil()
}

// Exec executes a command inside the runtime.
// Start must have been called before calling Exec.
func (r *Runtime) Exec(ctx context.Context, execID string, opts *executorv1.ExecOpts) (*runtime.ExecResult, error) {
	r.syslog.Infow("Executing command", "name", opts.Name, "args", opts.Args)
	r.Log().ExecSource(execID, true).Printf("Executing command: %s %v", opts.Name, opts.Args)
	execLog := r.Log().ExecSource(execID, false)
	execConfig := ExecConfig{
		ContainerID: r.state.containerID,
		Command:     append([]string{opts.Name}, opts.Args...),
		WorkingDir:  r.state.containerConfig.GuestWorkspaceDir,
		Env:         r.fixEnv(opts.Env),
	}

	w := execLog.Stdout()
	defer w.Close()
	execConfig.Stdout = w

	w = execLog.Stderr()
	defer w.Close()
	execConfig.Stderr = w

	err := r.containerManager.Execute(ctx, execConfig)
	if err != nil {
		var exitErr *exitError
		if errors.As(err, &exitErr) {
			return &runtime.ExecResult{ExitCode: int32(exitErr.exitCode)}, nil
		}
		return nil, fmt.Errorf("error running command: %w", err)
	}
	return &runtime.ExecResult{ExitCode: 0}, nil
}

func (r *Runtime) prepareJobContainerConfig(ctx context.Context) (*runtimeContainerConfig, error) {
	switch r.state.imageConfig.OS {
	case runtime.OSLinux:
		return r.prepareLinuxContainerConfig(ctx)
	case runtime.OSWindows:
		return r.prepareWindowsContainerConfig(ctx)
	default:
		return nil, fmt.Errorf("error unsupported image OS: %v", r.state.imageConfig.OS)
	}
}

func (r *Runtime) prepareWindowsContainerConfig(ctx context.Context) (*runtimeContainerConfig, error) {
	guestWorkingDir := "C:\\controlci\\workspace"
	binds := []string{
		fmt.Sprintf("%s:%s:rw", r.baseDir, guestWorkingDir),
		// Windows containers only run on Windows, so use the Windows pipe syntax
		"\\\\.\\pipe\\docker_engine:\\\\.\\pipe\\docker_engine",
	}
	return &runtimeContainerConfig{
		Name:              r.runtimeID,
		Binds:             binds,
		GuestWorkspaceDir: guestWorkingDir,
		PID0Command:       []string{"timeout", "/t", "-1"},
	}, nil
}

func (r *Runtime) prepareLinuxContainerConfig(ctx context.Context) (*runtimeContainerConfig, error) {
	guestWorkingDir := "/tmp/controlci/workspace"
	binds := []string{
		fmt.Sprintf("%s:%s:rw", r.baseDir, guestWorkingDir),
		// Linux containers run natively on Linux, and in a Linux VM on Windows and macOS,
		// so we can always refer to the Linux socket path here
		"/var/run/docker.sock:/var/run/docker.sock",
	}
	return &runtimeContainerConfig{
		Name:              r.runtimeID,
		Binds:             binds,
		GuestWorkspaceDir: guestWorkingDir,
		PID0Command:       []string{"/bin/sh", "-c", "while :; do sleep 2073600; done"},
	}, nil
}

func (r *Runtime) fixEnv(env []string) []string {
	for i, envVar := range env {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			path, changed, err := r.mapHostPath(runtime.GetHostOS(), parts[1])
			if err != nil {
				r.syslog.Warnf("Ignoring error mapping host path for %q: %v", parts[1], err)
			} else if changed {
				env[i] = fmt.Sprintf("%s=%s", parts[0], path)
			}
		}
	}
	return env
}

func (r *Runtime) mapHostPath(fromOS runtime.OS, path string) (string, bool, error) {
	if r.state.imageConfig.OS == "" {
		return "", false, fmt.Errorf("error runtime is not prepared")
	}
	var changed bool
	if strings.HasPrefix(path, r.baseDir) {
		path = strings.Replace(path, r.baseDir, r.state.containerConfig.GuestWorkspaceDir, 1)
		changed = true
	}
	if changed {
		switch fromOS {
		case runtime.OSMacOS:
			// macOS can only run Linux containers.
			// macOS's paths are compatible with Linux.
		case runtime.OSLinux:
			// Linux can only run Linux containers.
			// Linux to Linux does not need a conversion.
		case runtime.OSWindows:
			// Windows can run Windows or Linux containers.
			switch r.state.imageConfig.OS {
			case runtime.OSLinux:
				// Windows to Linux needs path separator tweaking after we've swapped out the path above.
				path = strings.Replace(path, "\\", "/", -1)
			case runtime.OSWindows:
				// Windows to Windows does not need a conversion.
			default:
				return "", false, fmt.Errorf("error unsupported container OS: %v", fromOS)
			}
		default:
			return "", false, fmt.Errorf("error unsupported OS: %v", fromOS)
		}
	}
	return path, changed, nil
}
