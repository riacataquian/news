// Package handler contains routes and HTTP handlers.
package handler

import (
	"context"

	"github.com/golang/protobuf/proto"

	pb "github.com/riacataquian/news/protos/api"
)

// Func describes a function that handles HTTP requests and responses.
type Func func(context.Context) (proto.Message, *pb.HTTPError)

// Routes is the lookup table for URL paths and their matching handlers.
var Routes = []struct {
	Path        string
	HandlerFunc Func
}{
	{"/api/news", News},
	{"/{*}", NotFound},
}
