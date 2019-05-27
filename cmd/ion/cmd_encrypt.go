// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lfaoro/pkg/encrypto"
	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

var encryptCmd = cli.Command{
	Name:    "encrypt",
	Aliases: []string{"e", "en"},
	Usage:   "$ ion e --backup --key genesis.txt",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "key",
			Usage:  "encrypts the file using a generated encryption key",
			EnvVar: "ION_KEY",
		},
		cli.BoolFlag{
			Name:   "backup",
			Usage:  "backups the file before encryption",
			EnvVar: "ION_BACKUP",
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
			k := encrypto.NewEncryptionKey()
			encodedKey := base64.StdEncoding.EncodeToString(k[:])
			fmt.Printf("🔑 Encryption-key: %s\n", encodedKey)
			key = k
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

		backupFlag := c.Bool("backup")
		for _, path := range buildFileList(c.Args()) {
			backupFile(backupFlag, path)

			data, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to open file")
			}

			if isEncrypted(data) {
				return errors.Errorf("%s is already encrypted", path)
			}

			err = encryptFile(path, data, engine)
			if err != nil {
				return err
			}

			fmt.Printf("🔒 Encrypted %s\n", path)
		}

		return nil
	},
}
