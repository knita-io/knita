FROM golang:1.21
MAINTAINER tylerpowerx@gmail.com

RUN apt-get update
RUN apt-get install -y curl unzip python3 python3-pip python3-venv

RUN PROTOC_VERSION=26.1 && \
    mkdir -p /tmp/protoc && \
    cd /tmp/protoc && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    cp -R bin/* /usr/local/bin/ && \
    cp -R include/* /usr/local/include/ && \
    rm -rf /tmp/protoc

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
RUN go install golang.org/x/tools/cmd/goimports@latest

RUN python3 -m venv /tmp/python3
ENV PATH="/tmp/python3/bin:$PATH"
RUN pip3 install build twine 'grpcio-tools==1.63.0'
