// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

import "github.com/mainflux/license/errors"

// ErrLicenseValidation wraps an error in case of unsuccessfull validation.
var ErrLicenseValidation = errors.New("license validation failed")

// Agent represents licensing agent.
// Licensing Agent is a service that handles License locally.
type Agent interface {
	// Validate validates service.
	Validate([]byte) ([]byte, error)

	// Load reads License from the location.
	Load() error

	// Save saves License to file.
	Save() error

	// Do runs the Agent.
	Do()
}
