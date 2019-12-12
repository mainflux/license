// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"context"
	errs "errors"

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
	// CreateLicense adds License that belongs to the
	// user identified by the provided key.
	CreateLicense(ctx context.Context, token string, l License) (string, error)

	// RetrieveLicense retrieves the License by given ID that belongs to
	//  the user identified by the provided token.
	RetrieveLicense(ctx context.Context, token, id string) (License, error)

	// UpdateLicense updates an existing License.
	UpdateLicense(ctx context.Context, l License) error

	// RemoveLicense removes a License with the given ID
	// that belongs to the given owner.
	RemoveLicense(ctx context.Context, token, id string) error
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

func (svc licenseService) CreateLicense(ctx context.Context, token string, l License) (string, error) {
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

	l.Issuer = issuer.GetValue()
	l.UpdatedAt = l.CreatedAt
	l.UpdatedBy = l.Issuer
	return svc.repo.Save(ctx, l)
}

func (svc licenseService) RetrieveLicense(ctx context.Context, token, id string) (License, error) {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return License{}, errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.repo.Retrieve(ctx, issuer.GetValue(), id)
}

func (svc licenseService) UpdateLicense(ctx context.Context, l License) error {
	return svc.repo.Update(ctx, l)
}

func (svc licenseService) RemoveLicense(ctx context.Context, token, id string) error {
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return errors.Wrap(ErrUnauthorizedAccess, err)
	}
	return svc.repo.Remove(ctx, issuer.GetValue(), id)
}
