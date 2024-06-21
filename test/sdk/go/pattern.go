package main

import (
	"bytes"
	"log"
	"os"

	"github.com/knita-io/knita/sdk/go/knita"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
	"github.com/knita-io/knita/sdk/go/knita/runtime/exec"
)

func main() {
	client := knita.MustNewClient()

	runtimesFactories := []func() *knita.Runtime{
		func() *knita.Runtime {
			return client.MustRuntime(
				runtime.WithType(runtime.TypeHost),
				runtime.WithTag("name", "host-test"))
		},
		func() *knita.Runtime {
			return client.MustRuntime(
				runtime.WithType(runtime.TypeDocker),
				runtime.WithImage("ubuntu:latest"),
				runtime.WithPullStrategy(runtime.PullStrategyNotExists),
				runtime.WithTag("name", "docker-test"))
		},
	}

	for _, factory := range runtimesFactories {
		func() {
			rt := factory()
			defer rt.MustClose()

			// Verify sysinfo is reported
			if rt.SysInfo() == nil {
				log.Fatalf("sysinfo not set")
			}
			if rt.SysInfo().Os == "" {
				log.Fatalf("sysinfo os not set")
			}
			if rt.SysInfo().Arch == "" {
				log.Fatalf("sysinfo arch not set")
			}
			if rt.SysInfo().TotalCpuCores <= 0 {
				log.Fatalf("sysinfo cpu cores not set")
			}
			if rt.SysInfo().TotalMemory <= 0 {
				log.Fatalf("sysinfo memory not set")
			}

			// Verify files can be imported
			expectedFilePath := "input/input.txt"
			buf, err := os.ReadFile(expectedFilePath)
			if err != nil {
				log.Fatalf("error reading input file: %v", err)
			}
			expectedContents := string(buf)
			rt.MustImport(expectedFilePath, "")
			rt.MustExec(
				exec.WithTag("name", "import-test"),
				exec.WithCommand("/bin/bash", "-c", `
				contents="$(cat `+expectedFilePath+`)"
				if [[ "$contents" != "`+expectedContents+`" ]]; then
					exit 1
				fi
			`))

			// Verify zero-byte files can be imported
			rt.MustImport("input/zero-bytes.txt", "")
			rt.MustExec(
				exec.WithTag("name", "zero-byte-import-test"),
				exec.WithCommand("/bin/bash", "-c", `stat input/zero-bytes.txt`))

			// Verify the remote work directory is reported correctly
			rt.MustExec(
				exec.WithTag("name", "work-directory-test"),
				exec.WithCommand("/bin/bash", "-c", `
				contents="$(cat `+rt.WorkDirectory(expectedFilePath)+`)"
				if [[ "$contents" != "`+expectedContents+`" ]]; then
					exit 1
				fi
			`))

			// Verify files can be exported
			expectedContents = "hello world\n"
			expectedFilePath = "output/host.txt"
			rt.MustExec(
				exec.WithTag("name", "export-test"),
				exec.WithCommand("/bin/bash", "-c", `
				mkdir output && echo -n '`+expectedContents+`' > `+expectedFilePath+`
			`))
			rt.MustExport(expectedFilePath, "")
			buf, err = os.ReadFile(expectedFilePath)
			if err != nil {
				log.Fatalf("error reading output file: %v", err)
			}
			if string(buf) != expectedContents {
				log.Fatalf("mismatched export contents")
			}

			// Verify stdout and stderr can be captured
			expectedOutput := "hello world\n"
			stdout := bytes.NewBuffer(make([]byte, 0))
			stderr := bytes.NewBuffer(make([]byte, 0))
			rt.MustExec(
				exec.WithTag("name", "io-test"),
				exec.WithStdout(stdout),
				exec.WithStderr(stderr),
				exec.WithCommand("/bin/bash", "-c", `
				echo -n "`+expectedOutput+`" | tee /dev/stderr
			`))
			if stdout.String() != expectedOutput {
				log.Fatalf("mismatched stdout output: %v", stdout.String())
			}
			if stderr.String() != expectedOutput {
				log.Fatalf("mismatched stderr output: %v", stderr.String())
			}
		}()
	}
}
