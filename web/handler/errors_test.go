package handler

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestNotFound(t *testing.T) {
	x, err := NotFound(ctx)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	t.Fatalf("res %v", x)
}
