WIP: this project is a Work-In-Progress, we're experimenting with the idea and its utility.
It is not released, nor stable, nor reliable; there are no functional tests yet.

ncrypt is for your data, what a vault is for your bank. Valuables should be protected.

## Easy to use

```bash
$ ncrypt genesis.doc
[+] Encrypted genesis.doc

$ ncrypt genesis.doc
[+] Decrypted genesis.doc

$ ncrypt upload genesis.doc
[+] Uploaded genesis.doc
[#] Download reference: 2E3fde2a-genesis.doc
[#] Expires: 24 hours

$ ncrypt download 2E3fde2a-genesis.doc
[+] Downloaded genesis.doc

$ ncrypt -key genesis.doc
[+] Password: *******
[+] Encrypted genesis.doc

$ ncrypt genesis.doc
[+] Password: ******
[+] Decrypted genesis.doc
```

## Super secure

Authenticated Encryption with Additional Authenticated Data (AEAD) couples confidentiality
and integrity. Using the most popular AEAD today, AES-GCM.

ref paper: https://eprint.iacr.org/2017/168.pdf

## Compliance

ncrypt stores the encryption keys in a `key` file, located in `$HOME/.config/ncrypt`

To comply with regulations you might need to generate encryption keys using a Hardware Security Module aka HSM. 
ncrypt comes with a HSM security plugin for GCP and AWS. These providers offer HSM as a service. 

Configure the GCP/AWS environment variables in order to activate Cloud HSM; ref: https://.

> In progress: https://github.com/lfaoro/ncrypt/issues/1

## Quick start

```bash
## developers
go get -u github.com/lfaoro/ncrypt

## macOS
brew install ncrypt

## linux
curl ncryp.to/i | sh
```

## Contributing

> Any help and suggestions are very welcome and appreciated.
> Start by opening an [issue](https://github.com/lfaoro/pkg/issues/new).
