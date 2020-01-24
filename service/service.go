// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/mainflux/license"
	"github.com/mainflux/license/errors"
	"github.com/mainflux/mainflux"
)

type licenseService struct {
	repo license.Repository
	idp  license.IdentityProvider
	auth mainflux.AuthNServiceClient
}

// New returns new instance of License Service.
func New(repo license.Repository, idp license.IdentityProvider, auth mainflux.AuthNServiceClient) license.Service {
	return licenseService{
		repo: repo,
		idp:  idp,
		auth: auth,
	}
}

func (svc licenseService) Create(ctx context.Context, token string, l license.License) (string, error) {
	if l.CreatedAt.IsZero() {
		return "", errors.Wrap(license.ErrMalformedEntity, license.ErrIssuedAt)
	}
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(license.ErrUnauthorizedAccess, err)
	}
	l.ID, err = svc.idp.ID()
	if err != nil {
		return "", err
	}
	l.Key, err = svc.idp.ID()
	if err != nil {
		return "", err
	}

	l.Issuer = issuer.GetValue()
	l.UpdatedAt = l.CreatedAt
	l.UpdatedBy = l.Issuer
	return svc.repo.Save(ctx, l)
}

func (svc licenseService) Retrieve(ctx context.Context, token, id string) (license.License, error) {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return license.License{}, errors.Wrap(license.ErrUnauthorizedAccess, err)
	}

	l, err := svc.repo.Retrieve(ctx, issuer.GetValue(), id)
	if err != nil {
		return license.License{}, err
	}

	return l, nil
}

func (svc licenseService) Fetch(ctx context.Context, key, id string) (license.License, error) {
	l, err := svc.repo.RetrieveByID(ctx, id)
	if err != nil {
		return license.License{}, err
	}
	if l.Key != key {
		return license.License{}, license.ErrUnauthorizedAccess
	}
	if err := l.Validate(); err != nil {
		return license.License{}, err
	}
	bytes, err := json.Marshal(l)
	if err != nil {
		return license.License{}, errors.Wrap(license.ErrMalformedEntity, err)
	}
	h := hmac.New(sha256.New, []byte(l.Key))
	if _, err := h.Write(bytes); err != nil {
		return license.License{}, license.ErrMalformedEntity
	}
	l.Signature = h.Sum(nil)

	return l, nil
}

func (svc licenseService) Update(ctx context.Context, token string, l license.License) error {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(license.ErrUnauthorizedAccess, err)
	}
	iss := issuer.GetValue()
	l.Issuer = iss
	l.UpdatedBy = iss
	l.UpdatedAt = time.Now().UTC()
	return svc.repo.Update(ctx, l)
}

func (svc licenseService) Remove(ctx context.Context, token, id string) error {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(license.ErrUnauthorizedAccess, err)
	}
	return svc.repo.Remove(ctx, issuer.GetValue(), id)
}

func (svc licenseService) ChangeActive(ctx context.Context, token, id string, active bool) error {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(license.ErrUnauthorizedAccess, err)
	}

	return svc.repo.ChangeActive(ctx, issuer.GetValue(), id, active)
}

func (svc licenseService) Validate(ctx context.Context, name, id string, payload []byte) error {
	l, err := svc.repo.RetrieveByID(ctx, id)
	if err != nil {
		return err
	}
	if err := l.Validate(); err != nil {
		return err
	}

	h := hmac.New(sha256.New, []byte(l.Key))
	if _, err := h.Write([]byte(l.DeviceID)); err != nil {
		return license.ErrMalformedEntity
	}
	if !hmac.Equal(payload, h.Sum(nil)) {
		return license.ErrMalformedEntity
	}

	for _, s := range l.Services {
		if s == name {
			return nil
		}
	}

	return license.ErrNotFound
}
