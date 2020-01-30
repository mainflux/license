// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"

	"github.com/mainflux/license"
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

type successRes struct{}

func (res successRes) Code() int {
	return http.StatusOK
}

func (res successRes) Headers() map[string]string {
	return map[string]string{}
}

func (res successRes) Empty() bool {
	return true
}

type licenseRes struct {
	license.License
	created bool
}

func (res licenseRes) Code() int {
	if res.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (res licenseRes) Headers() map[string]string {
	ret := make(map[string]string)
	if res.created {
		ret["Location"] = fmt.Sprintf("/licenses/%s", res.ID)
	}

	return ret
}

func (res licenseRes) Empty() bool {
	return res.created
}

type errorRes struct {
	Err string `json:"error"`
}
