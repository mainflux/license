package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/license"
)

type client struct {
	url     string
	crypto  license.Crypto
	handler license.Handler
}

var _ license.Client = (*client)(nil)

// New returns new license client.
func New(url string, crypto license.Crypto, handler license.Handler) license.Client {
	return client{
		url:     url,
		crypto:  crypto,
		handler: handler,
	}
}

func (c client) Validate(svcName string) {
	var err error
	defer func() {
		c.handler(err)
	}()

	url := fmt.Sprintf("%s/%s", c.url, svcName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	key := "0ecd741f-7aa1-494a-aa94-4b544ef41e4a"

	k, err := c.crypto.Encrypt([]byte(key))
	if err != nil {
		return
	}

	key = hex.EncodeToString(k)
	req.Header.Set("Authorization", key)
	res, err := http.DefaultClient.Do(req)

	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return
	}
	if res.StatusCode == http.StatusOK {
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = errors.New(string(data))
}
