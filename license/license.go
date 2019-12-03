// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import "time"

import "context"

// License represents single license object.
type License struct {
	ID       string
	Owner    string
	Active   bool
	Created  time.Time
	Duration *uint
	Expires  *time.Time
	Metadata map[string]interface{}
	Plan     map[string]interface{}
}

// Repository specifies a License persistence API.
type Repository interface {
	// Save stores a License.
	Save(context.Context, License) (string, error)

	// Retrieve the License by given ID that belongs to the given owner.
	Retrieve(context.Context, string, string) (License, error)

	// Update an existing License.
	Update(context.Context, License) error

	// Remove a License with the given ID that belongs to the given owner.
	Remove(context.Context, string, string) error
}
