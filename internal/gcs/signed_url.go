// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gcs provides abstractions to obtain a signed URL and upload
// data to it.
package gcs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

var c = http.Client{
	Timeout: time.Hour,
}

const streamURL = "https://us-central1-ncrypt.cloudfunctions.net/Stream"

func GetSignedURL(fileName, fileMD5 string) (string, error) {
	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("x-ion-filename", fileName)
	// req.Header.Add("x-ion-md5", fileMD5)

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		return "", fmt.Errorf(string(b))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var sURL string
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&sURL)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(sURL)
	if err != nil {
		return "", errors.New("unable to parse signed URL")
	}
	return u.String(), nil
}

func UploadToSignedURL(url string, data *os.File) error {
	info, _ := data.Stat()
	if info.Size() < 1 {
		return errors.New("no data to upload")
	}

	header := []byte("## ionized with love\n")
	buf := make([]byte, len(header))
	_, err := io.ReadFull(data, buf)
	if err != nil {
		return err
	}
	if !bytes.Equal(header, buf) {
		fmt.Printf("WARNING: %s is not encrypted\n", data.Name())
	}
	data.Seek(0, 0)

	bar := pb.New(int(info.Size()))
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond * 50)
	bar.SetMaxWidth(80)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.Start()
	rd := bar.NewProxyReader(data)

	req, err := http.NewRequest("PUT", url, rd)
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		b, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(b))
		os.Exit(1)
	}

	return nil
}
