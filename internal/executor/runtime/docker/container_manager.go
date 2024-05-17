package docker

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/errdefs"
	"github.com/hashicorp/go-multierror"
	"github.com/moby/moby/client"
	"github.com/moby/moby/pkg/stdcopy"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	runtime2 "github.com/knita-io/knita/internal/executor/runtime"
)

type ContainerConfig struct {
	Name       string
	ImageURI   string
	Entrypoint []string
	Command    []string
	WorkingDir string
	Env        []string
	Binds      []string
	// Networks is a list of network IDs the container should be connected to.
	// The container will always be connected to the host network by default.
	Networks []string
	// Aliases is a list of names this container will be resolvable on from
	// each of the configured networks (if any).
	Aliases []string
	Stdout  io.Writer
	Stderr  io.Writer
}

type ExecConfig struct {
	ContainerID string
	Command     []string
	WorkingDir  string
	Env         []string
	Stdout      io.Writer
	Stderr      io.Writer
}

type ContainerManager struct {
	client *client.Client
	log    *zap.SugaredLogger
}

func NewContainerManager(log *zap.SugaredLogger, client *client.Client) *ContainerManager {
	return &ContainerManager{
		client: client,
		log:    log.Named("container_manager"),
	}
}

// PullDockerImage pulls a Docker Image from a remote registry.
func (r *ContainerManager) PullDockerImage(ctx context.Context, log *runtime2.Log, opts *executorv1.DockerPullOpts) error {
	imageURI := parseDockerImageURI(opts.ImageUri)
	fil := filters.NewArgs()
	fil.Add("reference", imageURI.Reference())
	list, err := r.client.ImageList(ctx, types.ImageListOptions{
		All:     false,
		Filters: fil,
	})
	if err != nil {
		return fmt.Errorf("error listing images: %w", err)
	}

	alreadyExists := len(list) > 0
	if opts.PullStrategy == executorv1.DockerPullOpts_PULL_STRATEGY_NEVER {
		log.Printf("Docker pull strategy is %q; %q will not be pulled",
			DockerPullStrategyNever, imageURI.FQN())
		return nil
	}
	if opts.PullStrategy == executorv1.DockerPullOpts_PULL_STRATEGY_NOT_EXISTS && alreadyExists {
		log.Printf("Docker pull strategy is %q and image exists in cache; %q will not be pulled",
			DockerPullStrategyIfNotExists, imageURI.FQN())
		return nil
	}
	if alreadyExists && !imageURI.IsLatest() && opts.PullStrategy == executorv1.DockerPullOpts_PULL_STRATEGY_UNSPECIFIED {
		log.Printf("Docker pull strategy is %q, image exists in cache and is not latest; %q will not be pulled",
			DockerPullStrategyDefault, imageURI.FQN())
		return nil
	}

	log.Printf("Pulling image: %s", imageURI)

	// If authentication has been provided then pass it into the image pull
	imagePullOptions := types.ImagePullOptions{}
	if opts.Auth != nil && opts.Auth.GetBasic() != nil {
		log.Printf("Using Docker registry auth: Basic")
		jsonBytes, err := protojson.Marshal(opts.Auth.GetBasic())
		if err != nil {
			return fmt.Errorf("error encoding docker auth: %w", err)
		}
		imagePullOptions.RegistryAuth = base64.StdEncoding.EncodeToString(jsonBytes)
	} else if opts.Auth != nil && opts.Auth.GetAwsEcr() != nil {
		awsECR := opts.Auth.GetAwsEcr()
		log.Printf("Using Docker registry auth: AWS")
		cfg := &aws.Config{}
		if awsECR.Region != "" {
			cfg = cfg.WithRegion(awsECR.Region)
		}
		cfg = cfg.WithCredentials(credentials.NewStaticCredentials(awsECR.AwsAccessKeyId, awsECR.AwsSecretKey, ""))
		sess, err := session.NewSession(cfg)
		if err != nil {
			return fmt.Errorf("error creating AWS session: %w", err)
		}
		svc := ecr.New(sess)
		token, err := svc.GetAuthorizationTokenWithContext(ctx, &ecr.GetAuthorizationTokenInput{})
		if err != nil {
			return fmt.Errorf("error getting AWS ECR authorization token: %w", err)
		}
		if len(token.AuthorizationData) == 0 {
			return fmt.Errorf("error unexpected AWS ECR token format")
		}
		authData := token.AuthorizationData[0].AuthorizationToken
		data, err := base64.StdEncoding.DecodeString(*authData)
		if err != nil {
			return fmt.Errorf("error decoding AWS ECR token: %w", err)
		}
		parts := strings.SplitN(string(data), ":", 2)
		if len(parts) < 2 {
			return fmt.Errorf("error unexpected AWS ECR token data format")
		}
		basic := &executorv1.BasicAuth{
			Username: "AWS",
			Password: parts[1],
		}
		jsonBytes, err := protojson.Marshal(basic)
		if err != nil {
			return fmt.Errorf("error encoding docker auth: %w", err)
		}
		imagePullOptions.RegistryAuth = base64.StdEncoding.EncodeToString(jsonBytes)
	} else {
		log.Printf("Using Docker registry auth: None")
	}

	// TODO this error needs to go to the job log.go
	stream, err := r.client.ImagePull(ctx, imageURI.FQN(), imagePullOptions)
	if err != nil {
		return errors.Wrap(err, "error pulling image")
	}
	defer stream.Close()

	// TODO can output image pull info to the build logs here:
	// 	Do a list on the image to discover its size
	//  Intercept stream and output some progress information as it's being read
	//  Make this a generic util and use it for other pulls
	res, err := r.client.ImageLoad(ctx, stream, false)
	if err != nil {
		return errors.Wrap(err, "error loading image")
	}
	defer res.Body.Close()
	return nil
}

