// Package hsm provides Hardware Security Module connections.
package hsm

type HSM interface {
	Connect()
	Encrypt()
	Decrypt()
}
