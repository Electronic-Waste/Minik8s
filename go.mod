module minik8s.io

go 1.20

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect

	k8s.io/cri-api/pkg/apis/runtime/v1
)

replace (
	k8s.io/cri-api => ./staging/cri-api
)