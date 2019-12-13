// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"context"
	errs "errors"
	"time"

	"github.com/mainflux/license/errors"
	"github.com/mainflux/mainflux"
)

var (
	ErrConflict           = errors.New("entity already exists")
	ErrNotFound           = errors.New("entity does not exist")
	ErrMalformedEntity    = errors.New("malformed entity data")
	ErrUnauthorizedAccess = errors.New("unauthorized access")

	errIssuedAt = errs.New("invalid issue data")
)

// Service represents licensing service API specification.
type Service interface {
	// Create adds License that belongs to the
	// user identified by the provided token.
	Create(ctx context.Context, token string, l License) (string, error)

	// Retrieve retrieves the License by given ID that belongs to
	//  the user identified by the provided token.
	Retrieve(ctx context.Context, token, id string) (License, error)

	// Update updates an existing License that's issued
	// by the given issuer.
	Update(ctx context.Context, token string, l License) error

	// Remove removes a License with the given ID
	// that belongs to the given issuer.
	Remove(ctx context.Context, token, id string) error

	// ChangeActive a License with the given ID
	// that belongs to the given issuer.
	ChangeActive(ctx context.Context, token, id string, active bool) error
}

type licenseService struct {
	repo Repository
	idp  IdentityProvider
	auth mainflux.UsersServiceClient
}

// New returns new instance of License Service.
func New(repo Repository, idp IdentityProvider, auth mainflux.UsersServiceClient) Service {
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
	return svc.repo.Retrieve(ctx, issuer.GetValue(), id)
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
