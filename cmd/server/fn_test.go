package server

import (
	"bytes"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const streamURL = "https://us-central1-ncrypt.cloudfunctions.net/Stream"

var fileName = "testfile.txt"
var fileMD5 = "1e50210a0202497fb79bc38b6ade6c34"

func TestStream(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", Stream)

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	req.Header.Add("x-ncrypt-filename", fileName)
	req.Header.Add("x-ncrypt-md5", fileMD5)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, res.Body.String())
	assert.NotNil(t, res.Body.String())

	var sURL string
	err = gob.NewDecoder(bytes.NewReader(res.Body.Bytes())).Decode(&sURL)
	assert.Nil(t, err)
	u, err := url.Parse(sURL)
	assert.Nil(t, err)
	assert.NotNil(t, u)

	t.Log(u)
}
