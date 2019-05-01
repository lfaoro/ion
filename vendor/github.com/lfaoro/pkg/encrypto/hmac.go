// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encrypto

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Hmac512 ciphers the data extracted from the Reader and returns
// a b64 encoded string of a SHA512 hash.
func Hmac512(key string, data io.Reader) (b64 string, err error) {
	hash := hmac.New(sha512.New, []byte(key))
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return "", errors.Wrap(err, "unable to read data from reader")
	}
	if _, err := hash.Write(b); err != nil {
		return "", errors.Wrap(err, "unable to hash data")
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
