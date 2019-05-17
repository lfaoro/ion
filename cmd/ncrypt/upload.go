// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"

	"github.com/lfaoro/ncrypt/cmd/ncrypt/signedurl"
	"github.com/lfaoro/pkg/encrypto"
)

var uploadCmd = cli.Command{
	Name:    "upload",
	Aliases: []string{"u, up, upl"},
	Usage:   "uploads the encrypted file to cloud storage.",
	Action: func(c *cli.Context) error {
		fileName := c.Args().Get(0)
		if len(fileName) < 1 {
			return errors.New("what should we upload?")
		}

		path := constructPath(fileName)

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = uploadFile(fileName, data)
		return err
	},
}

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
	surl, err := signedurl.Get(uploadName, "")
	if err != nil {
		return "", err
	}

	err = signedurl.Upload(data, surl)
	if err != nil {
		return "", err
	}

	fmt.Println("⬆️ Uploaded   ", fileName)
	fmt.Println("ℹ️ Reference  ", uploadName)
	fmt.Println("ℹ️ Expiration ", "24 hours")

	return uploadName, nil
}
