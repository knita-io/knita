package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	stdexec "os/exec"
	"regexp"
	stdruntime "runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/knita-io/knita/sdk/go/knita"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
	"github.com/knita-io/knita/sdk/go/knita/runtime/exec"
)

func main() {
	client := knita.MustNewClient()
	knitaVersion := mustGetKnitaVersion()
	builderDockerImage := dockerImage(client)
	protobuf(client, builderDockerImage)
	unit(client, builderDockerImage)
	build(client, builderDockerImage, knitaVersion)
	testSDK(client)
	publishSDK(client, builderDockerImage, knitaVersion)
}

// dockerImage builds the knita/build Docker image that is used by subsequent build targets.
// This image is versioned based on the content of the Dockerfile, and published to a publicly
// readable repository. If you're building Knita and need to change this Dockerfile, you will
// need write permissions to the repo. Open a GitHub issue to discuss.
func dockerImage(s *knita.Client) string {
	fingerprint := mustFingerprint("build/docker/Dockerfile")
	builderDockerImage := fmt.Sprintf("ghcr.io/knita-io/knita/build:%s", fingerprint)

	host := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "docker"),
		runtime.WithType(runtime.TypeHost))
	defer host.MustClose()

	password := os.Getenv("KNITA_BUILD_DOCKER_PASSWORD")
	if password == "" {
		log.Printf("KNITA_BUILD_DOCKER_PASSWORD is not set, build will fail if %s "+
			"does not already exist in public registry\n", builderDockerImage)
	}

	host.MustImport("build/docker/*", "")
	host.MustExec(
		exec.WithTag(knita.NameTag, "knita/build"),
		exec.WithEnv("DOCKER_PASSWORD="+password),
		exec.WithCommand("/bin/bash", "-c", `
			export DOCKER_CLI_EXPERIMENTAL=enabled
			if ! docker manifest inspect `+builderDockerImage+` > /dev/null; then
				echo $DOCKER_PASSWORD | docker login ghcr.io -u USERNAME --password-stdin
				cd build/docker
				docker buildx build --push --tag `+builderDockerImage+` --platform linux/amd64,linux/arm64 .
			fi`),
	)

	return builderDockerImage
}

// protobuf generates protobuf bindings for the languages used in the Knita repo.
func protobuf(s *knita.Client, builderDockerImage string) {
	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "protobuf"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage(builderDockerImage))
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

// unit executes unit tests for the Knita repo.
func unit(s *knita.Client, builderDockerImage string) {
	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "unit-test"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage(builderDockerImage))
	defer container.MustClose()

	container.MustImport(".", ".")
	container.MustExec(
		exec.WithTag(knita.NameTag, "go"),
		exec.WithCommand("/bin/bash", "-c", "go test ./..."))
}

// build compiles the various binaries in the Knita repo, and saves them to ./build/output/
func build(s *knita.Client, builderDockerImage string, knitaVersion knitaVersion) {
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
		runtime.WithImage(builderDockerImage))
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
					exec.WithEnv(fmt.Sprintf("LDFLAGS=-X github.com/knita-io/knita/internal/version.Version=%s", knitaVersion)),
					exec.WithCommand("/bin/bash", "-c",
						fmt.Sprintf("cd cmd/knita && env GOOS=%[1]s GOARCH=%[2]s go build -ldflags \"$LDFLAGS\" -o ../../build/output/cli/knita-%[1]s-%[2]s .", os, arch)))
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

// publishSDK builds and publishes Knita SDK to relevant package repositories when running against a release tag.
// Requires the KNITA_BUILD_TWINE_PASSWORD env var be set for pushing the Python SDK to Pypi.
func publishSDK(s *knita.Client, builderDockerImage string, knitaVersion knitaVersion) {
	// See Python version string spec: https://peps.python.org/pep-0440/
	if !knitaVersion.IsPublic() {
		log.Printf("Build is not a release build, skipping SDK publishing: %s\n", knitaVersion)
		return
	}

	password := os.Getenv("KNITA_BUILD_TWINE_PASSWORD")
	if password == "" {
		log.Fatal("KNITA_BUILD_TWINE_PASSWORD must be set when publishing SDKs")
	}

	container := s.MustRuntime(
		runtime.WithTag(knita.NameTag, "sdk-publish"),
		runtime.WithType(runtime.TypeDocker),
		runtime.WithImage(builderDockerImage))
	defer container.MustClose()

	container.MustImport("sdk/python", "")
	container.MustExec(
		exec.WithTag(knita.NameTag, "python"),
		exec.WithEnv("TWINE_PASSWORD="+os.Getenv("KNITA_BUILD_TWINE_PASSWORD")),
		exec.WithCommand("/bin/bash", "-c", `
			cd sdk/python
			sed -i -e 's/version = "0.0.0"/version = "`+knitaVersion.String()+`"/g' pyproject.toml
			python3 -m build
			twine upload --non-interactive -u __token__ -p $TWINE_PASSWORD dist/*`))
}

// mustGetKnitaVersion returns the Knita software version as derived from Git. Exits the process on error.
func mustGetKnitaVersion() knitaVersion {
	version, err := getKnitaVersion()
	if err != nil {
		log.Fatal(err)
	}
	return version
}

// getKnitaVersion returns the Knita software version as derived from Git.
func getKnitaVersion() (knitaVersion, error) {
	describe := bytes.NewBuffer(make([]byte, 0))
	cmd := stdexec.Command("git", "describe", "--long", "--tags", "--always")
	cmd.Stdout = describe
	err := cmd.Run()
	if err != nil {
		return knitaVersion{}, fmt.Errorf("error running git describe: %w", err)
	}
	regex := regexp.MustCompile("(v?[0-9+]+\\.[0-9+]+\\.[0-9+]+)-([0-9]+)-(.*)$")
	matches := regex.FindAllStringSubmatch(strings.Trim(describe.String(), "\n"), -1)
	if len(matches) == 0 || len(matches[0]) != 4 {
		return knitaVersion{}, fmt.Errorf("error unexpected number of matches in git describe regex: %v", matches)
	}
	semver := matches[0][1]
	tagDistance := matches[0][2]
	shortSHA := matches[0][3]
	distance, err := strconv.ParseInt(tagDistance, 10, 64)
	if err != nil {
		return knitaVersion{}, fmt.Errorf("error parsing tag distance: %w", err)
	}
	return knitaVersion{semver: semver, tagDistance: int(distance), shortSHA: shortSHA}, nil
}

// mustFingerprint calculates a SHA256 hash of the contents of files.
// Same files and contents in, same hash out. Useful for content addressable versioning.
// Exits the process on error.
func mustFingerprint(files ...string) string {
	h := sha256.New()
	for _, path := range files {
		func() {
			f, err := os.Open(path)
			if err != nil {
				log.Fatalf("error opening file %s: %v", path, err)
			}
			defer f.Close()
			if _, err := io.Copy(h, f); err != nil {
				log.Fatalf("error reading file %s: %v", path, err)
			}
		}()
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

type knitaVersion struct {
	semver      string
	tagDistance int
	shortSHA    string
}

func (v knitaVersion) IsPublic() bool {
	return v.tagDistance == 0
}

func (v knitaVersion) String() string {
	if !v.IsPublic() {
		return fmt.Sprintf("%s-%d-%s", v.semver, v.tagDistance, v.shortSHA)
	}
	return v.semver
}