// GetDockerImageOS returns the type of underlying guest OS the specified Docker image
// is made from. The docker image must have been pulled first.
func (r *ContainerManager) GetDockerImageOS(ctx context.Context, imageURI string) (runtime2.OS, error) {
	image := parseDockerImageURI(imageURI)
	inspect, _, err := r.client.ImageInspectWithRaw(ctx, image.String())
	if err != nil {
		return "", fmt.Errorf("error inspecting image %q: %w", image, err)
	}
	return runtime2.OS(inspect.Os), nil
}

// StartContainer starts a new container in the background and returns its unique ID.
// Call StopContainer to stop it and free up resources.
func (r *ContainerManager) StartContainer(ctx context.Context, config ContainerConfig) (string, error) {
	image := parseDockerImageURI(config.ImageURI)
	cConfig := &container.Config{
		Image:      image.FQN(),
		Entrypoint: config.Entrypoint,
		Cmd:        config.Command,
		WorkingDir: config.WorkingDir,
		Env:        config.Env,
	}
	hConfig := &container.HostConfig{
		AutoRemove: false,
		Binds:      config.Binds,
	}
	nConfig := &network.NetworkingConfig{}
	res, err := r.client.ContainerCreate(ctx, cConfig, hConfig, nConfig, nil, config.Name) // platform is optional
	if err != nil {
		return "", errors.Wrap(err, "error creating container")
	}
	for _, networkID := range config.Networks {
		nConfig := &network.EndpointSettings{Aliases: config.Aliases}
		err = r.client.NetworkConnect(ctx, networkID, res.ID, nConfig)
		if err != nil {
			return "", fmt.Errorf("error connecting container to network: %w", err)
		}
	}
	err = r.client.ContainerStart(ctx, res.ID, container.StartOptions{})
	if err != nil {
		return "", errors.Wrap(err, "error starting container")
	}
	if config.Stdout != nil || config.Stderr != nil {
		ctx := context.TODO()
		opts := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: false}
		reader, err := r.client.ContainerLogs(ctx, res.ID, opts)
		if err != nil {
			return "", fmt.Errorf("error connecting to container log.go stream: %w", err)
		}
		r.pipeContainerLogAsync(reader, config.Stdout, config.Stderr)
	}
	return res.ID, nil
}

