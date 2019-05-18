// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var downloadCmd = cli.Command{
	Name:    "download",
	Aliases: []string{"d, do, down"},
	Usage:   "downloads the encrypted file using the reference-code.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, out, o",
			Usage: "output path / where to store the downloaded file",
		},
	},
	Action: func(c *cli.Context) error {
		fileName := c.Args().Get(0)
		err := downloadFile(fileName)
		return err
	},
}

func downloadFile(fileName string) error {
	// we append a len(6) string to every file
	if len(fileName) <= 6 {
		return errors.New("invalid ncrypt file")
	}

	uri, err := url.ParseRequestURI("https://storage.googleapis.com/ncrypt-users")
	if err != nil {
		return err
	}

	uri.Path = path.Join(uri.Path, fileName)
	fmt.Println("ℹ️ Reference URL ", uri.String())

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

	if res.StatusCode > 299 {
		return errors.Wrap(err, "download:")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !isEncrypted(b) {
		return errors.New("downloaded unecrypted content, contact: lfaoro+support@gmail.com")
	}

	fn := fileName[6:]
	err = ioutil.WriteFile(fn, b, 0600)
	if err != nil {
		return err
	}

	fmt.Println("⬇️ Downloaded ", fn)

	return nil
}
