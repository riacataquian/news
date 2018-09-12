// Package main is the entry point for News platform web interface.
package main // import "github.com/riacataquian/news"

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/riacataquian/news/internal/httperror"
	"github.com/riacataquian/news/internal/store"
	"github.com/riacataquian/news/web/handler"

	"github.com/gorilla/mux"
)

const (
	// DefaultErrStatusCode is the default status code for HTTP error responses.
	DefaultErrStatusCode = http.StatusInternalServerError
)

// main starts a web server and register routes and their matching handlers.
// It injects a context.Context argument for the route handlers to allow deadline and cancelation among HTTP requests.
// It also injects a data repository handler to be consumed by the HTTP handlers.
// Finally, it marshals successful and error JSON responses.
func main() {
	serve(context.Background(), store.New())
}

func serve(ctx context.Context, repo store.Store) {
	srv := mux.NewRouter().PathPrefix("/api").Subrouter()
	for _, route := range handler.Routes {
		srv.Handle(route.Path, middleware(ctx, repo, route.HandlerFunc))
	}

	if err := http.ListenAndServe(":8000", srv); err != nil {
		log.Fatalf("could not listen to port 8000: %v", err)
	}
}

// middleware transforms a handler.Func to http.HandlerFunc.
//
// The middleware does the repetitive yet necessary calculations for a handler:
// 1. Sets the response's content-type to application/json.
// 2. Sets the supplied status code in the response's header then finally encode the response for JSON rendering.
// 3. When an error is encountered, sets the proper response header, given an httperror or the DefaultErrStatusCode.
func middleware(ctx context.Context, repo store.Store, h handler.Func) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		resp, err := h(ctx, repo, r)
		if err == nil {
			w.WriteHeader(resp.Code)
			encode(w, resp)
			return
		}

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

// encode encodes `r` to `w` as JSON responses.
func encode(w io.Writer, r interface{}) {
	if err := json.NewEncoder(w).Encode(r); err != nil {
		log.Fatalf("error marshalling response: %v", err)
		os.Exit(1)
	}
}
