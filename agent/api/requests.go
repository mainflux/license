// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/mainflux/license"
	"github.com/mainflux/license/errors"
)

var errEmptyPayload = errors.New("validation payload must not be empty")

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
