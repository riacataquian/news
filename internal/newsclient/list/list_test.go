package list

import (
	"testing"
)

func TestEncode(t *testing.T) {
	in := &Params{
		Query:    "some-query",
		From:     "some-from",
		To:       "some-to",
		Sources:  "some-source1,some-source2",
		Domains:  "some-domain1,some-domain2",
		SortBy:   Popularity,
		Language: "en",
		Page:     2,
		PageSize: 10,
	}
	want := "domains=some-domain1%2Csome-domain2&from=some-from&language=en&page=2&pageSize=10&q=some-query&sortBy=popularity&sources=some-source1%2Csome-source2&to=some-to"
	if got, err := in.Encode(); got != want {
		desc := "returns the correct query params given valid params"
		t.Errorf("%s: Encode(): want (%v, nil), got (%v, %v)", desc, want, got, err)
	}
}

func TestEncodeErrors(t *testing.T) {
	tests := []struct {
		desc    string
		in      *Params
		wantErr error
	}{
		{
			desc:    "pageSize exceeded the maxPageSize",
			in:      &Params{PageSize: 500, Language: "en"},
			wantErr: ErrInvalidPageSize,
		},
		{
			desc:    "no parameter is supplied",
			wantErr: ErrNoRequiredParams,
		},
		{
			desc:    "no required parameter is supplied",
			in:      &Params{PageSize: 500, Page: 2},
			wantErr: ErrNoRequiredParams,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := test.in.Encode()
			if err == nil {
				t.Fatalf("%s: Encode(): want (nil, %v), got (%v, %v)", test.desc, test.wantErr, got, err)
			}

			if err != test.wantErr {
				t.Errorf("%s: Encode(): want (nil, %v), got (%v, %v)", test.desc, test.wantErr, got, err)
			}
		})
	}
}
