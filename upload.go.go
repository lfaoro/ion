package main

import (
  "bytes"
  "context"
  "fmt"
  "io"
  "log"
  "os"
  "path/filepath"

  "cloud.google.com/go/storage"
  "github.com/lfaoro/pkg/encrypto"
)

type gcs struct {
  secretPath string
  bucketName string
}

func newGCS(secretPath, bucketName string) gcs{
  if secretPath == "" || bucketName == "" {

  }
  return gcs{
    secretPath: "",
    bucketName: "",
  }
}

func Upload(filePath string) string {
  data := readFile(filePath)
  fileName := filepath.Base(filePath)

  const bucketName = "ncrypt-users"
  err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "secret.json")
  _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
  if !ok {
    log.Fatal("GOOGLE_APPLICATION_CREDENTIALS required to access Hardware Security Module")
  }
  log.Println("successfully set GOOGLE_APPLICATION_CREDENTIALS")
  check(err)
  ctx := context.Background()

  c, err := storage.NewClient(ctx)
  check(err)

  rs, err := encrypto.RandomString(5)
  check(err)
  objName := rs+"-"+filepath.Base(fileName)
  obj := c.Bucket(bucketName).Object(objName)
  wc := obj.NewWriter(ctx)
  _, err = io.Copy(wc, bytes.NewReader(data))
  check(err)
  wc.Close()
  err = obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
  check(err)

  link := fmt.Sprintf("http://storage.googleapis.com/%s/%s", bucketName, objName)
  return link
}

