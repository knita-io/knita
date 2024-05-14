module github.com/knita-io/knita/sdk/go/knita

go 1.21.9

replace github.com/knita-io/knita => ../../../

require (
	github.com/knita-io/knita v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.63.2
)

require (
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
