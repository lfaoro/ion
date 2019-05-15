package main

import (
	"fmt"
	"os"

	"github.com/lfaoro/pkg/logger"
	"github.com/pkg/errors"
)

var log = logger.New("debug", nil)

func check(err error) {
	if err != nil {
		errors.WithStack(err)
		fmt.Printf("🔥 Error: %v\n", err)
		os.Exit(1)
	}
}
