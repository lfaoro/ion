// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/lfaoro/pkg/encrypto"
	"github.com/lfaoro/pkg/encrypto/aesgcm"

	"github.com/stretchr/testify/assert"
)

var testKey *[32]byte

func Test_cryptCmd(t *testing.T) {
	key := encrypto.NewEncryptionKey()
	testKey = key

	fileName := "genesis.txt"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	engine, err := aesgcm.New(testKey)
	assert.Nil(t, err)
	err = cryptoCmd(engine, filePath, data)
	assert.Nil(t, err)
}

func Test_decryptCmd(t *testing.T) {
	engine, err := aesgcm.New(testKey)
	assert.Nil(t, err)

	fileName := "genesis.txt"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	err = cryptoCmd(engine, filePath, data)
	assert.Nil(t, err)
}
