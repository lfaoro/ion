// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logger is the way we do logging.
package logger

import (
	"log"
	"os"

	"cloud.google.com/go/errorreporting"
)

// Log is a wrapper on the stdlib log pkg which includes the Stackdriver
// Error Reporting client.
type Log struct {
	*log.Logger
	*errorreporting.Client
}

// New returns an initialized Log with defaults setup.
func New(prefix string, errorClient *errorreporting.Client) *Log {
	return &Log{
		log.New(os.Stdout, prefix+" ", log.Lshortfile),
		errorClient,
	}
}
