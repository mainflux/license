// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import "github.com/mainflux/license/license"

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

type createLicenseReq struct {
	token    string
	Duration uint                   `json:"duration,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Plan     map[string]interface{} `json:"plan,omitempty"`
}

func (req createLicenseReq) validate() error {
	if req.token == "" || req.Plan == nil || len(req.Plan) == 0 {
		return license.ErrUnauthorizedAccess
	}

	return nil
}

type updateLicenseReq struct {
	token    string
	id       string
	Duration uint                   `json:"duration,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Plan     map[string]interface{} `json:"plan,omitempty"`
}

func (req updateLicenseReq) validate() error {
	if req.token == "" {
		return license.ErrUnauthorizedAccess
	}
	if req.id == "" {
		return license.ErrNotFound
	}

	return nil
}
