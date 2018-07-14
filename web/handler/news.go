package handler

// This file contains HTTP handlers for news endpoint.

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"

	pb "github.com/riacataquian/news/protos/api"
)

// News ...
func News(_ context.Context) (proto.Message, *pb.HTTPError) {
	return nil, &pb.HTTPError{Code: http.StatusInternalServerError, Error: &pb.Error{Message: "news error"}}
}
