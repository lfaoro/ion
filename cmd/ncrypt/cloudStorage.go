// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/lfaoro/pkg/encrypto"
)

type gcs struct {
	secretPath string
	bucketName string
}

func (g gcs) Upload(filePath string) string {
	data := readFile(filePath)
	fileName := filepath.Base(filePath)

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	check(err)

	rs := encrypto.RandomString(5)
	check(err)

	objName := rs + "-" + filepath.Base(fileName)
	obj := c.Bucket(g.bucketName).Object(objName)
	wc := obj.NewWriter(ctx)
	_, err = io.Copy(wc, bytes.NewReader(data))
	check(err)
	wc.Close()
	err = obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	check(err)

	link := fmt.Sprintf("http://storage.googleapis.com/%s/%s", g.bucketName, objName)
	return link
}
