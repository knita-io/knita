package runtime

import (
	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

const (
	// TypeHost is a runtime that executes directly on the host executor without any containerization or virtualization.
	TypeHost Type = "host"
	// TypeDocker is a runtime that executes inside a Docker container.
	TypeDocker Type = "docker"
)

type Type string

type Opt interface {
	Apply(opts *executorv1.Opts)
}

type withType struct {
	runtimeType executorv1.RuntimeType
}

func (o *withType) Apply(opts *executorv1.Opts) {
	opts.Type = o.runtimeType
}

// WithType specifies the type of runtime to open.
func WithType(t Type) Opt {
	switch t {
	case TypeDocker:
		return &withType{runtimeType: executorv1.RuntimeType_RUNTIME_DOCKER}
	case TypeHost:
		return &withType{runtimeType: executorv1.RuntimeType_RUNTIME_HOST}
	default:
		panic("Unknown runtime type: " + t)
	}
}

type withLabel struct {
	label string
}

func (o *withLabel) Apply(opts *executorv1.Opts) {
	opts.Labels = append(opts.Labels, o.label)
}

// WithLabel specifies that this runtime can only be opened on an executor that supports the specified label.
func WithLabel(label string) Opt {
	return &withLabel{label: label}
}

type withLabels struct {
	labels []string
}

func (o *withLabels) Apply(opts *executorv1.Opts) {
	opts.Labels = append(opts.Labels, o.labels...)
}

// WithLabels specifies that this runtime can only be opened on an executor that supports the specified labels.
func WithLabels(labels ...string) Opt {
	return &withLabels{labels: labels}
}

type withTag struct {
	key   string
	value string
}

func (o *withTag) Apply(opts *executorv1.Opts) {
	if opts.Tags == nil {
		opts.Tags = map[string]string{}
	}
	opts.Tags[o.key] = o.value
}

// WithTag tags the runtime with a key value pair.
func WithTag(key string, value string) Opt {
	return &withTag{key: key, value: value}
}

type withTags struct {
	tags map[string]string
}

func (o *withTags) Apply(opts *executorv1.Opts) {
	if opts.Tags == nil {
		opts.Tags = map[string]string{}
	}
	for k, v := range o.tags {
		opts.Tags[k] = v
	}
}

// WithTags tags the runtime with one or more key value pairs.
// The expected format of tags is an array of alternating key value pairs.
// If tags are mismatched this function will panic.
func WithTags(tags ...string) Opt {
	m := map[string]string{}
	if len(tags)%2 != 0 {
		panic("expected alternating kv paris")
	}
	return &withTags{tags: m}
}

type withImage struct {
	imageURI string
}

func (o *withImage) Apply(opts *executorv1.Opts) {
	if opts.Opts == nil {
		opts.Opts = &executorv1.Opts_Docker{Docker: &executorv1.DockerOpts{}}
	}
	if opts.GetDocker().GetImage() == nil {
		opts.GetDocker().Image = &executorv1.DockerPullOpts{}
	}
	opts.GetDocker().Image.ImageUri = o.imageURI
}

// WithImage specifies the image to use when the runtime type is Docker.
func WithImage(imageURI string) Opt {
	return &withImage{imageURI: imageURI}
}

const (
	// PullStrategyAlways ensures that a Docker pull is performed prior to starting the runtime.
	PullStrategyAlways DockerPullStrategy = "always"
	// PullStrategyNever will skip the Docker pull entirely. This is useful if executors are
	// loaded with relevant Docker images out of band of any individual build.
	PullStrategyNever DockerPullStrategy = "never"
	// PullStrategyNotExists will run Docker pull only if no matching image is found on the executor.
	PullStrategyNotExists DockerPullStrategy = "not-exists"
)

type DockerPullStrategy string

type withPullStrategy struct {
	pullStrategy executorv1.DockerPullOpts_PullStrategy
}

func (o *withPullStrategy) Apply(opts *executorv1.Opts) {
	if opts.Opts == nil {
		opts.Opts = &executorv1.Opts_Docker{Docker: &executorv1.DockerOpts{}}
	}
	if opts.GetDocker().GetImage() == nil {
		opts.GetDocker().Image = &executorv1.DockerPullOpts{}
	}
	opts.GetDocker().Image.PullStrategy = o.pullStrategy
}

// WithPullStrategy specifies the behaviour of the docker pull command when the runtime type is Docker.
func WithPullStrategy(pullStrategy DockerPullStrategy) Opt {
	switch pullStrategy {
	case PullStrategyAlways:
		return &withPullStrategy{pullStrategy: executorv1.DockerPullOpts_PULL_STRATEGY_ALWAYS}
	case PullStrategyNever:
		return &withPullStrategy{pullStrategy: executorv1.DockerPullOpts_PULL_STRATEGY_NEVER}
	case PullStrategyNotExists:
		return &withPullStrategy{pullStrategy: executorv1.DockerPullOpts_PULL_STRATEGY_NOT_EXISTS}
	default:
		panic("Unknown pull strategy: " + pullStrategy)
	}
}

type withBasicAuth struct {
	basicAuth *executorv1.BasicAuth
}

func (o *withBasicAuth) Apply(opts *executorv1.Opts) {
	if opts.GetDocker() == nil {
		opts.Opts = &executorv1.Opts_Docker{Docker: &executorv1.DockerOpts{}}
	}
	if opts.GetDocker().GetImage() == nil {
		opts.GetDocker().Image = &executorv1.DockerPullOpts{}
	}
	opts.GetDocker().Image.Auth = &executorv1.DockerPullAuth{Auth: &executorv1.DockerPullAuth_Basic{Basic: o.basicAuth}}
}

// WithBasicAuth configures basic auth for the docker pull command when the runtime type is Docker.
func WithBasicAuth(username string, password string) Opt {
	return &withBasicAuth{basicAuth: &executorv1.BasicAuth{Username: username, Password: password}}
}

type withAWSAuth struct {
	awsECRAuth *executorv1.AWSECRAuth
}

func (o *withAWSAuth) Apply(opts *executorv1.Opts) {
	if opts.GetDocker() == nil {
		opts.Opts = &executorv1.Opts_Docker{Docker: &executorv1.DockerOpts{}}
	}
	if opts.GetDocker().GetImage() == nil {
		opts.GetDocker().Image = &executorv1.DockerPullOpts{}
	}
	opts.GetDocker().Image.Auth = &executorv1.DockerPullAuth{Auth: &executorv1.DockerPullAuth_AwsEcr{AwsEcr: o.awsECRAuth}}
}

// WithAWSECRAuth configures AWS ECR auth for the docker pull command when the runtime type is Docker.
func WithAWSECRAuth(region string, awsAccessKeyID string, awsSecretKey string) Opt {
	return &withAWSAuth{awsECRAuth: &executorv1.AWSECRAuth{Region: region, AwsAccessKeyId: awsAccessKeyID, AwsSecretKey: awsSecretKey}}
}
