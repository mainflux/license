package client

import "github.com/mainflux/license"

type client struct {
	url string
}

var _ license.Client = (*client)(nil)

// New returns new license client.
func New(url string) license.Client {
	return client{
		url: url,
	}
}
func (c client) Fetch(id, key string) (license.License, error) {
	return license.License{}, nil
}

func (c client) Validate(l license.License, svcName string) error {
	if err := l.Validate(); err != nil {
		return err
	}
	for _, svc := range l.Services {
		if svc == svcName {
			return nil
		}
	}
	return license.ErrNotFound
}
