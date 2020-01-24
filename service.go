// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"context"
	errs "errors"

	"github.com/mainflux/license/errors"
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

	// ErrIssuedAt represents invalid issue date.
	ErrIssuedAt = errs.New("invalid issue date")
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
