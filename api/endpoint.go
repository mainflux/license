// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/license"
)

func createLicenseEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createLicenseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		l := license.License{
			Services:  req.Services,
			Plan:      req.Plan,
			CreatedAt: time.Now().UTC(),
		}
		l.ExpiresAt = l.CreatedAt.Add(req.Duration * time.Second)

		saved, err := svc.CreateLicense(ctx, req.token, l)
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

func viewLicenseEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(licenseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		l, err := svc.RetrieveLicense(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		res := licenseRes{
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
		}
		return res, nil
	}
}

func updateLicenseEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateLicenseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		l := license.License{
			ID:       req.id,
			Services: req.Services,
			Plan:     req.Plan,
		}

		err := svc.UpdateLicense(ctx, req.token, l)
		if err != nil {
			return nil, err
		}

		return updateRes{}, nil

	}
}

func removeLicenseEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(licenseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.RemoveLicense(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}

		return removeRes{}, nil
	}
}
