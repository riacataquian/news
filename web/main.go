// Package main is the entry point for News platform web interface.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/riacataquian/news/internal/httperror"
	"github.com/riacataquian/news/web/handler"

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
	ctx := context.Background()
	serve(ctx)
}

func serve(ctx context.Context) {
	srv := mux.NewRouter()
	for _, route := range handler.Routes {
		srv.Handle(route.Path, middleware(ctx, route.HandlerFunc))
	}

	if err := http.ListenAndServe(":8000", srv); err != nil {
		log.Fatalf("could not listen to port 8000: %v", err)
	}
}

func middleware(ctx context.Context, h handler.Func) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		resp, err := h(ctx, r)
		if err == nil {
			// TODO: Work around hardcoded http.StatusOK.
			w.WriteHeader(http.StatusOK)
			encode(w, resp)
		} else {
			var code int
			if v, ok := err.(*httperror.HTTPError); ok {
				code = v.Code
			} else {
				code = DefaultErrStatusCode
			}

			w.WriteHeader(code)
			encode(w, err)
		}
	}
}

// encode encodes `r` as JSON responses.
func encode(w io.Writer, r interface{}) {
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Fatalf("error marshalling response: %v", err)
		os.Exit(1)
	}
}
