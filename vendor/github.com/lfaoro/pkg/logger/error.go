// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"net/http"
	"path/filepath"
	"runtime"

	"cloud.google.com/go/errorreporting"
)

// FatalIfErr panics if the error value is not nil.
func (log *Log) FatalIfErr(err error) {
	if err == nil {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		log.Fatal(err)
	}
	log.Fatalf("%v\n %v:%v", err, filepath.Base(file), line)
}

// LogIfErr logs if the error value is not nil.
func (log *Log) LogIfErr(err error) {
	if err == nil {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		log.Println(err)
	}
	log.Printf("%v\n %v:%v", err, filepath.Base(file), line)
}

// ReportIfErr sends an error report to Stackdriver Error Reporting.
func (log *Log) ReportIfErr(err error, user string, req *http.Request) {
	if err == nil {
		return
	}
	log.Report(errorreporting.Entry{
		Error: err,
		User:  user,
		Req:   req,
	})
	log.Printf("reported error: %v", err)
}
