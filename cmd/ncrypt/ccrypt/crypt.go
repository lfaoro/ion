// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crypt

import (
	"bytes"

	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

var header = []byte("## ncrypted with love")

func Encrypt(plaintext []byte, key *[32]byte) ([]byte, error) {
	engine, err := aesgcm.New(key)
	ciphertext, err := engine.Encrypt(plaintext)
	ciphertext = addHeader(ciphertext)
	return ciphertext, err
}

func addHeader(data []byte) []byte {
	return append(header, data...)
}

func Decrypt(ciphertext []byte, key *[32]byte) ([]byte, error) {
	engine, err := aesgcm.New(key)
	plaintext, err := engine.Decrypt(ciphertext)
	plaintext = removeHeader(plaintext)
	return plaintext, err
}

func removeHeader(data []byte) []byte {
	if !bytes.Contains(data, header) {
		panic("this file is not encrypted")
	}
	i := bytes.IndexByte(data, byte('\n'))
	if i == -1 {
		panic("invalid ncrypt file")
	}
	return data[i+1:]
}
