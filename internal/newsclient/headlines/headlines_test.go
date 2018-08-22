package headlines

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		desc string
		in   *Params
		want string
	}{
		{
			desc: "returns the encoded params",
			in:   &Params{Query: "bitcoin", Country: "us", Category: "tech", Page: 1, PageSize: 10},
			want: "category=tech&country=us&page=1&pageSize=10&q=bitcoin",
		},
		{
			desc: "returns the correct query params given valid params",
			in:   &Params{Sources: "some-source", Query: "bitcoin", PageSize: 50, Page: 2},
			want: "page=2&pageSize=50&q=bitcoin&sources=some-source",
		}}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if got, err := test.in.Encode(); got != test.want {
				t.Errorf("%s: Encode(): want (%v, nil), got (%v, %v)", test.desc, test.want, got, err)
			}
		})
	}
}

func TestEncodeErrors(t *testing.T) {
	tests := []struct {
		desc    string
		in      *Params
		wantErr error
	}{
		{
			desc:    "country can't be mixed with sources param",
			in:      &Params{Country: "us", Sources: "the-times-of-india"},
			wantErr: ErrMixParams,
		},
		{
			desc:    "category can't be mixed with sources param",
			in:      &Params{Category: "technology", Sources: "the-times-of-india"},
			wantErr: ErrMixParams,
		},
		{
			desc:    "pageSize exceeded the maxPageSize",
			in:      &Params{PageSize: 500, Sources: "the-times-of-india"},
			wantErr: ErrInvalidPageSize,
		},
		{
			desc:    "no required parameter is supplied",
			in:      &Params{PageSize: 50, Page: 2},
			wantErr: ErrNoRequiredParams,
		},
		{
			desc:    "no parameter is supplied",
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
