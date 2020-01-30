// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/mainflux/license"
	"github.com/mainflux/license/errors"
)

var errEmptyPayload = errors.New("validation payload must not be empty")

type validationReq struct {
	svcID  string
	client string
}

func (req validationReq) validate() error {
	if req.svcID == "" || req.client == "" {
		return license.ErrNotFound
	}

	return nil
}
