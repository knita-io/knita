module github.com/knita-io/knita/test/sdk/go

go 1.22.2

replace github.com/knita-io/knita/sdk/go/knita => ./../../../sdk/go/knita

replace github.com/knita-io/knita => ./../../../

require github.com/knita-io/knita/sdk/go/knita v0.0.0-00010101000000-000000000000

require (
	github.com/knita-io/knita v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240808171019-573a1156607a // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
