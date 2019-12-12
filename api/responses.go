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
	ID       string                 `json:"id"`
	Created  time.Time              `json:"created,omitempty"`
	Expires  time.Time              `json:"expires,omitempty"`
	Duration uint                   `json:"duration,omitempty"`
	Active   bool                   `json:"active,omitempty"`
	Plan     map[string]interface{} `json:"plan,omitempty"`
	created  bool
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
