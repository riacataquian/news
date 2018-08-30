package newsclient

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/api/news"
)

type fakeParams struct {
	lang    string
	wantErr bool
}

func (fp fakeParams) Encode() (string, error) {
	if fp.wantErr {
		return "", errors.New("encoding error")
	}
	return "language=" + fp.lang, nil
}

func TestNewFromContext(t *testing.T) {
	ctx := context.Background()
	testSE := ServiceEndpoint{
		RequestURL: "http://some-request-url",
		DocsURL:    "some-docs-url",
	}
	got := NewFromContext(ctx, testSE)
	want := &Client{ctx: ctx, ServiceEndpoint: testSE}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a new newsclient"
		t.Errorf("%s: NewFromContext(): Diff (-got +want)\n%s", desc, diff)
	}
}

func TestGet(t *testing.T) {
	server := setupStubServer(t, true)
	defer server.Close()

	client := setupFakeClient(server.URL)

	want := fakeResponse
	testAuthKey := "test-auth-key"
	params := fakeParams{lang: "en"}

	got, err := client.Get(context.Background(), testAuthKey, params)
	if err != nil {
		t.Fatalf("Get(%s, %v): want (%v, nil), got (%v, %v)", testAuthKey, params, want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a news.Response and nil error"
		t.Errorf("%s: Get(%s, %v) diff: (-got +want)\n%s", desc, testAuthKey, params, diff)
	}
}

func TestGetErrors(t *testing.T) {
	tests := []struct {
		desc          string
		isServerValid bool
		params        fakeParams
	}{
		{
			desc:   "returns an error when server errored",
			params: fakeParams{lang: "en"},
		},
		{
			desc:          "returns an error when params errored",
			isServerValid: true,
			params:        fakeParams{lang: "en", wantErr: true},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			server := setupStubServer(t, test.isServerValid)
			defer server.Close()

			client := setupFakeClient(server.URL)
			testAuthKey := "test-auth-key"
			if got, err := client.Get(context.Background(), testAuthKey, test.params); err == nil {
				t.Errorf("%s: Get(%s, %v) want (nil, error), got (%v, %v)", test.desc, testAuthKey, test.params, got, err)
			}
		})
	}
}

func TestDispatchReq(t *testing.T) {
	want := fakeResponse
	server := setupStubServer(t, true)
	defer server.Close()

	r, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("dispatchReq(_): error creating a new request: %v", err)
	}

	got, err := dispatchReq(r)
	if err != nil {
		t.Fatalf("dispatchReq(_): want (%v, nil), got (%v, %v)", want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a news.Response and nil error"
		t.Errorf("%s: dispatchReq(_) diff: (-got +want)\n%s", desc, diff)
	}
}

func TestDispatchReqErrors(t *testing.T) {
	want := &news.ErrorResponse{
		Code:    "500",
		Message: "some error",
	}

	server := setupStubServer(t, false)
	defer server.Close()

	r, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("dispatchReq(_): error creating a new request: %v", err)
	}

	got, err := dispatchReq(r)
	if err == nil {
		t.Fatalf("dispatchReq(_): want (nil, error), got (%v, %v)", got, err)
	}

	if diff := pretty.Compare(err, want); diff != "" {
		desc := "returns a news.ErrorResponse when error is encountered"
		t.Errorf("%s: dispatchReq(_) diff: (-got +want)\n%s", desc, diff)
	}
}
