// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

// Client represents licensing service client API specification.
type Client interface {
	// Fetch fetches the License using license id and license key for authentication.
	Fetch(id, key string) (License, error)

	// Validate validates service against provided license.
	Validate(l License, svcName string) error
}
