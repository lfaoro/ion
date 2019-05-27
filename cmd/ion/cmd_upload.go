// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli"

	"github.com/lfaoro/ncrypt/internal/gcs"
	"github.com/lfaoro/pkg/encrypto"
)

var uploadCmd = cli.Command{
	Name:    "upload",
	Aliases: []string{"u", "up", "upl"},
	Usage:   "uploads the encrypted file to cloud storage.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "to",
			Usage: "send a download link via email",
		},
	},
	Action: func(c *cli.Context) error {
		fileName := c.Args().Get(0)
		if len(fileName) < 1 {
			return errors.New("what should we upload?")
		}

		err := uploadFile(fileName)
		return err
	},
}

func uploadFile(fileName string) error {
	var downloadURL = "http://s.apionic.com"

	fpath := guessPath(fileName)
	log.Println("guessed path", fpath)

	rs := encrypto.RandomString(5)
	uploadName := fmt.Sprintf("%v_%v", rs, fileName)
	surl, err := gcs.GetSignedURL(uploadName, "")
	if err != nil {
		return err
	}

	data, err := os.Open(fpath)
	if err != nil {
		return err
	}
	err = gcs.UploadToSignedURL(surl, data)
	if err != nil {
		return err
	}

	url := path.Join(downloadURL, uploadName)

	fmt.Printf("\nDownload from: %s\n", url)

	return nil
}
