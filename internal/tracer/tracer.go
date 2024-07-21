package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"log"
)

func MustSetup(serviceName string) io.Closer {
	cfg, err := traceconfig.FromEnv()
	if err != nil {
		panic(err)
	}
	cfg.Sampler = &traceconfig.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	cfg.ServiceName = serviceName

	tracer, closer, err := cfg.NewTracer(traceconfig.Logger(jaeger.StdLogger), traceconfig.Metrics(prometheus.New()))
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger %s", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}
