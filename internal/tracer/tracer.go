package tracer

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"log"
)

func MustSetup(ctx context.Context, serviceName string) {
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

	go func() {
		<-ctx.Done()
		if err = closer.Close(); err != nil {
			log.Printf("error closing tracer: %s\n", err)
		}
	}()

	opentracing.SetGlobalTracer(tracer)
}
