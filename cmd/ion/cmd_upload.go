// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

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

		fpath, err := filepath.Abs(fileName)
		if err != nil {
			return err
		}

		err = uploadFile(fpath)
		return err
	},
}

func uploadFile(fpath string) error {
	var downloadURL = "http://s.apionic.com"

	rs := encrypto.RandomString(5)
	uploadName := fmt.Sprintf("%v_%v", rs, filepath.Base(fpath))
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

	_url := path.Join(downloadURL, uploadName)
	rawurl, err := url.Parse(_url)
	if err != nil {
		return err
	}

	fmt.Printf("\nDownload from: %s\n", rawurl.String())

	return nil
}
