// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

// Client represents licensing service client API specification.
type Client interface {
	// Validate validates service against provided license.
	Validate(svcName string) error
}
