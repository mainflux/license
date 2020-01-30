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
		req := request.(validationReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		ret, err := agent.Validate(req.svcID, req.client)
		if err != nil {
			return nil, err
		}

		return ret, nil
	}
}
