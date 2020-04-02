package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"go.opentelemetry.io/otel/sdk/resource/resourcekeys"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
	exporter, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:55680"),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResourceAttributes(core.Key(resourcekeys.ServiceKeyName).String("go-service")),
	)
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}

func main() {
	initTracer()
	tr := global.Tracer("go-demo")

	s := &server{
		tracer: tr,
	}

	http.HandleFunc("/", s.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type server struct {
	tracer trace.Tracer
}

func (s *server) handler(w http.ResponseWriter, req *http.Request) {
	attrs, _, spanCtx := httptrace.Extract(req.Context(), req)

	// ctx, span := tracer.Start(
	ctx, span := s.tracer.Start(
		trace.ContextWithRemoteSpanContext(req.Context(), spanCtx),
		"hello",
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	response := "hello from go\n"
	if pyBody, err := s.fetchFromPythonService(ctx); err == nil {
		response += string(pyBody)
	} else {
		response += "error fetching from python"
	}

	_, _ = io.WriteString(w, response)
}

func (s *server) fetchFromPythonService(ctx context.Context) ([]byte, error) {
	client := http.DefaultClient
	var body []byte
	err := s.tracer.WithSpan(ctx, "fetch-from-python",
		func(ctx context.Context) error {
			req, err := http.NewRequest("GET", "http://localhost:8082/", nil)
			if err != nil {
				return err
			}

			ctx, req = httptrace.W3C(ctx, req)
			httptrace.Inject(ctx, req)

			fmt.Printf("Sending request...\n")
			res, err := client.Do(req)
			if err != nil {
				return err
			}
			body, err = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			return err
		})
	return body, err
}
