// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/license/license"
)

func createLicenseEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createLicenseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		l := license.License{
			Owner:    req.token,
			Metadata: req.Metadata,
			Plan:     req.Plan,
		}
		saved, err := svc.CreateLicense(ctx, l)
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
			ID:      l.ID,
			created: false,
		}

		return res, nil
	}
}

func updateLicenseEndpoint(svc license.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createLicenseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		l := license.License{
			Owner:    req.token,
			Metadata: req.Metadata,
			Plan:     req.Plan,
		}

		err := svc.UpdateLicense(ctx, l)
		if err != nil {
			return nil, err
		}

		res := licenseRes{
			ID:      l.ID,
			created: true,
		}

		return res, nil
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
