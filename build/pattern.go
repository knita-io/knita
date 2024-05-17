package main

import (
	"fmt"
	"log"
	"os"
	stdruntime "runtime"
	"sync"

	"github.com/knita-io/knita/sdk/go/knita"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
	"github.com/knita-io/knita/sdk/go/knita/runtime/exec"
)

func main() {
	client := knita.MustNewClient()
	docker(client)
	protobuf(client)
	unit(client)
	build(client)
	testSDK(client)
	publishSDK(client)
}

func docker(s *knita.Client) {
	host := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "docker"),
		runtime.WithType(runtime.TypeHost))
	defer host.MustClose()

	host.MustImport("build/docker/*", "")
	host.MustExec(
		exec.WithTag(knita.NameTag, "knita/build"),
		exec.WithCommand("/bin/bash", "-c", `
			cd build/docker
			docker build -t knita/build:latest .`),
	)
}

func protobuf(s *knita.Client) {
	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "protobuf"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer container.MustClose()

	container.MustImport("api/**/*.proto", "")

	container.MustExec(
		exec.WithTag(knita.NameTag, "python"),
		exec.WithCommand("/bin/bash", "-c", `
			python \
			-m grpc_tools.protoc \
			-I api \
			--python_out=api \
			--pyi_out=api \
			--grpc_python_out=api \
			executor/v1/executor.proto \
			director/v1/director.proto

			sed -i -e 's/from director.v1/from ./g' api/*/v1/*.py*
			sed -i -e 's/from executor.v1/from ./g' api/*/v1/*.py*`))
	container.MustExport("api/**/*.py*", "sdk/python/knita/")

	container.MustExec(
		exec.WithTag(knita.NameTag, "go"),
		exec.WithCommand("/bin/bash", "-c", `
			protoc \
			--proto_path=api \
			--go_out=api \
			--go_opt=paths=source_relative \
			--go-grpc_out=api \
			--go-grpc_opt=paths=source_relative \
			broker/v1/broker.proto \
			executor/v1/executor.proto \
			director/v1/director.proto \
			observer/v1/observer.proto`))
	container.MustExport("api/**/*.pb.go", "")
}

func unit(s *knita.Client) {
	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "unit-test"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer container.MustClose()

	container.MustImport(".", ".")
	container.MustExec(
		exec.WithTag(knita.NameTag, "go"),
		exec.WithCommand("/bin/bash", "-c", "go test ./..."))
}

func build(s *knita.Client) {
	var targets = []struct {
		os   string
		arch []string
	}{
		{os: "darwin", arch: []string{"arm64", "amd64"}},
		{os: "linux", arch: []string{"arm", "arm64", "amd64"}},
		{os: "windows", arch: []string{"arm64", "amd64"}},
	}

	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "build"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer container.MustClose()

	container.MustImport(".", ".")
	wg := sync.WaitGroup{}
	for _, target := range targets {
		for _, arch := range target.arch {
			wg.Add(1)
			go func(os string, arch string) {
				defer wg.Done()
				container.MustExec(
					exec.WithTag(knita.NameTag, fmt.Sprintf("knita-%[1]s-%[2]s", os, arch)),
					exec.WithCommand("/bin/bash", "-c",
						fmt.Sprintf("cd cmd/knita && env GOOS=%[1]s GOARCH=%[2]s go build -o ../../build/output/cli/knita-%[1]s-%[2]s .", os, arch)))
			}(target.os, arch)
		}
	}
	wg.Wait()
	container.MustExport("build/output/cli/knita-*", "build/output/cli/")
}

// testSDK runs e2e tests for the Knita SDKs. This is a little meta as we're invoking `knita build`
// from inside an existing `knita build`, but targeting a per-sdk test pattern file.
func testSDK(s *knita.Client) {
	host := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "sdk-test"),
		runtime.WithType(runtime.TypeHost))
	defer host.MustClose()

	host.MustImport(fmt.Sprintf("build/output/cli/knita-%s-%s", stdruntime.GOOS, stdruntime.GOARCH), "knita")
	host.MustImport("go.mod", "")
	host.MustImport("api", "")
	host.MustImport("test", "")
	host.MustImport("sdk", "")

	host.MustExec(
		exec.WithTag(knita.NameTag, "python"),
		exec.WithCommand("/bin/bash", "-c", `
			python3 -m venv python-sdk-test
			source python-sdk-test/bin/activate
			python3 -m pip install sdk/python
			export PATH=$PATH:$(pwd)
			cd test/sdk/python
			knita build ./pattern.py`))

	host.MustExec(
		exec.WithTag(knita.NameTag, "go"),
		exec.WithCommand("/bin/bash", "-c", `
			export PATH=$PATH:$(pwd)
			cd test/sdk/go
			knita build ./pattern.sh`))
}

func publishSDK(s *knita.Client) {
	if os.Getenv("KNITA_PUBLISH_SDK") == "" {
		log.Printf("\n\nKNITA_PUBLISH_SDK is unset; Skipping SDK publishing\n\n")
		return
	}

	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "sdk-publish"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer container.MustClose()

	container.MustImport("sdk/python", "")
	container.MustExec(
		exec.WithTag(knita.NameTag, "python"),
		exec.WithEnv("TWINE_PASSWORD="+os.Getenv("TWINE_PASSWORD")),
		exec.WithCommand("/bin/bash", "-c", `
			cd sdk/python
			python3 -m build
			twine upload --non-interactive -u __token__ -p $TWINE_PASSWORD dist/*`))
}
