package runtime

import (
	directorv1 "github.com/knita-io/knita/api/director/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

const (
	// TypeHost is a runtime that executes directly on the host executor.
	TypeHost Type = "host"
	// TypeDocker is a runtime that executes inside a Docker container.
	TypeDocker Type = "docker"
)

type Type string

// Opt configures a directorv1.OpenRequest.
type Opt func(opts *directorv1.OpenRequest)

// WithType specifies the type of runtime to open.
func WithType(t Type) Opt {
	return func(o *directorv1.OpenRequest) {
		switch t {
		case TypeDocker:
			o.Opts.Type = executorv1.RuntimeType_RUNTIME_DOCKER
		case TypeHost:
			o.Opts.Type = executorv1.RuntimeType_RUNTIME_HOST
		default:
			panic("unknown runtime type: " + string(t))
		}
	}
}

type Operator string

const (
	OperatorIn           Operator = "in"
	OperatorNotIn        Operator = "not-in"
	OperatorExists       Operator = "exists"
	OperatorDoesNotExist Operator = "not-exists"
)

// Requirement is a single expression in a selector.
type Requirement struct {
	Key      string
	Operator Operator
	Values   []string // used only by In/NotIn
}

// WithRunsOn specifies that this runtime can only be opened on an executor
// whose labels satisfy the given matchLabels AND matchExpressions.
func WithRunsOn(matchLabels map[string]string, matchExpressions ...Requirement) Opt {
	return func(o *directorv1.OpenRequest) {
		// convert to protobuf LabelSelector
		var exprs []*executorv1.LabelSelectorRequirement
		for _, r := range matchExpressions {
			var op executorv1.LabelSelectorRequirement_Operator
			switch r.Operator {
			case OperatorIn:
				op = executorv1.LabelSelectorRequirement_IN
			case OperatorNotIn:
				op = executorv1.LabelSelectorRequirement_NOT_IN
			case OperatorExists:
				op = executorv1.LabelSelectorRequirement_EXISTS
			case OperatorDoesNotExist:
				op = executorv1.LabelSelectorRequirement_DOES_NOT_EXIST
			default:
				panic("unknown label selector operator: " + string(r.Operator))
			}
			exprs = append(exprs, &executorv1.LabelSelectorRequirement{
				Key:      r.Key,
				Operator: op,
				Values:   r.Values,
			})
		}
		o.Opts.LabelSelector = &executorv1.LabelSelector{
			MatchLabels:      matchLabels,
			MatchExpressions: exprs,
		}
	}
}

// WithImage specifies the Docker image URI to use.
func WithImage(imageURI string) Opt {
	return func(o *directorv1.OpenRequest) {
		if o.Opts.Opts == nil {
			o.Opts.Opts = &executorv1.RuntimeOpts_Docker{Docker: &executorv1.DockerOpts{}}
		}
		d := o.Opts.GetDocker()
		if d.Image == nil {
			d.Image = &executorv1.DockerPullOpts{}
		}
		d.Image.ImageUri = imageURI
	}
}

type DockerPullStrategy string

const (
	PullStrategyAlways    DockerPullStrategy = "always"
	PullStrategyNever     DockerPullStrategy = "never"
	PullStrategyNotExists DockerPullStrategy = "not-exists"
)

// WithPullStrategy configures how Docker pulls the image.
func WithPullStrategy(ps DockerPullStrategy) Opt {
	return func(o *directorv1.OpenRequest) {
		if o.Opts.Opts == nil {
			o.Opts.Opts = &executorv1.RuntimeOpts_Docker{Docker: &executorv1.DockerOpts{}}
		}
		d := o.Opts.GetDocker()
		if d.Image == nil {
			d.Image = &executorv1.DockerPullOpts{}
		}
		switch ps {
		case PullStrategyAlways:
			d.Image.PullStrategy = executorv1.DockerPullOpts_PULL_STRATEGY_ALWAYS
		case PullStrategyNever:
			d.Image.PullStrategy = executorv1.DockerPullOpts_PULL_STRATEGY_NEVER
		case PullStrategyNotExists:
			d.Image.PullStrategy = executorv1.DockerPullOpts_PULL_STRATEGY_NOT_EXISTS
		default:
			panic("unknown pull strategy: " + string(ps))
		}
	}
}

// WithBasicAuth configures basic auth for Docker pulls.
func WithBasicAuth(username, password string) Opt {
	return func(o *directorv1.OpenRequest) {
		if o.Opts.Opts == nil {
			o.Opts.Opts = &executorv1.RuntimeOpts_Docker{Docker: &executorv1.DockerOpts{}}
		}
		d := o.Opts.GetDocker()
		if d.Image == nil {
			d.Image = &executorv1.DockerPullOpts{}
		}
		d.Image.Auth = &executorv1.DockerPullAuth{
			Auth: &executorv1.DockerPullAuth_Basic{
				Basic: &executorv1.BasicAuth{
					Username: username,
					Password: password,
				},
			},
		}
	}
}

// WithAWSECRAuth configures AWS ECR auth for Docker pulls.
func WithAWSECRAuth(region, accessKeyID, secretKey string) Opt {
	return func(o *directorv1.OpenRequest) {
		if o.Opts.Opts == nil {
			o.Opts.Opts = &executorv1.RuntimeOpts_Docker{Docker: &executorv1.DockerOpts{}}
		}
		d := o.Opts.GetDocker()
		if d.Image == nil {
			d.Image = &executorv1.DockerPullOpts{}
		}
		d.Image.Auth = &executorv1.DockerPullAuth{
			Auth: &executorv1.DockerPullAuth_AwsEcr{
				AwsEcr: &executorv1.AWSECRAuth{
					Region:         region,
					AwsAccessKeyId: accessKeyID,
					AwsSecretKey:   secretKey,
				},
			},
		}
	}
}

// WithDisplayName sets the display name for the runtime.
func WithDisplayName(displayName string) Opt {
	return func(o *directorv1.OpenRequest) {
		o.Opts.DisplayName = displayName
	}
}

// WithLabel sets a single label.
func WithLabel(key, value string) Opt {
	return WithLabels(key, value)
}

// WithLabels sets multiple labels from an alternating key/value list.
// Panics if you pass an odd number of args.
func WithLabels(kv ...string) Opt {
	m := KVMap("WithLabels", kv)
	return func(o *directorv1.OpenRequest) {
		o.Opts.Meta = MergeLabels(o.Opts.Meta, m)
	}
}

// WithAnnotation sets a single annotation.
func WithAnnotation(key, value string) Opt {
	return WithAnnotations(key, value)
}

// WithAnnotations sets multiple annotations from an alternating key/value list.
// Panics if you pass an odd number of args.
func WithAnnotations(kv ...string) Opt {
	m := KVMap("WithAnnotations", kv)
	return func(o *directorv1.OpenRequest) {
		o.Opts.Meta = MergeAnnotations(o.Opts.Meta, m)
	}
}
