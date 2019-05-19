// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encrypto

import "crypto/tls"

func DefaultTLSConfig() *tls.Config {
	return &tls.Config{
		// Avoids most of the memorably-named TLS attacks
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13, // Go 1.12 only
		// Causes servers to use Go's default cipher suite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients.
			// tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		// Only use curves which have assembly implementations.
		CurvePreferences: []tls.CurveID{
			tls.CurveP256, // ECDHE
			tls.X25519,    // Go 1.8 only
		},
	}
}
