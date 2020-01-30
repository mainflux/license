// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/license"
)

func validateEndpoint(agent license.Agent) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		if err := agent.Validate(req.service); err != nil {
			return nil, err
		}

		return successRes{}, nil
	}
}
