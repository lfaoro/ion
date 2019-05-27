// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var downloadCmd = cli.Command{
	Name:    "download",
	Aliases: []string{"do", "dow", "down"},
	Usage:   "downloads the encrypted file using the reference-code.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, out, o",
			Usage: "output path / where to store the downloaded file",
		},
	},
	Action: func(c *cli.Context) error {
		fileName := c.Args().Get(0)
		output := c.String("output")

		err := downloadFile(fileName, output)
		return err
	},
}

func downloadFile(fileName, output string) error {
	// we append a len(6) string to every file
	if len(fileName) <= 6 {
		return errors.New("invalid ncrypt file")
	}

	if strings.Contains(fileName, "http"){
		fileName = path.Base(fileName)
	}

	uri, err := url.ParseRequestURI("https://storage.googleapis.com/"+bucketName)
	if err != nil {
		return err
	}

	uri.Path = path.Join(uri.Path, fileName)

	res, err := http.Get(uri.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return errors.Wrap(err, "download:")
	}

	if output == "" {
		output, _ = filepath.Abs(fileName[6:])
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0600)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", output)
	return err
}
