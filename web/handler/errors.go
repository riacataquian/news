package handler

// This file contains error handlers.

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"

	pb "github.com/riacataquian/news/protos/api"
)

// NotFound handles HTTP requests for missing or not found pages and resources.
func NotFound(_ context.Context) (proto.Message, *pb.HTTPError) {
	return nil, &pb.HTTPError{Code: http.StatusNotFound, Error: &pb.Error{Message: "page not found"}}
}
