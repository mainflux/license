// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import "time"

import "context"

// License represents single license object.
type License struct {
	ID        string
	Issuer    string
	DeviceID  string
	Active    bool
	CreatedAt time.Time
	ExpiresAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
	Services  []string
	Plan      map[string]interface{}
}

// Repository specifies a License persistence API.
type Repository interface {
	// Save stores a License.
	Save(ctx context.Context, l License) (string, error)

	// Retrieve the License by given ID that belongs to the given owner.
	Retrieve(ctx context.Context, issuer, id string) (License, error)

	// Update an existing License.
	Update(ctx context.Context, l License) error

	// Remove a License with the given ID that belongs to the given owner.
	Remove(ctx context.Context, issuer, id string) error

	// ChangeActive a License with the given ID
	// that belongs to the given issuer.
	ChangeActive(ctx context.Context, token, id string, active bool) error
}
