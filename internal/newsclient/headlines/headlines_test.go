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
		desc string
		in   *Params
	}{
		{
			desc: "country can't be mixed with sources param",
			in:   &Params{Country: "us", Sources: "the-times-of-india"},
		},
		{
			desc: "category can't be mixed with sources param",
			in:   &Params{Category: "technology", Sources: "the-times-of-india"}},
		{
			desc: "pageSize exceeded the maxPageSize",
			in:   &Params{PageSize: 500, Category: "technology", Sources: "the-times-of-india"},
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
