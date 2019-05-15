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

// AESGCM mplements the Encrypt/Decrypt methods
// using AES-GCM: https://eprint.iacr.org/2015/102.pdf
type AESGCM struct {
	// either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	cipher cipher.Block
}

// validate interface conformity.
var _ encrypto.Cryptor = &AESGCM{}

// New makes a new AES/GCM Cryptor.
//
// The key argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func New(key string) (*AESGCM, error) {
	switch len(key) {
	default:
		return nil, errors.Errorf("AESGCM: key size invalid %d", len(key))
	case 16, 24, 32:
	}
	_key := []byte(key)
	_cipher, err := aes.NewCipher(_key)
	if err != nil {
		return nil, errors.Wrap(err, "AESGCM: unable to create a new cipher")
	}
	return &AESGCM{cipher: _cipher}, nil
}

// Encrypt ciphers the plainText using the provided 16, 24 or 32 bytes key
// with AES/GCM and returns a base64 encoded string.
func (ag *AESGCM) Encrypt(plainText []byte) (cypherText []byte, err error) {
	gcm, err := cipher.NewGCM(ag.cipher)
	if err != nil {
		return []byte{}, errors.Wrap(err, "unable to wrap cipher in GCM")
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, errors.Wrap(err, "unable to read random nonce")
	}
	b := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return b, nil
}

// Decrypt deciphers the provided base64 encoded and AES/GCM ciphered
// data returning the original plainText string.
func (ag *AESGCM) Decrypt(cipherText []byte) (plainText []byte, err error) {
	gcm, err := cipher.NewGCM(ag.cipher)
	if err != nil {
		return nil, errors.Wrap(err, "unable to wrap cipher in GCM")
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.Wrap(err, "unable to read random nonce")
	}
	nonce, cipherplainText := cipherText[:nonceSize], cipherText[nonceSize:]
	return gcm.Open(nil, nonce, cipherplainText, nil)
}
