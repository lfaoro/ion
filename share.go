package main

import (
  "fmt"

  "github.com/urfave/cli"
)

func shareCommand(c *cli.Context) error {
  // read the file
  // check the bytes for the NCRYPT tag / exit if not present
  // Can't take the responsibility to share unencrypted files
  // upload bytes to bucket w/ expiration date
  // print link to file
  fileName := c.Args().First()
  filePath := filePath(fileName)
  link := Upload(filePath)
  fmt.Printf("%s is available at: %s", fileName, link)
  return nil
}
