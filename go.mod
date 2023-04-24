module minik8s.io

go 1.20

require minik8s.io/cri-api v1.0.0

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace minik8s.io/cri-api => ./staging/cri-api