// StopContainer stops and removes a previously started docker container.
func (r *ContainerManager) StopContainer(ctx context.Context, containerID string) error {
	var results *multierror.Error
	err := r.client.ContainerKill(ctx, containerID, "kill")
	if err != nil && !errdefs.IsNotFound(err) {
		results = multierror.Append(results, fmt.Errorf("error killing container %q: %w", containerID, err))
	}
	err = r.client.ContainerRemove(ctx, containerID, container.RemoveOptions{RemoveVolumes: true, Force: true})
	if err != nil && !errdefs.IsNotFound(err) {
		results = multierror.Append(results, fmt.Errorf("error removing container %q: %w", containerID, err))
	}
	return results.ErrorOrNil()
}

// Execute a command inside the container.
// StartContainer must have previously been called.
func (r *ContainerManager) Execute(ctx context.Context, config ExecConfig) error {
	eConfig := types.ExecConfig{
		Cmd:          config.Command,
		Env:          config.Env,
		WorkingDir:   config.WorkingDir,
		Detach:       false,
		AttachStderr: true,
		AttachStdout: true,
	}
	createRes, err := r.client.ContainerExecCreate(ctx, config.ContainerID, eConfig)
	if err != nil {
		return fmt.Errorf("error creating exec: %w", err)
	}
	resp, err := r.client.ContainerExecAttach(ctx, createRes.ID, types.ExecStartCheck{})
	if err != nil {
		return fmt.Errorf("error attaching exec: %w", err)
	}
	defer resp.Close()
	if config.Stdout != nil || config.Stderr != nil {
		err = r.pipeContainerLog(resp.Reader, config.Stdout, config.Stderr)
		if err != nil {
			return fmt.Errorf("error piping container log.go: %w", err)
		}
	}
	var exitCode int
	for {
		res, err := r.client.ContainerExecInspect(ctx, createRes.ID)
		if err != nil {
			return fmt.Errorf("error inspecting script exec: %w", err)
		}
		if res.Running {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		exitCode = res.ExitCode
		break
	}
	if exitCode != 0 {
		return &exitError{exitCode: exitCode}
	}
	return nil
}

type exitError struct {
	err      error
	exitCode int
}

func (e *exitError) Error() string {
	return e.err.Error()
}

// CreateNetwork creates a new private network and returns its ID.
func (r *ContainerManager) CreateNetwork(ctx context.Context, name string) (string, error) {
	res, err := r.client.NetworkCreate(ctx, name, types.NetworkCreate{})
	if err != nil {
		return "", fmt.Errorf("error creating network: %w", err)
	}
	return res.ID, nil
}

// DeleteNetwork deletes a previously created network.
func (r *ContainerManager) DeleteNetwork(ctx context.Context, networkID string) error {
	err := r.client.NetworkRemove(ctx, networkID)
	if err != nil {
		return errors.Wrap(err, "error removing network")
	}
	return nil
}

func (r *ContainerManager) pipeContainerLog(from io.Reader, stdout io.Writer, stderr io.Writer) error {
	// https://github.com/docker/cli/blob/master/cli/command/container/logs.go
	// https://github.com/docker/cli/blob/ebca1413117a3fcb81c89d6be226dcec74e5289f/vendor/github.com/docker/docker/pkg/stdcopy/stdcopy.go#L94
	_, err := stdcopy.StdCopy(stdout, stderr, from)
	if err != nil && err != io.EOF && !errors.Is(err, io.ErrClosedPipe) {
		return err
	}
	return nil
}

func (r *ContainerManager) pipeContainerLogAsync(from io.Reader, stdout io.Writer, stderr io.Writer) <-chan struct{} {
	doneC := make(chan struct{})
	go func() {
		defer close(doneC)
		err := r.pipeContainerLog(from, stdout, stderr)
		if err != nil {
			r.log.Warnf("Ignoring error piping container logs; Logs may be incomplete: %s", err)
		}
	}()
	return doneC
}
