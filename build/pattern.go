package main

import (
	"fmt"

	"github.com/knita-io/knita/sdk/go/knita"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
	"github.com/knita-io/knita/sdk/go/knita/runtime/exec"
)

var targets = []struct {
	os   string
	arch []string
}{
	{os: "darwin", arch: []string{"arm64", "amd64"}},
	{os: "linux", arch: []string{"arm", "arm64", "amd64"}},
	{os: "windows", arch: []string{"arm64", "amd64"}},
}

func main() {
	client := knita.MustNewClient()
	docker(client)
	generate(client)
	unit(client)
	build(client)
}

func docker(s *knita.Client) {
	host := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "docker"),
		runtime.WithType(runtime.TypeHost))
	defer host.MustClose()

	host.MustImport("build/docker/*", "")
	host.MustExec(
		exec.WithTag(knita.NameTag, "build"),
		exec.WithCommand("/bin/bash", "-c", `
			cd build/docker && \
			docker build -t knita/build:latest .`),
	)
}

func generate(s *knita.Client) {
	golang := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "generate"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer golang.MustClose()

	golang.MustImport("api/**/*.proto", "")
	golang.MustExec(
		exec.WithTag(knita.NameTag, "protobuf"),
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
	golang.MustExport("api/**/*.pb.go", "")
}

func unit(s *knita.Client) {
	golang := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "unit"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer golang.MustClose()

	golang.MustImport(".", ".")
	golang.MustExec(
		exec.WithTag(knita.NameTag, "go"),
		exec.WithCommand("/bin/bash", "-c", "go test ./..."))
}

func build(s *knita.Client) {
	golang := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "build"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage("knita/build:latest"),
		runtime.WithPullStrategy(runtime.PullStrategyNever))
	defer golang.MustClose()

	golang.MustImport(".", ".")
	for _, target := range targets {
		for _, arch := range target.arch {
			golang.MustExec(
				exec.WithTag(knita.NameTag, fmt.Sprintf("knita-%[1]s-%[2]s", target.os, arch)),
				exec.WithCommand("/bin/bash", "-c",
					fmt.Sprintf("GOOS=%[1]s GOARCH=%[2]s cd cmd/knita && go build -o ../../build/output/cli/knita-%[1]s-%[2]s .", target.os, arch)))
		}
	}
	golang.MustExport("build/output/cli/knita-*", "build/output/cli/")
}
