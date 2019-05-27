// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var name string

func Test_uploadFile(t *testing.T) {
	t.Skip()
	fileName := "genesis.txt.encrypted"

	// TODO(leo): add corrupted data test case.

	err := uploadFile(fileName)
	assert.Nil(t, err)
}

func Test_downloadFile(t *testing.T) {
	if name == "" {
		t.Skip("no file name provided")
	}

	tmpPath := path.Join(os.TempDir(), "ncrypt")
	err := os.MkdirAll(tmpPath, 0700)
	assert.Nil(t, err)
	tmpPath = path.Join(tmpPath, name)
	err = downloadFile(name, tmpPath)
	assert.Nil(t, err)
}
