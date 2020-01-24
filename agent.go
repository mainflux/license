// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package license

import "github.com/mainflux/license/errors"

// ErrLicenseValidation wraps an error in case of unsuccessfull validation.
var ErrLicenseValidation = errors.New("license validation failed")

// Agent represents licensing agent.
// Licensing Agent is a service that handles License locally.
type Agent interface {
	// Load reads License from the location.
	Load() error

	// Save saves License to file.
	Save() error

	// Validate validates service against provided license.
	Validate(svcNames []string) error

	// Do runs the Agent.
	Do()
}
