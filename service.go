// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	errs "errors"
	"time"

	"github.com/mainflux/license/errors"
	"github.com/mainflux/mainflux"
)

var (
	// ErrConflict represents unique identifier violation.
	ErrConflict = errors.New("entity already exists")

	// ErrNotFound represents non-existing entity request.
	ErrNotFound = errors.New("entity does not exist")

	// ErrMalformedEntity represents malformed entity specification.
	ErrMalformedEntity = errors.New("malformed entity data")

	// ErrUnauthorizedAccess represents missing or invalid credentials.
	ErrUnauthorizedAccess = errors.New("unauthorized access")

	// ErrExpired represents expired license error.
	ErrExpired = errs.New("the license is expired")

	errIssuedAt = errs.New("invalid issue date")
)

// Service represents licensing service API specification.
type Service interface {
	// Create adds License that belongs to the
	// user identified by the provided token.
	Create(ctx context.Context, token string, l License) (string, error)

	// Retrieve retrieves the License by given ID that belongs to
	//  the user identified by the provided token.
	Retrieve(ctx context.Context, token, id string) (License, error)

	// Fetch retrieves License using license ID and Key.
	Fetch(ctx context.Context, key, id string) (License, error)

	// Update updates an existing License that's issued
	// by the given issuer.
	Update(ctx context.Context, token string, l License) error

	// Remove removes a License with the given ID
	// that belongs to the given issuer.
	Remove(ctx context.Context, token, id string) error

	// ChangeActive a License with the given ID
	// that belongs to the given issuer.
	ChangeActive(ctx context.Context, token, id string, active bool) error

	// Validate checks if the license is valid for the given service name.
	Validate(ctx context.Context, svcName, id string, payload []byte) error
}

type licenseService struct {
	repo Repository
	idp  IdentityProvider
	auth mainflux.AuthNServiceClient
}

// New returns new instance of License Service.
func New(repo Repository, idp IdentityProvider, auth mainflux.AuthNServiceClient) Service {
	return licenseService{
		repo: repo,
		idp:  idp,
		auth: auth,
	}
}

func (svc licenseService) Create(ctx context.Context, token string, l License) (string, error) {
	if l.CreatedAt.IsZero() {
		return "", errors.Wrap(ErrMalformedEntity, errIssuedAt)
	}
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(ErrUnauthorizedAccess, err)
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

func (svc licenseService) Retrieve(ctx context.Context, token, id string) (License, error) {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return License{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}

	l, err := svc.repo.Retrieve(ctx, issuer.GetValue(), id)
	if err != nil {
		return License{}, err
	}

	return l, nil
}

func (svc licenseService) Fetch(ctx context.Context, key, id string) (License, error) {
	l, err := svc.repo.RetrieveByID(ctx, id)
	if err != nil {
		return License{}, err
	}
	if l.Key != key {
		return License{}, ErrUnauthorizedAccess
	}
	if err := l.Validate(); err != nil {
		return License{}, err
	}
	bytes, err := json.Marshal(l)
	if err != nil {
		return License{}, errors.Wrap(ErrMalformedEntity, err)
	}
	h := hmac.New(sha256.New, []byte(l.Key))
	if _, err := h.Write(bytes); err != nil {
		return License{}, ErrMalformedEntity
	}
	l.Signature = h.Sum(nil)

	return l, nil
}

func (svc licenseService) Update(ctx context.Context, token string, l License) error {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
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
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.repo.Remove(ctx, issuer.GetValue(), id)
}

func (svc licenseService) ChangeActive(ctx context.Context, token, id string, active bool) error {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
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
		return ErrMalformedEntity
	}
	if !hmac.Equal(payload, h.Sum(nil)) {
		return ErrMalformedEntity
	}

	for _, s := range l.Services {
		if s == name {
			return nil
		}
	}

	return ErrNotFound
}
