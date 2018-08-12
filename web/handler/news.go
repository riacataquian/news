package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/list"

	"github.com/gorilla/schema"
)

// This file contains handlers for news endpoint.

var client newsclient.Client

// List is the HTTP handler for news requests.
// See https://newsapi.org/docs/endpoints/everything for the official documentation.
func List(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = list.NewClient()
	params := new(list.Params)
	err := schema.NewDecoder().Decode(params, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	// Requests to external services should timeout for 5 seconds.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := client.NewGetRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("establishing a new GET request: %v", err)
	}

	key := lookupAPIKey()
	client.AuthorizeReq(req, key)

	return fetch(req, client, params)
}

// TopHeadlines is the HTTP handler for top-headlines news requests.
// See https://newsapi.org/docs/endpoints/top-headlines for the official documentation.
// func TopHeadlines(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
// 	r.ParseForm()

// 	client = headlines.NewClient()
// 	dst := new(headlines.Params)
// 	err := schema.NewDecoder().Decode(dst, r.Form)
// 	if err != nil {
// 		return nil, fmt.Errorf("error decoding params: %v", err)
// 	}

// 	return fetch(ctx, r, client, dst)
// }

// fetch performs the request to the client given params.
func fetch(req *http.Request, client newsclient.Client, params newsclient.Params) (*SuccessResponse, error) {
	news, err := client.Get(req, params)
	if err != nil {
		return nil, err
	}
	resp := &SuccessResponse{Data: news}

	qp := req.URL.Query() // nani dafuq? is it because the client encodes the params to the request URL?
	ps := qp.Get("pageSize")
	if ps != "" {
		count, err := strconv.Atoi(ps)
		if err != nil {
			return nil, fmt.Errorf("invalid page size value: %v", ps)
		}
		resp.Count = count
	}

	p := qp.Get("page")
	if p != "" {
		page, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid page value: %v", p)
		}
		resp.Page = page
	}

	resp.Code = http.StatusOK
	resp.RequestURL = req.URL.String()
	return resp, nil
}

// lookupAPIKey sets the env variable API_KEY in the supplied request.
func lookupAPIKey() string {
	k, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatal("missing API_KEY set as environment variable")
		os.Exit(1)
	}
	return k
}
