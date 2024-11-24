package tracing

import (
	"errors"

	"github.com/uber/jaeger-client-go/config"
)

// Init инициализирует трейсы
func Init(serviceName string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "localhost:6831",
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return errors.New("failed to init tracing")
	}

	return nil
}
