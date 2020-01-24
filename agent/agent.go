// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package agent

import (
	"encoding/json"
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
	params []string
}

const (
	read action = iota
	write
	validate
)

type agent struct {
	svcURL   string
	location string
	commands chan command
	errs     chan error
	license  *license.License
}

// New returns new License agent.
func New(svcURL, location string) license.Agent {
	return &agent{
		svcURL:   svcURL,
		location: location,
		commands: make(chan command),
		errs:     make(chan error),
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
			}
		case validate:
			err = a.validate(cmd.params)
		case write:
			err = a.save()
		}
		a.errs <- err
	}
}

func (a *agent) Load() error {
	return a.command(read)
}

func (a *agent) Save() error {
	return a.command(write)
}

func (a *agent) Validate(services []string) error {
	return a.command(validate, services...)
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
	data, err = Dec(data)
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
	data, err = Enc(data)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(a.location, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (a *agent) validate(params []string) error {
	if a.license == nil {
		return errors.Wrap(license.ErrLicenseValidation, errLicenseNotLoaded)
	}
	if err := a.license.Validate(); err != nil {
		return err
	}
	valid := true
	for _, p := range params {
		if !exists(p, a.license.Services) {
			valid = false
			break
		}
	}
	if !valid {
		return errors.Wrap(license.ErrLicenseValidation, errServiceNotAllowed)
	}
	return nil
}

func (a *agent) command(act action, params ...string) error {
	cmd := command{
		action: act,
		params: params,
	}
	a.commands <- cmd
	return <-a.errs
}

func (a *agent) fetch() ([]byte, error) {
	res, err := http.DefaultClient.Get(a.svcURL)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	var data []byte
	if _, err := res.Body.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}

func exists(p string, services []string) bool {
	for _, s := range services {
		if string(p) == s {
			return true
		}
	}
	return false
}
