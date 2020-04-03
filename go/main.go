package main

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/plugin/httptrace"
)

func main() {
	s := &server{}

	var mux http.ServeMux
	mux.Handle("/", http.HandlerFunc(s.handler))
	check(http.ListenAndServe(":8080", &mux))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type server struct{}

func (s *server) handler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	response := "hello from go\n"
	if pyBody, err := s.fetchFromPythonService(ctx); err == nil {
		response += string(pyBody)
	} else {
		response += "error fetching from python"
	}

	_, _ = io.WriteString(w, response)
}

func (s *server) fetchFromPythonService(ctx context.Context) ([]byte, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	var body []byte

	req, err := http.NewRequest("GET", "http://localhost:8082/", nil)
	if err != nil {
		return body, err
	}

	ctx, req = httptrace.W3C(ctx, req)
	httptrace.Inject(ctx, req)

	res, err := client.Do(req)
	if err != nil {
		return body, err
	}
	body, err = ioutil.ReadAll(res.Body)
	err = res.Body.Close()

	return body, err
}
