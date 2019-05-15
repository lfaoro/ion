// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encrypto

import (
	"crypto/sha256"
	"fmt"
	"log"
)

// Hash is a convenience function for sha256 hashing
// that returns base16 encoded data.
func Hash(data []byte) string {
	h := sha256.New()
	n, err := h.Write(data)
	if err != nil {
		log.Println("hashing failed:", err)
		return ""
	}
	if len(data) < n {
		log.Println("hashing failed: invalid data length")
		return ""
	}
	sum := h.Sum(nil)
	// %x - converts the data into base16
	return fmt.Sprintf("%x", sum)
}
