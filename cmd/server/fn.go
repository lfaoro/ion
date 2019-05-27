// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

//go:generate gcloud config set project ncrypt
//go:generate gcloud functions deploy Stream --memory 128 --runtime go111 --trigger-http
// generate gcloud functions delete Stream

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
)

var (
	// BucketName is the name of the bucket.
	BucketName = "s.apionic.com"
	// SecretPath is the path where the system can find the secret json file.
	SecretPath = ""
	privateKey []byte
)

// TODO: use env variable instead of pem file.
// TODO: deploy automatically via gitlab-ci in order to guarantee the code is
//  not manipulated before reaching the server.

func init() {
	b, err := ioutil.ReadFile("./secret/ncrypt-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	privateKey = b

	log.SetPrefix("stream: ")
}

const (
	errHeaderFilename = "x-lsh-filename header missing"
	errHeaderMD5      = "x-lsh-md5 header missing"
)

// TODO(leo): add request limit
// TODO(leo): add file size limit?

// Stream is a proxy in front of GCS that returns signed URLs.
// https://godoc.org/cloud.google.com/go/storage#SignedURL
//
// Thank you Antti Kupila https://github.com/akupila for the suggestion.
func Stream(w http.ResponseWriter, r *http.Request) {
	fileName := r.Header.Get("x-ion-filename")
	if fileName == "" {
		hError(w, errHeaderFilename, http.StatusBadRequest)
		return
	}

	// TODO: add file integrity check
	// md5Hash := r.Header.Get("x-ncrypt-md5")
	// if md5Hash == "" || len(md5Hash) != 32{
	// 	hError(w, errHeaderMD5, http.StatusBadRequest)
	// 	return
	// }

	// m := base64.StdEncoding.EncodeToString([]byte(md5Hash))
	// s, _ := base64.StdEncoding.DecodeString(m)
	// log.Println(md5Hash, len(md5Hash), m, len(s))

	url, err := storage.SignedURL(BucketName, fileName, &storage.SignedURLOptions{
		GoogleAccessID: "cloudfunction@ncrypt.iam.gserviceaccount.com",
		PrivateKey:     privateKey,
		Method:         "PUT",
		Expires:        time.Now().Add(time.Minute * 10),
		// MD5:            base64.StdEncoding.EncodeToString([]byte(md5Hash)),
	})
	if err != nil {
		hError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = gob.NewEncoder(w).Encode(url)
	if err != nil {
		hError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func hError(w http.ResponseWriter, error string, code int) {
	log.Printf("%v: %v", code, error)
	http.Error(w, error, code)
}
