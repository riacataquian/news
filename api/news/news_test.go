package news

import (
	"testing"
)

func TestError(t *testing.T) {
	want := "some error message"
	err := &ErrorResponse{
		Code:    "400",
		Message: want,
	}

	if got := err.Error(); got != want {
		t.Errorf("Error() = %s, want %s", got, want)
	}
}
