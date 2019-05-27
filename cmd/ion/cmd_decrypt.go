// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

var decryptCmd = cli.Command{
	Name:    "decrypt",
	Aliases: []string{"d", "de", "dec"},
	Usage:   "encrypts a file",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "key",
			Usage:  "decrypts the file using the provided encryption key",
			EnvVar: "ION_KEY",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			return errors.New("which file?")
		}

		keyFlag := c.Bool("key")
		var key = new([32]byte)
		// if --key -- generate new key
		if keyFlag {
			r := bufio.NewReader(os.Stdin)
			fmt.Print("🔑 Encryption-key: ")
			encodedKey, err := r.ReadString('\n')
			if err != nil {
				return err
			}
			if len(encodedKey) < 32 {
				return errors.New("invalid encryption key")
			}
			k, err := base64.StdEncoding.DecodeString(encodedKey)
			copy(key[:], k)
		} else {
			// otherwise retrieve the key from the config
			k, err := keyFromConfig()
			if err != nil {
				return err
			}
			copy(key[:], k)
		}

		engine, err := aesgcm.New(key)
		if err != nil {
			return err
		}

		for _, path := range buildFileList(c.Args()) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to open file")
			}

			if !isEncrypted(data) {
				return errors.New("nothing to decrypt")
			}

			err = decryptFile(path, data, engine)
			if err != nil {
				return err
			}

			fmt.Printf("🔓 Decrypted %s\n", path)
		}

		return nil
	},
}
