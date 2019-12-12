// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
	"time"

	"github.com/mainflux/mainflux"
)

var (
	_ mainflux.Response = (*removeRes)(nil)
	_ mainflux.Response = (*licenseRes)(nil)
)

type removeRes struct{}

func (res removeRes) Code() int {
	return http.StatusNoContent
}

func (res removeRes) Headers() map[string]string {
	return map[string]string{}
}

func (res removeRes) Empty() bool {
	return true
}

type licenseRes struct {
	created   bool
	ID        string                 `json:"id,omitempty"`
	Issuer    string                 `json:"issuer,omitempty"`
	DeviceID  string                 `json:"device_id,omitempty"`
	Active    bool                   `json:"active,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitempty"`
	ExpiresAt time.Time              `json:"expires_at,omitempty"`
	UpdatedBy string                 `json:"updated_by,omitempty"`
	UpdatedAt time.Time              `json:"updated_at,omitempty"`
	Services  []string               `json:"services,omitempty"`
	Plan      map[string]interface{} `json:"plan,omitempty"`
}

func (res licenseRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res licenseRes) Headers() map[string]string {
	return map[string]string{}
}

func (res licenseRes) Empty() bool {
	return false
}

type errorRes struct {
	Err string `json:"error"`
}
