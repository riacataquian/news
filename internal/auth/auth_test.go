package auth

import (
	"os"
	"testing"
)

func TestLookupAPIAuthKey(t *testing.T) {
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

		got, err := LookupAPIAuthKey()
		if test.wantErr != nil && err != nil {
			if test.wantErr != err {
				t.Fatalf("%s: LookupAPIAuthKey() = (_, %v), got (_, %v)", test.desc, test.wantErr, err)
			}
		}

		if test.wantKey != got {
			t.Errorf("%s: LookupAPIAuthKey() = (%v, _), got (%v, _)", test.desc, test.wantKey, got)
		}
	}
}
