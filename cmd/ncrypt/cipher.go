// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

func cryptoCmd(engine *aesgcm.AESGCM, filePath string, data []byte) error {
	fileName := filepath.Base(filePath)
	if isEncrypted(data) {
		err := decryptFile(filePath, data, engine)
		if err != nil {
			return err
		}

		fmt.Printf("🔓 Decrypted %s\n", fileName)
		return nil
	}

	err := encryptFile(filePath, data, engine)
	if err != nil {
		return err
	}

	fmt.Printf("🔒 Encrypted %s\n", fileName)
	return nil
}

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
	// remove ncrypt header
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
