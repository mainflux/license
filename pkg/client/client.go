package client

import (
	"github.com/mainflux/license"
)

type client struct {
	url    string
	crypto license.Crypto
}

var _ license.Client = (*client)(nil)

// New returns new license client.
func New(url string, crypto license.Crypto) license.Client {
	return client{
		url:    url,
		crypto: crypto,
	}
}

func (c client) Validate(svcName string) error {
	return license.ErrNotFound
}
