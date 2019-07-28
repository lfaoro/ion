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
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
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

func downloadFile(filename, output string) error {
	// we append a len(6) string to every file
	if len(filename) <= 6 {
		return errors.New("invalid ncrypt file")
	}

	if strings.Contains(filename, "http") {
		filename = path.Base(filename)
	}

	uri, err := url.ParseRequestURI("https://storage.googleapis.com/" + bucketName)
	if err != nil {
		return err
	}

	uri.Path = path.Join(uri.Path, filename)

	res, err := http.Get(uri.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return errors.Wrap(err, "download:")
	}

	if output == "" {
		output, _ = filepath.Abs(filename[6:])
	}

	size, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	bar := pb.New(size)
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond * 50)
	bar.SetMaxWidth(80)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.Start()
	rd := bar.NewProxyReader(res.Body)

	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded: %s\n", output)
	return err
}
