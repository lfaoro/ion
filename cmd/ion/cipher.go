// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

func encryptFile(filePath string, data []byte, engine *aesgcm.AESGCM) error {
	cipherText, err := engine.Encrypt(data)
	if err != nil {
		return err
	}

	cipherText = addHeader(cipherText)

	err = ioutil.WriteFile(filePath, cipherText, 0600)
	if err != nil {
		return err
	}

	return nil
}

func decryptFile(filePath string, data []byte, engine *aesgcm.AESGCM) error {
	cipherText, err := removeHeader(data)
	if err != nil {
		return err
	}

	plainText, err := engine.Decrypt(cipherText)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, plainText, 0600)
	if err != nil {
		return err
	}

	return nil
}

func removeHeader(data []byte) ([]byte, error) {
	if !isEncrypted(data) {
		return []byte{}, errors.New("this file is not encrypted")
	}
	i := bytes.IndexByte(data, byte('\n'))
	if i == -1 {
		return []byte{}, errors.New("invalid ncrypt file")
	}
	return data[i+1:], nil
}

func addHeader(data []byte) []byte {
	header := getHeader()
	return append(header, data...)
}

func isEncrypted(data []byte) bool {
	return bytes.Contains(data, getHeader())
}
