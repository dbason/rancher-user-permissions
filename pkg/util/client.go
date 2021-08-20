package util

import (
	"fmt"
	"os"
	"time"

	managementv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	TimeoutFlagValue         time.Duration
	NamespaceFlagValue       string
	ContextOverrideFlagValue string
	K8sClient                client.Client
	RestConfig               *rest.Config
)

// ClientOptions can be passed to some of the functions in this package when
// creating clients and/or client configurations.
type ClientOptions struct {
	overrides *clientcmd.ConfigOverrides
}

type ClientOption func(*ClientOptions)

func (o *ClientOptions) Apply(opts ...ClientOption) {
	for _, op := range opts {
		op(o)
	}
}

// WithConfigOverrides allows overriding specific kubeconfig fields from the
// user's loaded kubeconfig.
func WithConfigOverrides(overrides *clientcmd.ConfigOverrides) ClientOption {
	return func(o *ClientOptions) {
		o.overrides = overrides
	}
}

// CreateClientOrDie constructs a new controller-runtime client, or exit
// with a fatal error if an error occurs.
func CreateClientOrDie(opts ...ClientOption) (*rest.Config, client.Client) {
	scheme := CreateScheme()
	clientConfig := LoadClientConfig(opts...)

	cli, err := client.New(clientConfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return clientConfig, cli
}

// LoadClientConfig loads the user's kubeconfig using the same logic as kubectl.
func LoadClientConfig(opts ...ClientOption) *rest.Config {
	options := ClientOptions{}
	options.Apply(opts...)

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig, err := rules.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	clientConfig, err := clientcmd.NewDefaultClientConfig(
		*kubeconfig, options.overrides).ClientConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return clientConfig
}

// CreateScheme creates a new scheme with the types necessary for opnictl.
func CreateScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(managementv3.AddToScheme(scheme))
	return scheme
}

// MaybeContextOverride will return 0 or 1 ClientOptions, depending on if the
// user provided a specific kubectl context or not.
func MaybeContextOverride() []ClientOption {
	if ContextOverrideFlagValue != "" {
		return []ClientOption{
			WithConfigOverrides(&clientcmd.ConfigOverrides{
				CurrentContext: NamespaceFlagValue,
			}),
		}
	}
	return []ClientOption{}
}

func LoadDefaultClientConfig() {
	RestConfig, K8sClient = CreateClientOrDie(MaybeContextOverride()...)
}
