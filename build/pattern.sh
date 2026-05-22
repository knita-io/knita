#!/bin/bash

# GOTOOLCHAIN=auto lets `go` auto-download the toolchain pinned in go.mod (currently 1.23.0)
# even when the locally-installed Go is older.
export GOTOOLCHAIN=auto

cd build && go build -o output/knita-pattern . && cd .. && ./build/output/knita-pattern
