// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lfaoro/ncrypt/cmd/lsh/ccrypt"
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

	// check if file has ?key= pattern
	var key = new([32]byte)
	i := strings.Index(fileName, "?=")
	t := i > 0
	if t {
		// decode key
		k, err := base64.URLEncoding.DecodeString(fileName[i+2:])
		if err != nil {
			return err
		}
		copy(key[:], k)

		// remove key hash from filename
		fileName = fileName[:i]
		fmt.Println(fileName)
	}

	uri, err := url.ParseRequestURI("https://storage.googleapis.com/ncrypt-users")
	if err != nil {
		return err
	}

	uri.Path = path.Join(uri.Path, fileName)
	fmt.Println("ℹ️ Reference URL:", uri.String())

	res, err := http.Get(uri.String())
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return errors.Wrap(err, "download:")
	}

	ciphertext, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(ciphertext[:10]))

	var plaintext = []byte{}
	if t {
		pt, err := ccrypt.Decrypt(ciphertext, key)
		if err != nil {
			return err
		}
		plaintext = pt
	}

	if output == "" {
		output, _ = filepath.Abs(fileName[6:])
	}

	err = ioutil.WriteFile(output, plaintext, 0600)
	if err != nil {
		return err
	}

	fmt.Println("⬇️ Downloaded:", output)
	o, err := exec.Command("head", "-n2", output).CombinedOutput()
	fmt.Println("👀 Preview:\n", string(o))
	return err
}
