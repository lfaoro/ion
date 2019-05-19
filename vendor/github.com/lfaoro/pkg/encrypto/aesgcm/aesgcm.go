// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package AESGCM implement AES encryption with GCM authentication according
// to the paper at ref: https://eprint.iacr.org/2015/102.pdf
package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"

	"github.com/lfaoro/pkg/encrypto"
)

// AESGCM implements the Encrypt/Decrypt methods
// using AES-GCM: https://eprint.iacr.org/2015/102.pdf
type AESGCM struct {
	block cipher.Block
}

// validate interface conformity.
var _ encrypto.Cryptor = &AESGCM{}

// New makes a new AES/GCM Cryptor. In order to select AES-256, a 32-byte key
// is enforced.
//
// ref: https://github.com/gtank/cryptopasta/blob/master/encrypt.go
func New(key *[32]byte) (*AESGCM, error) {
	if key == nil {
		return nil, errors.New("aesgcm: invalid key")
	}
	_block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, errors.Wrap(err, "aesgcm: unable to create a new cipher")
	}

	return &AESGCM{block: _block}, nil
}

// Encrypt encrypts data using 256-bit AES-GCM. This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (ag *AESGCM) Encrypt(plainText []byte) (cypherText []byte, err error) {
	gcm, err := cipher.NewGCM(ag.block)
	if err != nil {
		return nil, errors.Wrap(err, "unable to wrap cipher in GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read random nonce")
	}

	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM. This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (ag *AESGCM) Decrypt(cipherText []byte) (plainText []byte, err error) {
	gcm, err := cipher.NewGCM(ag.block)
	if err != nil {
		return nil, errors.Wrap(err, "unable to wrap cipher in GCM")
	}

	if len(cipherText) < gcm.NonceSize() {
		return nil, errors.Wrap(err, "unable to read random nonce")
	}

	return gcm.Open(nil,
		cipherText[:gcm.NonceSize()],
		cipherText[gcm.NonceSize():],
		nil,
	)
}
