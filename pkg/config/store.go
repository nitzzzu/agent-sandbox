package config

import (
	"context"
	"sync/atomic"

	"knative.dev/pkg/configmap"
)

type cfgKey struct{}

const configMapName = "config-templates"

// FromContext obtains a Config injected into the passed context.
func FromContext(ctx context.Context) *Config {
	return ctx.Value(cfgKey{}).(*Config)
}

// TemplatesContent loads/unloads our untyped configuration.
type TemplatesContent struct {
	*configmap.UntypedStore

	// current is the current Config.
	current atomic.Value
}

// NewStore creates a new configuration Store.
func NewStore(logger configmap.Logger, onAfterStore ...func(name string, value interface{})) *TemplatesContent {
	s := &TemplatesContent{}

	// Append an update function to run after a ConfigMap has updated to update the
	// current state of the Config.
	onAfterStore = append(onAfterStore, func(_ string, _ interface{}) {
		//c := &TemplatesContent{}
		// this allows dynamic updating of the config-network
		// this is necessary for not requiring activator restart for system-internal-tls in the future
		// however, current implementation is not there yet.
		// see https://github.com/knative/serving/issues/13754
		tpls := s.UntypedLoad(configMapName)
		s.current.Store(tpls)
	})
	s.UntypedStore = configmap.NewUntypedStore(
		"agent-sandbox",
		logger,
		configmap.Constructors{
			configMapName: Templates,
		},
		onAfterStore...,
	)
	return s
}

// ToContext stores the configuration Store in the passed context.
func (s *TemplatesContent) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cfgKey{}, s.current.Load())
}
