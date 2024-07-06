<div align="center">

# Knita
## The Distributed Build Engine

</div>

<div align="center">
  <img src="https://github.com/knita-io/knita/raw/main/docs/images/knita-build-demo.gif" alt="Knita Build Demo"/>
</div>

## Overview

Knita overhauls the traditional software build process by combining build and continuous integration into one cohesive platform.

* **Real Code**: Replace cumbersome CI YAML files with real code. Matrices are just for loops, conditions are just if statements, input/outputs are just variables etc.
* **Local Builds**: Run and test builds entirely locally without the painful change-commit-wait cycle typical of traditional CI systems.
* **Distributed Builds**: Run builds across distributed build infrastructure, even when running outside your CI environment. Mix and match your local machine with remote build servers to minimize queue time.
* **Flexible Environments**: Knita can run builds in a variety of different runtime environments. Docker and direct host execution is currently supported, with VM, Kubernetes and Podman planned to follow.
* **Dynamic Builds**: Builds are now just code - no more static YAML files. Adapt the behaviour of your builds at runtime to achieve:
  * **Adaptive Test Splitting**: Dynamically calculate the distribution of tests across multiple parallel executors to optimize run time.
  * **Conditional Retries**: Automatically re-run failed build targets based on their outputs, such as standard output or specific error messages.
  * **Manual Introspection**: Automatically pause builds to manually attach a debugger or to inspect the environment via SSH. Send an email or Slack message to relevant engineers to let them know that hard-to-trigger bug is waiting for them to root cause.
  * **External API Calls**: Enhance build orchestration with API calls to external systems, such as reserving hardware in specialized environments.

Anything you can code, Knita can execute as part of the build process.

## Documentation


### Getting Started

Download the latest Knita CLI from the [release page](https://github.com/knita-io/knita/releases) and make sure it's in your path.

To define your first pattern, see the getting started guide for your preferred language:
* [Golang](docs/guides/sdk/go/getting-started.md)
* [Python](docs/guides/sdk/python/getting-started.md)

_Don't see your language? Open a GitHub issue to request it._

### Config Reference

* [Knita CLI](docs/guides/cli/config.md)
* [Knita Executor](docs/guides/executor/config.md)

## How It Works

Knita has three main components:

* **Director**: Orchestrates the build lifecycle, and is programmed using the Knita SDK.
* **Executor**: Executes builds, similar to runners or agents in CI systems.
* **Broker**: Connects the Director to Executors.

The Knita CLI includes a Director, Broker, and Executor, so you can run builds entirely on your local machine. For more scalability, it can connect to external Brokers and Executors for distributed builds.

Builds are defined using the Knita SDK in your preferred programming language. The code for builds is called a pattern, which runs as an executable process and interacts with the Knita CLI.
<div align="center">
  <img src="https://github.com/knita-io/knita/raw/main/docs/images/knita-architecture.png" width="700" height="auto" alt="Knita Architecture"/>
</div>

### SDK Example

The example below has been extracted from Knita's own build code and annotated with helpful comments. For the full and functional code see [`build/pattern.go`](build/pattern.go) (yes, Knita is used to build Knita!).

```go
package main

import (
	"fmt"

	"github.com/knita-io/knita/sdk/go/knita"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
	"github.com/knita-io/knita/sdk/go/knita/runtime/exec"
)

func main() {
	// Get a handle on the Knita SDK client, which will be automatically configured to communicate with
	// the Knita CLI process
	client := knita.MustNewClient()
	// Execute the build steps
	// NOTE: Knita's protobufs must be generated before the binaries can be compiled. That dependency
	// is expressed here simply by running the build targets in series. More complex parallelized builds
	// could be achieved simply by using goroutines.
	generateProtos(client)
	buildBinaries(client)
}

func generateProtos(client *knita.Client) {
	// Start a new Docker-based runtime using the Golang Docker image
	golang := client.MustRuntime(
		runtime.WithTag("name", "generate"),
		runtime.WithType("docker"),
		runtime.WithImage("golang:1.22"))
	defer golang.MustClose()

	// Import the Knita protobuf definitions into the runtime
	golang.MustImport("internal/api/**/*.proto", "")

	// Execute the protobuf compiler
	golang.MustExec(
		exec.WithTag("name", "protobuf"),
		exec.WithCommand("/bin/bash", "-c", `
                  protoc \
                  --proto_path=internal/api \
                  --go_out=internal/api \
                  --go_opt=paths=source_relative \
                  --go-grpc_out=internal/api \
                  --go-grpc_opt=paths=source_relative \
                  broker/v1/broker.proto \
                  executor/v1/executor.proto \
                  director/v1/director.proto \
                  observer/v1/observer.proto`))

	// Export the generated Golang protobuf models back into the local source tree
	golang.MustExport("internal/api/**/*.pb.go", "")
}

func buildBinaries(client *knita.Client) {
	// Start a new Docker-based runtime using the Golang Docker image
	golang := client.MustRuntime(
		runtime.WithTag("name", "build"),
		runtime.WithType("docker"),
		runtime.WithImage("golang:1.22"))
	defer golang.MustClose()

	// Import the entire source tree into the runtime
	golang.MustImport(".", ".")

	// Cross compile the Knita CLI for each target OS and Architecture
	// NOTE: Matrices are expressed by using a for loop.
	for _, target := range []struct {
		os   string
		arch []string
	}{
		{os: "darwin", arch: []string{"arm64", "amd64"}},
		{os: "linux", arch: []string{"arm", "arm64", "amd64"}},
		{os: "windows", arch: []string{"arm64", "amd64"}},
	} {
		for _, arch := range target.arch {
			// NOTE: Each `MustExec` could be invoked on a separate goroutine to parallelize the binary compilation
			golang.MustExec(
				exec.WithTag("name", fmt.Sprintf("knita-%[1]s-%[2]s", target.os, arch)),
				exec.WithCommand("/bin/bash", "-c",
					fmt.Sprintf("GOOS=%[1]s GOARCH=%[2]s cd cmd/knita && go build -o ../../build/output/knita-%[1]s-%[2]s .", target.os, arch)))
		}
	}

	// Export the compiled Golang binaries back into the local source tree
	golang.MustExport("build/output/knita-*", "build/output/")
}
```

### Contributing


Thank you for considering contributing to Knita! All kinds of contributions are welcome, whether you're fixing bugs, adding new features, improving documentation, or helping others.

#### Guidelines
* **Code Style**: Please follow the coding style used in the project. Consistent code style helps maintain readability and quality.
* **Write Tests**: If you add a new feature, please write tests to cover it. If you're fixing a bug, consider adding a test that verifies the fix.
* **Documentation**: Update the documentation to reflect any changes you make. Good documentation helps others understand how to use and contribute to the project.
  ( **Pull Request Reviews**: Be open to feedback and revisions. Pull requests are a conversation, and improving the codebase is a collaborative effort.

#### Getting Help
If you need any help or have questions, feel free to open an issue or join the community discussions.
