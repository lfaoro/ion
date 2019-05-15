// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/lfaoro/pkg/encrypto"
)

func uploadFile(fileName string, data []byte) (string, error) {
	// bail if file is not encrypted
	if !isEncrypted(data) {
		// generate encryption key
		// encrypt
		// upload
		// display key
		return "", errors.New("won't upload unencrypted data")
	}

	rs := encrypto.RandomString(5)
	uploadName := fmt.Sprintf("%v_%v", rs, fileName)
	surl, err := getSignedURL(uploadName, "")
	if err != nil {
		return "", err
	}

	err = uploadToSignedURL(data, surl)
	if err != nil {
		return "", err
	}

	fmt.Println("⬆️ Uploaded   ", fileName)
	fmt.Println("ℹ️ Reference  ", uploadName)
	fmt.Println("ℹ️ Expiration ", "24 hours")

	return uploadName, nil
}

func downloadCmd(fileName string) error {
	if len(fileName) <= 6 {
		return errors.New("invalid ncrypt file")
	}

	uri, err := url.ParseRequestURI("https://storage.cloud.google.com/ncrypt-users")
	if err != nil {
		return err
	}

	uri.Path = path.Join(fileName)
	// fmt.Println("ℹ️ Reference URL ", uri.String())

	res, err := http.Get(uri.String())
	if err != nil {
		return err
	}

	// avoid writing the file during unit tests.
	if strings.HasSuffix(fileName, ".encrypted") {
		tmpPath := path.Join(os.TempDir(), "ncrypt", fileName)
		err = os.MkdirAll(tmpPath, 0700)
		if err != nil {
			return err
		}
		fileName = path.Join(tmpPath, fileName)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, b, 0600)
	if err != nil {
		return err
	}

	fn := fileName[6:]
	fmt.Println("⬇️ Downloaded ", fn)

	return nil
}
