// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/license"
)

type validator struct {
	url     string
	key     string
	crypto  license.Crypto
	handler license.Handler
}

var _ license.Validator = (*validator)(nil)

// New returns new license validator.
func New(url, key string, crypto license.Crypto, handler license.Handler) license.Validator {
	return validator{
		url:     url,
		key:     key,
		crypto:  crypto,
		handler: handler,
	}
}

func (v validator) Validate(svcName string) (err error) {
	defer func() {
		v.handler(err)
	}()

	url := fmt.Sprintf("%s/%s", v.url, svcName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	k, err := v.crypto.Encrypt([]byte(v.key))
	if err != nil {
		return err
	}

	key := hex.EncodeToString(k)
	req.Header.Set("Authorization", key)
	res, err := http.DefaultClient.Do(req)

	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return errors.New(string(data))
}
