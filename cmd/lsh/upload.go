// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"

	"github.com/lfaoro/ncrypt/cmd/lsh/ccrypt"
	"github.com/lfaoro/ncrypt/cmd/lsh/signedurl"
	"github.com/lfaoro/pkg/encrypto"
)

var uploadCmd = cli.Command{
	Name:    "upload",
	Aliases: []string{"u, up, upl"},
	Usage:   "uploads the encrypted file to cloud storage.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "key",
			Usage: "generate a shareable encryption key",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "send download link via email",
		},
	},
	Action: func(c *cli.Context) error {
		fileName := c.Args().Get(0)
		if len(fileName) < 1 {
			return errors.New("what should we upload?")
		}

		err := uploadFile(fileName, c.Bool("key"))
		return err
	},
}

func uploadFile(fileName string, encrypt bool) error {
	path := constructPath(fileName)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var b64_key string
	if encrypt {
		k := encrypto.NewEncryptionKey()

		ciphertext, err := ccrypt.Encrypt(data, k)
		if err != nil {
			return err
		}

		data = ciphertext

		b64_key = base64.URLEncoding.EncodeToString(k[:])
	}

	rs := encrypto.RandomString(5)
	uploadName := fmt.Sprintf("%v_%v", rs, fileName)
	surl, err := signedurl.Get(uploadName, "")
	if err != nil {
		return err
	}

	err = signedurl.Upload(data, surl)
	if err != nil {
		return err
	}

	fmt.Println("⬆️ Uploaded   ", fileName)
	if encrypt {
		fmt.Println("ℹ️ Reference  ", uploadName+"?="+b64_key)
	} else {
		fmt.Println("ℹ️ Reference  ", uploadName)
	}
	fmt.Println("ℹ️ Expiration ", "24 hours")

	return nil
}
