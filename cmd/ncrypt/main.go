// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lfaoro/pkg/encrypto"
	"github.com/lfaoro/pkg/logger"
)

var log = logger.New("debug", nil)

var (
	// Header is what we append to encrypted files, in order to launch
	// an encrypt or decrypt operation.
	Header = "## ncrypted with love"

	// version is injected during the release process.
	version = "dev"
	// commit is injected during the release process.
	commit = "none"
	// date is injected during the release process.
	date = "unknown"
)

const configPath = ".config/ncrypt"

func getHeader() []byte {
	header := fmt.Sprintf("%s\n", Header)
	return []byte(header)
}

func main() {
	app := cli.NewApp()
	app.Name = "ncrypt"
	app.Usage = "a geeky & friendly way to simply encrypt locally & share"
	app.Version = fmt.Sprintf("%s %s %s", version, commit, date)
	app.EnableBashCompletion = true

	app.Before = func(c *cli.Context) error {
		return checkConfig()
	}

	// Global flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "key",
			Usage:  "encrypts the file using the provided encryption key",
			EnvVar: "NCRYPT_KEY",
		},
		cli.BoolFlag{
			Name:   "backup",
			Usage:  "backups the file before encryption",
			EnvVar: "NCRYPT_BACKUP",
		},
	}

	app.Commands = []cli.Command{
		uploadCmd,
		downloadCmd,
		lockCmd,
		unlockCmd,
	}

	app.Action = func(c *cli.Context) error {
		keyFlag := c.Bool("key")
		backupFlag := c.Bool("backup")

		var key []byte
		if keyFlag {
			key, _ = encrypto.RandomBytes(32)
			fmt.Println("🔑 Encryption-key: ", string(key))
		}
		engine, err := newCryptoEngine(key)
		if err != nil {
			return err
		}

		for _, fileName := range c.Args() {
			if fileName == "" {
				return errors.New("file/s to encrypt not provided")
			}

			info, err := os.Stat(fileName)
			if err != nil {
				e := err.(*os.PathError)
				return errors.Errorf("%s - %s", e.Path, e.Err.Error())
			}
			if info.IsDir() {
				fmt.Printf("Sorry, I can't upload a whole directory. \nYou can tar it first though:\n\n")
				fmt.Printf("$ tar -cvf %s.tar %s/\n", info.Name(), info.Name())
				fmt.Printf("$ ncrypt %s.tar\n", info.Name())
				os.Exit(1)
			}

			path := constructPath(fileName)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to open file")
			}

			if backupFlag {
				tmp := os.TempDir()
				backup := filepath.Join(tmp, "ncrypt", fileName)
				err := os.MkdirAll(filepath.Dir(backup), 0700)
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(backup, data, 0600)
				if err != nil {
					return err
				}
				fmt.Printf("💾 Backed up %s", fileName)
			}

			err = cryptoCmd(engine, path, data)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "🔥 Error: %v\n", err)
		os.Exit(1)
	}
}

func checkConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	keyFile := filepath.Join(home, ".config/ncrypt/key")

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(keyFile), 0700)
		if err != nil {
			return errors.Wrap(err, "config")
		}
		_, err := os.Create(keyFile)
		if err != nil {
			return errors.Wrap(err, "config")
		}
		err = os.Chmod(keyFile, 0600)
		if err != nil {
			return errors.Wrap(err, "chmod")
		}
	}

	n, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return errors.Wrap(err, "config")
	}

	if len(n) < 2 {
		key := encrypto.RandomString(32)
		err = ioutil.WriteFile(keyFile, []byte(key), 0600)
		if err != nil {
			return errors.Wrap(err, "config")
		}
	}

	return nil
}

func keyFromConfig() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		check(err)
	}

	keyFile := filepath.Join(home, configPath, "key")

	key, err := ioutil.ReadFile(keyFile)
	check(err)

	return key, nil
}

func constructPath(fileName string) string {
	// if the fileName contains more than 1 '/' char
	// we can assume it's a full canonical path.
	sub := strings.Split(fileName, "/")
	if len(sub) > 1 {
		return fileName
	}

	// obtain canonical path
	wd, err := os.Getwd()
	check(err)
	return filepath.Join(wd, fileName)
}

func check(err error) {
	if err != nil {
		fmt.Printf("🔥 Error: %v\n", err)
		os.Exit(1)
	}
}
