module github.com/knita-io/knita/build

go 1.21.9

replace github.com/knita-io/knita/sdk/go/knita => ./../sdk/go/knita

replace github.com/knita-io/knita => ./../

require (
	github.com/knita-io/knita v0.0.0-00010101000000-000000000000
	github.com/knita-io/knita/sdk/go/knita v0.0.0-00010101000000-000000000000
)

require (
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
