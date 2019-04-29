package main

import (
  "bytes"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"

  "github.com/lfaoro/pkg/encrypto"
  "github.com/lfaoro/pkg/encrypto/aesgcm"
  "github.com/pkg/errors"
  "github.com/urfave/cli"
)

func main() {

  app := cli.NewApp()
  app.Name = "ncrypt"
  app.Usage = "is for your data, what a vault is for your bank."
  app.Version = "1.0.0"
  app.EnableBashCompletion = true
  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name:  "key,k",
      Usage: "encrypt using a custom encryption key",
    },
    cli.BoolFlag{
      Name:  "share,s",
      Usage: "upload encrypted file",
    },
  }

  app.Commands = []cli.Command{
    {
      Name:    "share",
      Aliases: []string{"s"},
      Usage:   "shares the encrypted files.",
      Action: func(c *cli.Context) error {
        return shareCommand(c)
      },
    },
  }

  app.Action = func(c *cli.Context) error {

    err := checkConfig()
    check(err)

    aes := newAESCrypter(c.String("key"))

    // loop the args
    // just get the first arg for now
    // open the filePath
    // run the bytes through encryption manipulation
    // get the encrypted bytes
    // put them back into a filePath adding .ncrypt extension
    first := c.Args().Get(0)
    b := readFile(first)
    filePath := filePath(first)

    // decrypt if it's an ncrypt filePath.
    if bytes.Contains(b, []byte("## NCRYPT")) {
      i := bytes.IndexByte(b, byte('\n'))
      if i == -1 {
        fmt.Println("damn")
      }
      cipherText := b[i+1:] // all the bytes till /n +1(/n)
      plainText, err := aes.Decrypt(cipherText)
      check(err)
      fmt.Println(string(plainText))
      os.Exit(0)
    }

    ct, err := aes.Encrypt(b)
    check(err)
    newFile := filePath + ".ncrypt"
    f, err := os.Create(newFile);
    defer f.Close()
    check(err)
    // make string a global var
    // add program version to string
    _, err = io.Copy(f, bytes.NewReader([]byte("## NCRYPT v0.1.0\n")))
    _, err = io.Copy(f, bytes.NewReader(ct))
    check(err)

    fileName := filepath.Base(filePath)
    newFileName := filepath.Base(newFile)

    fmt.Println(fileName + " -> " + newFileName)

    return nil

  }

  err := app.Run(os.Args)
  if err != nil {
    fmt.Println(err)
  }

}

func filePath(fileName string) string {
  wd, err := os.Getwd()
  check(err)
  return filepath.Join(wd, fileName)
}

func readFile(filePath string) (data []byte) {
  b, err := ioutil.ReadFile(filePath)
  check(err)
  return b
}

func newAESCrypter(key string) encrypto.Cryptor {
  if key != "" {
    aes, err := aesgcm.New(string(key))
    check(err)
    return aes
  }
  keyPath := getKeyPath()
  fileKey, err := ioutil.ReadFile(keyPath)
  fmt.Println("using key:", string(fileKey))
  check(err)
  aes, err := aesgcm.New(string(fileKey))
  check(err)
  return aes
}

func checkConfig() error {
  keyFile := getKeyPath()
  if _, err := os.Stat(keyFile); os.IsNotExist(err) {
    err = os.MkdirAll(filepath.Dir(keyFile), 0700)
    if err != nil {
      return errors.Wrap(err, "config")
    }
    _, err := os.Create(keyFile)
    if err != nil {
      return errors.Wrap(err, "config")
    }
  }
  f, err := os.OpenFile(keyFile, os.O_RDWR, os.ModeAppend)
  if err != nil {
    return errors.Wrap(err, "config")
  }
  defer f.Close()
  n, err := ioutil.ReadFile(keyFile)
  if len(n) < 2 {
    fmt.Println("debug: generating secret key")
    key, err := encrypto.RandomString(32)
    if err != nil {
      return errors.Wrap(err, "config")
    }
    _, err = f.WriteString(key)
  }
  return err
}

func getKeyPath() string {
  home, err := os.UserHomeDir()
  if err != nil {
    fmt.Println(err)
  }
  keyFile := filepath.Join(home, ".config/ncrypt/key")
  return keyFile
}

func check(err error) {
  if err != nil {
    fmt.Errorf("[-] Error: %v", err)
    os.Exit(1)
  }
}
