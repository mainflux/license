// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

// Handler handles validation response.
type Handler func(error)

// Client represents licensing service client API specification.
type Client interface {
	// Validate validates service against provided license.
	Validate(svcName string)
}
