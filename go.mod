module github.com/dbason/rancher-user-permissions

go 1.16

replace k8s.io/client-go => k8s.io/client-go v0.21.2

require (
	github.com/rancher/rancher/pkg/apis v0.0.0-20210820000421-50c33b29b096
	github.com/spf13/cobra v1.1.3
	github.com/ttacon/chalk v0.0.0-20160626202418-22c06c80ed31
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.17.0
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.9.0-beta.0
)
