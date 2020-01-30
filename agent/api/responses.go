// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import "net/http"

type errorRes struct {
	Err string `json:"error"`
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
