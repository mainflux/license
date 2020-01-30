// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

// Handler handles validation result.
type Handler func(error)

// Validator represents licensing service validator specification.
type Validator interface {
	// Validate validates service against provided license.
	Validate(svcName, client string) error
}
