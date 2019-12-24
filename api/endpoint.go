// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/license"
)

func createEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		l := license.License{
			DeviceID:  req.DeviceID,
			Services:  req.Services,
			Plan:      req.Plan,
			CreatedAt: time.Now().UTC(),
		}
		l.ExpiresAt = l.CreatedAt.Add(req.Duration * time.Second)

		saved, err := svc.Create(ctx, req.token, l)
		if err != nil {
			return nil, err
		}

		res := licenseRes{
			ID:      saved,
			created: true,
		}

		return res, nil
	}
}

func viewEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(licenseReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		l, err := svc.Retrieve(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := vewRes{
			licenseRes: licenseRes{
				ID:        l.ID,
				created:   false,
				Issuer:    l.Issuer,
				DeviceID:  l.DeviceID,
				Active:    l.Active,
				CreatedAt: &l.CreatedAt,
				ExpiresAt: &l.ExpiresAt,
				UpdatedAt: &l.UpdatedAt,
				UpdatedBy: l.UpdatedBy,
				Services:  l.Services,
				Plan:      l.Plan,
			},
			Key: l.Key,
		}
		return res, nil
	}
}

func fetchEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(licenseReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		l, err := svc.Fetch(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := fetchRes{
			licenseRes: licenseRes{
				ID:        l.ID,
				created:   false,
				Issuer:    l.Issuer,
				DeviceID:  l.DeviceID,
				Active:    l.Active,
				CreatedAt: &l.CreatedAt,
				ExpiresAt: &l.ExpiresAt,
				UpdatedAt: &l.UpdatedAt,
				UpdatedBy: l.UpdatedBy,
				Services:  l.Services,
				Plan:      l.Plan,
			},
			Signature: l.Signature,
		}
		return res, nil
	}
}

func updateEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		l := license.License{
			ID:       req.id,
			Services: req.Services,
			Plan:     req.Plan,
		}

		err := svc.Update(ctx, req.token, l)
		if err != nil {
			return nil, err
		}

		return successRes{}, nil

	}
}

func removeEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(licenseReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		err := svc.Remove(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}

func validateEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		if err := svc.Validate(ctx, req.service, req.id, req.Payload); err != nil {
			return nil, err
		}

		return successRes{}, nil
	}
}

func activationEndpoint(svc license.Service, active bool) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(licenseReq)

		if err := req.validate(); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}

		err := svc.ChangeActive(ctx, req.token, req.id, active)
		if err != nil {
			return nil, err
		}

		return successRes{}, nil
	}
}
