// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"context"
	"time"

	"github.com/mainflux/license/errors"
)

// License represents single license object.
type License struct {
	ID        string                 `json:"id"`
	Key       string                 `json:"key"`
	Issuer    string                 `json:"issuer"`
	DeviceID  string                 `json:"device_id"`
	Active    bool                   `json:"active"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	UpdatedBy string                 `json:"updated_by"`
	UpdatedAt time.Time              `json:"updated_at"`
	Services  []string               `json:"services"`
	Plan      map[string]interface{} `json:"plan"`
	Signature []byte                 `json:"signature"`
}

// Validate validates the license.
func (l License) Validate() error {
	now := time.Now().UTC()
	if l.ExpiresAt.UTC().Before(now) {
		return ErrExpired
	}
	if l.CreatedAt.UTC().After(now) {
		return errors.Wrap(ErrMalformedEntity, ErrIssuedAt)
	}
	if !l.Active {
		return ErrLicenseValidation
	}
	return nil
}

// Repository specifies a License persistence API.
type Repository interface {
	// Save stores a License.
	Save(ctx context.Context, l License) (string, error)

	// Retrieve the License by given ID that belongs to the given owner.
	Retrieve(ctx context.Context, issuer, id string) (License, error)

	// RetrieveByID retrives the license by device ID.
	RetrieveByDeviceID(ctx context.Context, deviceID string) (License, error)

	// Update an existing License.
	Update(ctx context.Context, l License) error

	// Remove a License with the given ID that belongs to the given owner.
	Remove(ctx context.Context, issuer, id string) error

	// ChangeActive a License with the given ID
	// that belongs to the given issuer.
	ChangeActive(ctx context.Context, token, id string, active bool) error
}
