// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

// Handler handles validation response.
type Handler func(error)

// Validator represents licensing service validator specification.
type Validator interface {
	// Validate validates service against provided license.
	Validate(svcName string) error
}
