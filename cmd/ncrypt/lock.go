// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/urfave/cli"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/lfaoro/pkg/encrypto"
	"github.com/lfaoro/pkg/encrypto/aesgcm"
)

var lockCmd = cli.Command{
	Name:    "lock",
	Aliases: []string{"lo, loc"},
	Usage:   "locks the encryption key with a user provided password.",
	Action: func(c *cli.Context) error {
		fmt.Print("Encryption-key: ")
		pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Print("\n")

		home, _ := os.UserHomeDir()
		saltPath := path.Join(home, configPath, "salt")

		if _, err = os.Stat(saltPath); os.IsNotExist(err) {
			salt := encrypto.RandomBytes(16)
			err = ioutil.WriteFile(saltPath, salt, 0600)
			if err != nil {
				return err
			}
		}

		salt, err := ioutil.ReadFile(saltPath)
		if err != nil {
			return err
		}

		skey, err := scrypt.Key([]byte(pwd), salt, 32768, 8, 1, 32)
		if err != nil {
			return err
		}

		var engineKey *[32]byte
		copy(engineKey[:], skey)
		engine, err := aesgcm.New(engineKey)
		if err != nil {
			return err
		}

		return lockKey(engine)
	},
}

func lockKey(engine encrypto.Cryptor) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	keyFile := filepath.Join(home, configPath, "key")

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	// Check if file is already locked.
	if bytes.Contains(key, []byte(Header)) {
		info, _ := os.Stat(keyFile)
		fmt.Printf("Locked on %v\n", info.ModTime().Format("Mon Jan 2 15:04:05"))
		return nil
	}

	cipherText, err := engine.Encrypt(key)
	if err != nil {
		return err
	}

	cipherText = addHeader(cipherText)

	err = ioutil.WriteFile(keyFile, cipherText, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Locked %s\n", keyFile)
	return nil
}
