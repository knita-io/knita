module github.com/knita-io/knita/sdk/go/knita

go 1.21.9

replace github.com/knita-io/knita => ../../../

require (
	github.com/knita-io/knita v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.65.0
)

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240808171019-573a1156607a // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
