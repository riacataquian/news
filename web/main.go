// Package main is the entry point for News platform web interface.
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/riacataquian/news/web/handler"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
)

const (
	// DefaultErrStatusCode is the default status code for HTTP error responses.
	DefaultErrStatusCode = http.StatusInternalServerError
)

// main starts a web server and register routes and their matching handlers.
// It injects a context.Context argument for the route handlers to allow deadline and cancelation among HTTP requests.
// Finally, it marshals successful and error JSON responses.
func main() {
	srv := mux.NewRouter()
	ctx := context.Background()

	for _, route := range handler.Routes {
		srv.Handle(route.Path, middleware(ctx, route.HandlerFunc))
	}

	if err := http.ListenAndServe(":8000", srv); err != nil {
		log.Fatalf("could not listen to port 8000: %v", err)
	}
}

func middleware(ctx context.Context, h handler.Func) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		resp, err := h(ctx)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			encode(w, resp)
		} else {
			var code int
			if err.GetCode() == 0 {
				code = DefaultErrStatusCode
			} else {
				code = int(err.GetCode())
			}

			w.WriteHeader(code)
			encode(w, err)
		}
	}
}

// encode encodes proto messages as binary-encoded responses.
func encode(w io.Writer, res proto.Message) {
	b, err := proto.Marshal(res)
	if err != nil {
		log.Fatalf("error marshalling response: %v", err)
	}

	_, err = w.Write(b)
	if err != nil {
		log.Fatalf("error writing response: %v", err)
		os.Exit(1)
	}
}
