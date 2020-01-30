// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/license"
)

type validator struct {
	url     string
	crypto  license.Crypto
	handler license.Handler
}

var _ license.Validator = (*validator)(nil)

// New returns new license validator.
func New(url string, crypto license.Crypto, handler license.Handler) license.Validator {
	return validator{
		url:     url,
		crypto:  crypto,
		handler: handler,
	}
}

func (v validator) Validate(svcName, client string) (err error) {
	defer func() {
		v.handler(err)
	}()

	url := fmt.Sprintf("%s/%s", v.url, svcName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", client)
	res, err := http.DefaultClient.Do(req)

	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	dec, err := v.crypto.Decrypt(data)
	if err != nil {
		return err
	}
	var r validateResponse
	if err := json.Unmarshal(dec, &r); err != nil {
		return err
	}
	if r.Status == http.StatusOK {
		return nil
	}

	return errors.New(r.Message)
}

type validateResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
