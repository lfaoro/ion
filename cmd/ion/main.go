// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lfaoro/pkg/encrypto"
)

var (
	// Header is what we append to encrypted files, in order to launch
	// an encrypt or decrypt operation.
	Header = "## ionized with love"

	// version is injected during the release process.
	version = "dev"
	// commit is injected during the release process.
	commit = "none"
	// date is injected during the release process.
	date = "unknown"
)

const configPath = ".config/ion"

var configKey = path.Join(configPath, "key")

const bucketName = "s.apionic.com"

func getHeader() []byte {
	header := fmt.Sprintf("%s\n", Header)
	return []byte(header)
}

func main() {
	app := cli.NewApp()
	app.Name = "ion"
	app.Usage = "a geeky & friendly way to simply encrypt & share"
	app.Version = fmt.Sprintf("%s %s %s", version, commit, date)
	app.EnableBashCompletion = true
	app.Authors = []cli.Author{
		{
			Name:  "Leonardo Faoro",
			Email: "lfaoro@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "displays debug information.",
			EnvVar: "ION_DEBUG",
			Hidden: true,
		},
	}

	app.Before = func(c *cli.Context) error {
		return checkConfig()
	}

	app.Commands = []cli.Command{
		encryptCmd,
		decryptCmd,
		uploadCmd,
		downloadCmd,
		lockCmd,
		unlockCmd,
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, 2)
		return nil
	}

	err := app.Run(os.Args)
	exitIfError(err)
}

func checkConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	keyFile := filepath.Join(home, configKey)

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
		return nil, err
	}

	keyFile := filepath.Join(home, configPath, "key")

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	if isEncrypted(key) {
		return nil, errors.New("the encryption key is Locked, use:\n$ ncrypt unlock")
	}

	return key, nil
}

func guessPath(fp string) string {
	// if the path begins with a '/' character
	// we can assume it's a canonical path.
	if strings.HasPrefix(fp, "/") {
		return fp
	}

	// otherwise we retrieve the working directory.
	wd, err := os.Getwd()
	exitIfError(err)
	return filepath.Join(wd, fp)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func buildFileList(args cli.Args) []string {
	var files []string
	for _, fileName := range args {
		if fileName == "" {
			exitIfError(errors.New("file/s to encrypt not provided"))
		}

		info, err := os.Stat(fileName)
		if err != nil {
			e := err.(*os.PathError)
			exitIfError(errors.Errorf("%s - %s", e.Path, e.Err.Error()))
		}
		if info.IsDir() {
			// TODO: add directory taring
			// Do you wish to upload the whole directory? Y/n
			// tar directory -- encrypt tar data -- upload
			fmt.Printf("Sorry, I can't upload a whole directory. \nYou can tar it first though:\n\n")
			fmt.Printf("$ tar -cvf %s.tar %s/\n", info.Name(), info.Name())
			fmt.Printf("$ ncrypt %s.tar\n", info.Name())
			os.Exit(1)
		}

		fp := guessPath(fileName)
		files = append(files, fp)
	}
	return files
}

func backupFile(yes bool, fp string) {
	if !yes {
		return
	}
	filename := path.Base(fp)
	tmp := os.TempDir()
	backup := filepath.Join(tmp, "ion", filename)
	err := os.MkdirAll(filepath.Dir(backup), 0700)
	exitIfError(err)

	data, err := ioutil.ReadFile(fp)
	exitIfError(err)

	err = ioutil.WriteFile(backup, data, 0600)
	exitIfError(err)

	fmt.Printf("💾 Backed up %s", filename)
}
