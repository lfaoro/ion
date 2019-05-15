package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lfaoro/pkg/encrypto"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	// Header is what we append to encrypted files, in order to launch
	// an encrypt or decrypt operation.
	Header = "## ncrypt"
	// Version is injected during the release process.
	Version = "v1.0.0-alpha"
)

func getHeader() []byte {
	header := fmt.Sprintf("%s %s\n", Header, Version)
	return []byte(header)
}

func main() {

	app := cli.NewApp()
	app.Name = "ncrypt"
	app.Usage = "a geeky & friendly way to simply encrypt locally & share"
	app.Version = Version
	app.EnableBashCompletion = true

	app.Before = func(c *cli.Context) error {
		return checkConfig()
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "key,k",
			Usage: "encrypts the file using the provided encryption key",
		},
		cli.BoolFlag{
			Name:  "backup,b",
			Usage: "backups the file before encryption",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "upload",
			ShortName: "u",
			Aliases:   []string{"u,up"},
			Usage:     "uploads the encrypted file to cloud storage.",
			Action: func(c *cli.Context) error {
				fileName := c.Args().Get(0)
				if len(fileName) < 1 {
					return errors.New("what should we upload?")
				}
				path := filePath(fileName)
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				_, err = uploadFile(fileName, data)
				return err
			},
		},
		{
			Name:      "download",
			ShortName: "d",
			Aliases:   []string{"d, do, down"},
			Usage:     "downloads the encrypted file using the reference-code.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output,out,o",
					Usage: "output path / where to store the downloaded file",
				},
			},
			Action: func(c *cli.Context) error {
				fileName := c.Args().Get(0)
				err := downloadCmd(fileName)
				return err
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		// key flag
		var keyFlag = c.Bool("key")

		var key string
		if keyFlag {
			rs := encrypto.RandomString(32)
			fmt.Println("🔑 Encryption-key: ", rs)
			key = rs
		}
		ce, err := newCryptoEngine(key)
		if err != nil {
			return err
		}

		for _, fileName := range c.Args() {
			// TODO: add check to identify if the fileName is a file
			// if not a valid file fail with error
			// Please, specify a file name.
			path := filePath(fileName)
			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to open file")
			}
			// VC (Visual Cue): Action
			err = crypt(c, ce, fileName, path, fileData)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		check(err)
	}

}

// TODO: refactor
func checkConfig() error {
	keyFile := keyPath()
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
	f, err := os.OpenFile(keyFile, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return errors.Wrap(err, "config")
	}
	defer f.Close()
	n, err := ioutil.ReadFile(keyFile)
	if len(n) < 2 {
		key := encrypto.RandomString(32)
		if err != nil {
			return errors.Wrap(err, "config")
		}
		_, err = f.WriteString(key)
	}
	return err
}

func keyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		check(err)
	}
	return filepath.Join(home, ".config/ncrypt/key")
}

func configKey() (string, error) {
	keyFile := keyPath()
	key := readFile(keyFile)
	return string(key), nil
}

func filePath(fileName string) string {
	//TODO: if path contains more than 1 '/' return
	wd, err := os.Getwd()
	check(err)
	return filepath.Join(wd, fileName)
}

func readFile(filePath string) (data []byte) {
	b, err := ioutil.ReadFile(filePath)
	check(err)
	return b
}

func dataFromFile(f *os.File) ([]byte, error) {
	var data []byte
	_, err := f.Read(data)
	if err != nil {
		return []byte{}, err
	}
	if len(data) <= 0 {
		return []byte{}, errors.New("file is empty")
	}
	return data, nil
}
