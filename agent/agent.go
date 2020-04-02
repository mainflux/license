// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mainflux/license"
	"github.com/mainflux/license/errors"
)

var (
	errServiceNotAllowed = errors.New("service not allowed")
	errLicenseNotLoaded  = errors.New("license not loaded")
)
var _ license.Agent = (*agent)(nil)

type action uint

type command struct {
	action action
	param  string
}

const (
	read action = iota
	write
	validate
)

type agent struct {
	svcURL    string
	location  string
	id        string
	key       string
	license   *license.License
	crypto    license.Crypto
	validator license.Validator
	commands  chan command
	errs      chan error
}

// New returns new License agent.
func New(svcURL, location string, id, key string, crypto license.Crypto, validator license.Validator) license.Agent {
	return &agent{
		svcURL:    svcURL,
		location:  location,
		id:        id,
		key:       key,
		crypto:    crypto,
		validator: validator,
		commands:  make(chan command),
		errs:      make(chan error),
	}
}

func (a *agent) Do() {
	for {
		cmd := <-a.commands
		var err error
		switch cmd.action {
		case read:
			var l license.License
			l, err = a.load()
			if err == nil {
				a.license = &l
				err = a.save()
			}
		case validate:
			err = a.validate(cmd.param)
		case write:
			err = a.save()
		}
		a.errs <- err
	}
}

func (a *agent) Load() error {
	a.commands <- command{action: read}
	return <-a.errs
}

func (a *agent) Save() error {
	a.commands <- command{action: write}
	return <-a.errs
}

func (a *agent) Validate(r []byte) ([]byte, error) {
	val, err := a.req(r)
	if err != nil {
		return nil, err
	}
	cmd := command{
		action: validate,
		param:  val.SvcID,
	}

	// Validate service against license.
	a.commands <- cmd
	ret := validateResponse{Status: http.StatusForbidden}
	err = <-a.errs
	if err != nil {
		ret.Message = err.Error()
	}

	// Optional custom validation.
	if a.validator != nil && err == nil {
		if err := a.validator.Validate(val.SvcID, val.Client); err != nil {
			ret = validateResponse{
				Status:  http.StatusForbidden,
				Message: err.Error(),
			}
		}
	}

	resp, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}

	return a.crypto.Encrypt(resp)
}

// Unlike their exported counterparts, methods load, save, and validate are not thread-safe.
// These methods are meant to be executed as command handlers in Do method.
func (a *agent) load() (license.License, error) {
	data, err := ioutil.ReadFile(a.location)
	switch {
	case err == nil:
		break
	case os.IsNotExist(err):
		data, err = a.fetch()
		if err != nil {
			return license.License{}, err
		}
	default:
		return license.License{}, err
	}
	data, err = a.crypto.Decrypt(data)
	if err != nil {
		return license.License{}, err
	}
	l := license.License{}
	err = json.Unmarshal(data, &l)
	return l, err
}

func (a *agent) save() error {
	if a.license == nil {
		return errLicenseNotLoaded
	}
	data, err := json.Marshal(*a.license)
	if err != nil {
		return err
	}
	data, err = a.crypto.Encrypt(data)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(a.location, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (a *agent) validate(svcName string) error {
	if a.license == nil {
		return errors.Wrap(license.ErrLicenseValidation, errLicenseNotLoaded)
	}
	if err := a.license.Validate(); err != nil {
		return err
	}
	for _, svc := range a.license.Services {
		if svcName == svc {
			return nil
		}
	}
	return errors.Wrap(license.ErrLicenseValidation, errServiceNotAllowed)
}

func (a *agent) fetch() ([]byte, error) {
	id, err := a.crypto.Encrypt([]byte(a.id))
	if err != nil {
		return nil, err
	}
	q := hex.EncodeToString(id)
	req, err := http.NewRequest(http.MethodGet, a.svcURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", q)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

type validateResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
