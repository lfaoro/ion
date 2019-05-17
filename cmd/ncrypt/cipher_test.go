// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/lfaoro/pkg/encrypto"

	"github.com/stretchr/testify/assert"
)

var testKey []byte

func Test_cryptCmd(t *testing.T) {
	key := encrypto.RandomString(32)
	t.Log(key)
	testKey = []byte(key)

	fileName := "genesis.txt"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	engine, err := newCryptoEngine(testKey)
	assert.Nil(t, err)
	err = cryptoCmd(engine, filePath, data)
	assert.Nil(t, err)
}

func Test_decryptCmd(t *testing.T) {
	t.Log(string(testKey))
	engine, err := newCryptoEngine(testKey)
	assert.Nil(t, err)

	fileName := "genesis.txt"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	err = cryptoCmd(engine, filePath, data)
	assert.Nil(t, err)
}
