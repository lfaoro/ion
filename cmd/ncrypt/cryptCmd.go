package main

import (
	"bytes"
	"fmt"
	"github.com/lfaoro/pkg/encrypto/aesgcm"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
)

func cryptCmd(c *cli.Context) error {
	return nil
}

func crypt(c *cli.Context, ce *aesgcm.AESGCM, fileName, filePath string, data []byte) error {
	if fileName == "" {
		return errors.New("file/s to encrypt not provided")
	}
	var backupFlag bool
	if c != nil {

		backupFlag = c.Bool("backup")
	}

	if isEncrypted(data) {
		// remove ncrypt header
		cipherText, err := removeHeader(data)
		if err != nil {
			return err
		}

		plainText, err := ce.Decrypt(cipherText)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filePath, plainText, 0600)
		if err != nil {
			return err
		}

		fmt.Printf("🔓 Decrypted %s\n", fileName)
		return nil
	} else {
		if backupFlag {
			tmp := os.TempDir()
			backup := filepath.Join(tmp, "ncrypt", fileName)
			err := ioutil.WriteFile(backup, data, 0600)
			if err != nil {
				return err
			}
			fmt.Printf("💾 Backed up %s", fileName)
		}

		cipherText, err := ce.Encrypt(data)
		if err != nil {
			return err
		}

		cipherText = addHeader(cipherText)

		err = ioutil.WriteFile(filePath, cipherText, 0600)
		if err != nil {
			return err
		}

		fmt.Printf("🔒 Encrypted %s\n", fileName)
		return nil
	}
}

func removeHeader(data []byte) ([]byte, error) {
	if !isEncrypted(data) {
		return []byte{}, errors.New("invalid Helix2 file")
	}
	i := bytes.IndexByte(data, byte('\n'))
	if i == -1 {
		return []byte{}, errors.New("invalid Helix2 file")
	}
	return data[i+1:], nil
}

func addHeader(data []byte) []byte {
	header := getHeader()
	return append(header, data...)
}

func newCryptoEngine(key string) (*aesgcm.AESGCM, error) {
	if key != "" {
		aes, err := aesgcm.New(string(key))
		if err != nil {
			return nil, err
		}
		return aes, nil
	}

	key, err := configKey()
	if err != nil {
		return nil, err
	}

	aes, err := aesgcm.New(string(key))
	if err != nil {
		return nil, err
	}

	return aes, nil
}

func isEncrypted(data []byte) bool {
	return bytes.Contains(data, getHeader())
}
