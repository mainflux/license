// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/mainflux/license"
	"time"
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
	Plan     map[string]interface{} `json:"plan,omitempty"`
}

func (req createReq) validate() error {
	if req.token == "" {
		return license.ErrUnauthorizedAccess
	}
	if req.Services == nil || len(req.Services) == 0 {
		return license.ErrMalformedEntity
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
		return license.ErrUnauthorizedAccess
	}

	return nil
}
