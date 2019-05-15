// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/lfaoro/pkg/encrypto"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var key string

func Test_cryptCmd(t *testing.T) {
	key = encrypto.RandomString(32)

	fileName := "genesis.txt"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	ce, err := newCryptoEngine(key)
	assert.Nil(t, err)
	err = crypt(nil, ce, fileName, filePath, data)
	assert.Nil(t, err)
}

func Test_decryptCmd(t *testing.T) {
	ce, err := newCryptoEngine(key)
	assert.Nil(t, err)

	fileName := "genesis.txt"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	err = crypt(nil, ce, fileName, filePath, data)
	assert.Nil(t, err)
}
