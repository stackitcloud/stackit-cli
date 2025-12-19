package testutils

import "github.com/google/go-cmp/cmp"

type Option struct {
	cmpOptions []cmp.Option
}

type TestingOption func(options *Option) error

func WithCmpOptions(cmpOptions ...cmp.Option) TestingOption {
	return func(options *Option) error {
		options.cmpOptions = append(options.cmpOptions, cmpOptions...)
		return nil
	}
}
