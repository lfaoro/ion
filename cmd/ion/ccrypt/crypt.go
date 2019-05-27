// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccrypt

import (
	"bytes"
	"fmt"

	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

var header = []byte("## ncrypted with love\n")

func Encrypt(plaintext []byte, key *[32]byte) ([]byte, error) {
	engine, err := aesgcm.New(key)
	if err != nil {
		return nil, err
	}

	ciphertext, err := engine.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	ciphertext = addHeader(ciphertext)
	fmt.Println(string(ciphertext[:20]))

	return ciphertext, err
}

func Decrypt(ciphertext []byte, key *[32]byte) ([]byte, error) {
	fmt.Println("key", string(key[:]))
	engine, err := aesgcm.New(key)
	if err != nil {
		return nil, err
	}
	ciphertext = removeHeader(ciphertext)
	plaintext, err := engine.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

func addHeader(data []byte) []byte {
	return append(header, data...)
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
