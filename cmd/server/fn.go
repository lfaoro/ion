package server

//go:generate gcloud config set project ncrypt
//go:generate gcloud functions deploy Stream --memory 128 --runtime go111 --trigger-http
//generate gcloud functions delete Stream

import (
	"cloud.google.com/go/storage"
	"encoding/gob"
	"go.opencensus.io/plugin/ochttp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	// BucketName is the name of the bucket.
	BucketName = "ncrypt-users"
	// SecretPath is the path where the system can find the secret json file.
	SecretPath = ""
	privateKey []byte
)

func init() {
	b, err := ioutil.ReadFile("./secret/ncrypt-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	privateKey = b
	log.SetPrefix("stream: ")
}

const (
	errHeaderFilename = "x-ncrypt-filename header missing"
	errHeaderMD5      = "x-ncrypt-md5 header missing"
)

// TODO(leo): add request limit
// TODO(leo): add file size limit?

// Stream .
// todo: alt+p/n -- triggers multicursor
// client: i’d like to upload a file with name X and md5 Y
// cloud func: ok, upload here (url is signed with `Content-MD5: xxx` header)
// client: uploads file with header
// object storage: verifies the header so you don’t get corrupted data
// https://godoc.org/cloud.google.com/go/storage#SignedURL
//
// Thank you Antti Kupila https://github.com/akupila for the suggestion.
func Stream(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fileName := r.Header.Get("x-ncrypt-filename")
		if fileName == "" {
			hError(w, errHeaderFilename, http.StatusBadRequest)
			return
		}
		//md5Hash := r.Header.Get("x-ncrypt-md5")
		//if md5Hash == "" || len(md5Hash) != 32{
		//	hError(w, errHeaderMD5, http.StatusBadRequest)
		//	return
		//}

		//m := base64.StdEncoding.EncodeToString([]byte(md5Hash))
		//s, _ := base64.StdEncoding.DecodeString(m)
		//log.Println(md5Hash, len(md5Hash), m, len(s))

		url, err := storage.SignedURL(BucketName, fileName, &storage.SignedURLOptions{
			GoogleAccessID: "cloudfunction@ncrypt.iam.gserviceaccount.com",
			PrivateKey:     privateKey,
			Method:         "PUT",
			Expires:        time.Now().Add(time.Minute * 10),
			//MD5:            base64.StdEncoding.EncodeToString([]byte(md5Hash)),
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
	traced := ochttp.Handler{
		Handler:          http.HandlerFunc(fn),
		IsPublicEndpoint: true,
	}
	traced.ServeHTTP(w, r)
}

func hError(w http.ResponseWriter, error string, code int) {
	log.Printf("%v: %v", code, error)
	http.Error(w, error, code)
}
