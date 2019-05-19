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

var unlockCmd = cli.Command{
	Name:    "unlock",
	Aliases: []string{"un, unl, unlo"},
	Usage:   "WIP: unlocks the encryption key.",
	Action: func(c *cli.Context) error {
		fmt.Print("Encryption-key: ")
		pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Print("\n")

		home, _ := os.UserHomeDir()
		saltPath := path.Join(home, configPath, "salt")

		salt, err := ioutil.ReadFile(saltPath)
		if err != nil {
			return err
		}

		skey, err := scrypt.Key([]byte(pwd), salt, 32768, 8, 1, 32)
		if err != nil {
			return err
		}

		var engineKey = new([32]byte)
		copy(engineKey[:], skey)
		engine, err := aesgcm.New(engineKey)
		if err != nil {
			return err
		}

		return unlockKey(engine)
	},
}

func unlockKey(engine encrypto.Cryptor) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	keyFile := filepath.Join(home, configPath, "key")

	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	// Check if file is already locked.
	if !bytes.Contains(data, []byte(Header)) {
		info, _ := os.Stat(keyFile)
		fmt.Printf("Unlocked on %v\n", info.ModTime().Format("Mon Jan 2 15:04:05"))
		return nil
	}

	data, err = removeHeader(data)
	if err != nil {
		return err
	}

	plainText, err := engine.Decrypt(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(keyFile, plainText, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Unlocked %s\n", keyFile)
	return nil
}
