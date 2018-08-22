package auth

import (
	"os"
	"testing"
)

func TestLookupAndSetAuth(t *testing.T) {
	tests := []struct {
		desc    string
		wantKey string
		wantErr error
	}{
		{
			desc:    "returns API_KEY env var if present",
			wantKey: "test-api-key",
		},
		{
			desc:    "returns ErrMissingAPIKey if API_KEY env var is not present",
			wantErr: ErrMissingAPIKey,
		},
	}

	for _, test := range tests {
		os.Setenv("API_KEY", test.wantKey)
		defer func() {
			os.Clearenv()
		}()

		got, err := LookupAndSetAuth()
		if test.wantErr != nil && err != nil {
			if test.wantErr != err {
				t.Fatalf("%s: LookupAndSetAuth() = (_, %v), got (_, %v)", test.desc, test.wantErr, err)
			}
		}

		if test.wantKey != got {
			t.Errorf("%s: LookupAndSetAuth() = (%v, _), got (%v, _)", test.desc, test.wantKey, got)
		}
	}
}
