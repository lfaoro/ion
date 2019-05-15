package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/lfaoro/pkg/encrypto"
	"github.com/stretchr/testify/assert"
)

var testFile = "testfile.txt"
var fileMD5 = "1e50210a0202497fb79bc38b6ade6c34"

func Test_getSignedURL(t *testing.T) {

	u, err := getSignedURL(testFile, fileMD5)
	assert.Nil(t, err)
	assert.NotNil(t, u)
}

func Test_uploadToSignedURL(t *testing.T) {
	fileName := "genesis.txt.encrypted"
	filePath := filepath.Join("testdata", fileName)
	data, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	fmt.Println("data size:", len(data)/1024, " Kilobytes")

	hasher := md5.New()
	hasher.Write(data)
	md5hash := hex.EncodeToString(hasher.Sum(nil))

	rs := encrypto.RandomString(5)
	fileName = fmt.Sprintf("%v-%v", rs, fileName)
	u, err := getSignedURL(fileName, string(md5hash))
	assert.Nil(t, err)
	assert.NotNil(t, u)

	err = uploadToSignedURL(data, u)
	assert.Nil(t, err)
}
