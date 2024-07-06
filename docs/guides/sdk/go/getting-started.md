# Getting Started with the Knita SDK for Go

In this guide we'll create a simple "hello world" build using the Knita SDK for Go. In order to run this, you will need
Docker installed.

1. Download the latest Knita CLI from the [release page](https://github.com/knita-io/knita/releases) and make sure it's
   in your path

2. Create a new `build` directory to hold your Knita pattern, and initialize a new Go module

    ```bash
    mkdir build
    cd build
    go mod init github.com/myorg/myproject/build
    ```

3. Create a new `build/pattern.go` Go file and copy in the example below:

    ```go
   package main
   
   import (
       "github.com/knita-io/knita/sdk/go/knita"
       "github.com/knita-io/knita/sdk/go/knita/runtime"
       "github.com/knita-io/knita/sdk/go/knita/runtime/exec"
   )
   
   func main() {
       client := knita.MustNewClient()
   
       golang := client.MustRuntime(
           runtime.WithTag(knita.NameTag, "example"),
           runtime.WithType(runtime.TypeDocker),
           runtime.WithImage("alpine:latest"))
       defer golang.MustClose()
   
       golang.MustExec(
           exec.WithTag(knita.NameTag, "hello-world"),
           exec.WithCommand("/bin/sh", "-c", `
               echo 'hello world'`))
   }
    ```

4. Update your go.mod and fetch the Knita SDK:
    ```bash
    go mod tidy
    ```

5. Create a new `build/pattern.sh` shell script to build and execute the pattern:

   Go is a compiled language, so the pattern must be built before it can be executed. An easy way to do this is to
   create a wrapper script that can be passed to the Knita CLI. In other SDK languages like Python, this step is
   not necessary.

    ```bash
    #!/bin/bash
    pushd build
      mkdir output
      go build -o output/knita-pattern .
    popd
    ./build/output/knita-pattern
    ```

    ```bash
    chmod +x build/pattern.sh
    ```

6. Run your new pattern using the Knita CLI:
    ```bash
    knita build ./build/pattern.sh
    ```

    You should see the following output:

    ```
    > knita build ./build/pattern.sh
    example: finished
     ✓ hello-world (0s)
    
    Build log available at: /var/folders/bj/8x5csq0s15sgfpk2sq8ld7kw0000gn/T/knita/knita-build-20240513T195847Z.log
    ```

    And if you have a look at the build log:

    ```
    ~/Development/knita/build ~/Development/knita
    mkdir: output: File exists
    ~/Development/knita
    Pulling Docker image...
    Docker pull strategy is "default", image exists in cache and is not latest; "docker.io/library/alpine:latest" will not be pulled
    Executing command: /bin/bash [-c echo 'hello world']
    hello world
    ```
   
