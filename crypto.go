// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package license

// Crypto provides an interface for reading and storing of the License.
type Crypto interface {
	// Encrypt encrypts license before storing in file.
	Encrypt([]byte) ([]byte, error)

	// Descrypt decrypts license retrieved from file or cloud.
	Decrypt([]byte) ([]byte, error)
}
