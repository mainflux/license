// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"time"

	"github.com/mainflux/license"
	"github.com/mainflux/license/errors"
)

var (
	errEmptyServices = errors.New("the list of services must not be empty")
	errEmptyDeviceID = errors.New("device id must not be empty")
	errEmptyPayload  = errors.New("validation payload must not be empty")
)

type apiReq interface {
	validate() error
}

type licenseReq struct {
	token string
	id    string
}

func (req licenseReq) validate() error {
	if req.token == "" || req.id == "" {
		return license.ErrMalformedEntity
	}
	return nil
}

type createReq struct {
	token    string
	Duration time.Duration          `json:"duration,omitempty"`
	Services []string               `json:"services,omitempty"`
	DeviceID string                 `json:"device_id,omitempty"`
	Plan     map[string]interface{} `json:"plan,omitempty"`
}

func (req createReq) validate() error {
	if req.token == "" {
		return license.ErrUnauthorizedAccess
	}
	if req.Services == nil || len(req.Services) == 0 {
		return errors.Wrap(errEmptyServices, license.ErrMalformedEntity)
	}

	if req.DeviceID == "" {
		return errors.Wrap(errEmptyDeviceID, license.ErrMalformedEntity)
	}

	return nil
}

type updateReq struct {
	token    string
	id       string
	Services []string               `json:"services,omitempty"`
	Duration uint                   `json:"duration,omitempty"`
	Plan     map[string]interface{} `json:"plan,omitempty"`
}

func (req updateReq) validate() error {
	if req.token == "" {
		return license.ErrUnauthorizedAccess
	}
	if req.id == "" {
		return license.ErrNotFound
	}
	if req.Services == nil || len(req.Services) == 0 {
		return errors.Wrap(errEmptyServices, license.ErrMalformedEntity)
	}

	return nil
}

type validateReq struct {
	id      string
	service string
	Payload []byte `json:"payload,omitempty"`
}

func (req validateReq) validate() error {
	if req.id == "" {
		return license.ErrNotFound
	}
	if req.service == "" {
		return license.ErrMalformedEntity
	}
	if req.Payload == nil || len(req.Payload) == 0 {
		return errors.Wrap(errEmptyPayload, license.ErrMalformedEntity)
	}

	return nil
}
