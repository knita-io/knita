FROM debian:bookworm
MAINTAINER tylerpowerx@gmail.com

RUN apt update
RUN apt install -y curl unzip python3 python3-pip python3-venv

RUN set -eux; \
	PROTOC_VERSION=26.1; \
	dpkgArch="$(dpkg --print-architecture)"; \
	case "$dpkgArch" in \
		amd64)  arch='x86_64' ;; \
		arm64) arch='aarch_64' ;; \
		*) echo >&2 "error: unsupported architecture: $dpkgArch"; exit 1 ;; \
	esac; \
	mkdir -p /tmp/protoc; \
	cd /tmp/protoc; \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-${arch}.zip; \
    unzip protoc-${PROTOC_VERSION}-linux-${arch}.zip; \
    cp -R bin/* /usr/local/bin/; \
    cp -R include/* /usr/local/include/; \
    rm -rf /tmp/protoc

RUN set -eux; \
    GOLANG_VERSION=1.22.4; \
    dpkgArch="$(dpkg --print-architecture)"; \
    case "$dpkgArch" in \
        amd64)  arch='amd64' ;; \
        arm64) arch='arm64' ;; \
        *) echo >&2 "error: unsupported architecture: $dpkgArch"; exit 1 ;; \
    esac; \
    curl -LO https://go.dev/dl/go${GOLANG_VERSION}.linux-${arch}.tar.gz; \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-${arch}.tar.gz; \
    rm go${GOLANG_VERSION}.linux-${arch}.tar.gz
ENV PATH="${PATH}:/go/bin:/usr/local/go/bin"
ENV GOPATH=/go

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
RUN go install golang.org/x/tools/cmd/goimports@latest

RUN python3 -m venv /tmp/python3
ENV PATH="/tmp/python3/bin:$PATH"
RUN pip3 install build twine 'grpcio-tools==1.63.0'
