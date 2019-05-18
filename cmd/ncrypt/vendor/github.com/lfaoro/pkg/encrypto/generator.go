// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encrypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Cryptographically secure pseudo-random number generator (CSPRNG)
// https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator
//
// BEST PRACTICE:
//
// Do not use package math/rand to generate keys, even throwaway ones.
// Unseeded, the generator is completely predictable. Seeded with
// time.Nanoseconds() there are just a few bits of entropy. Instead, use
// crypto/rand's Reader, and if you need text, print to hexadecimal or base64
//
// do NOT use rand.Seed(time.Now().UnixNano())

// RandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
//
// In the rare case you get an error there's something seriously
// wrong with your operating system.
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	return b
}

// RandomString returns a URL-safe, base64 encoded
// securely generated random string.
//
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandomString(s int) string {
	b := RandomBytes(s + 1)
	return base64.URLEncoding.EncodeToString(b)[:len(b)-1]
}

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}
