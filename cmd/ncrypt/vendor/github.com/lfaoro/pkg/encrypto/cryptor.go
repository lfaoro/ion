// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encrypto

// Cryptor defines the methods for encryption and decryption of data.
type Cryptor interface {
	Encrypt(plainText []byte) (cipherText []byte, err error)
	Decrypt(cipherText []byte) (plainText []byte, err error)
}
