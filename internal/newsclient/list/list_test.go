package list

import (
	"testing"
)

func TestEncode(t *testing.T) {
	in := &Params{
		Query:    "some-query",
		Sources:  "some-source1,some-source2",
		Domains:  "some-domain1,some-domain2",
		SortBy:   Popularity,
		Language: "en",
	}
	want := "domains=some-domain1%2Csome-domain2&language=en&q=some-query&sortBy=popularity&sources=some-source1%2Csome-source2"
	if got, err := in.Encode(); got != want {
		desc := "returns the correct query params given valid params"
		t.Errorf("%s: Encode(): want (%v, nil), got (%v, %v)", desc, want, got, err)
	}
}

func TestEncodeErrors(t *testing.T) {
	tests := []struct {
		desc string
		in   *Params
	}{
		{
			desc: "pageSize exceeded the maxPageSize",
			in:   &Params{PageSize: 500, Language: "en"},
		},
		{
			desc: "no required parameter is supplied",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if got, err := test.in.Encode(); err == nil {
				t.Errorf("%s: Encode(): want (nil, error), got (%v, %v)", test.desc, got, err)
			}
		})
	}
}
