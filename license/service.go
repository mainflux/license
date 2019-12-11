// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"context"
	"errors"

	"github.com/mainflux/mainflux"
)

var (
	ErrConflict           = errors.New("entity already exists")
	ErrNotFound           = errors.New("entity does not exist")
	ErrMalformedEntity    = errors.New("malformed entity data")
	ErrUnauthorizedAccess = errors.New("unauthorized access")
)

// Service represents licensing service API specification.
type Service interface {
	// CreateLicense adds License that belongs to the
	// user identified by the provided key.
	CreateLicense(context.Context, string, License) (string, error)

	// RetrieveLicense retrieves the License by given ID that belongs to
	//  the user identified by the provided key.
	RetrieveLicense(context.Context, string, string) (License, error)

	// UpdateLicense updates an existing License.
	UpdateLicense(context.Context, License) error

	// RemoveLicense removes a License with the given ID
	// that belongs to the given owner.
	RemoveLicense(context.Context, string, string) error
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
	issuer, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", err
	}

	l.ID, err = svc.idp.ID()
	l.Issuer = issuer.GetValue()
	if err != nil {
		return "", err
	}
	return svc.repo.Save(ctx, l)
}

func (svc licenseService) RetrieveLicense(ctx context.Context, id, owner string) (License, error) {
	return svc.repo.Retrieve(ctx, id, owner)
}

func (svc licenseService) UpdateLicense(ctx context.Context, l License) error {
	return svc.repo.Update(ctx, l)
}

func (svc licenseService) RemoveLicense(ctx context.Context, id, owner string) error {
	return svc.repo.Remove(ctx, id, owner)
}
