ncrypt is for your data, what a vault is for your bank. Valuables should be protected.

## Easy to use

```bash
# Encrypt using default key
$ ncrypt thedoc.docx
thedoc.docx successfully encrypted.

$ ncrypt thedoc.docx.safe
thedoc.docx successfully decrypted.

$ ncrypt -key='abracadabra' info.txt
info.txt successfully encrypted using key 'abracadabra'

$ ncrypt -key info.txt
info.txt successfully encrypted using generated-key '53JzhILXUERPYozNX4C/M2gfooGJL1D6UWA2yB3HgHM'

$ ncrypt share info.txt
info.txt available for download: https://ncrypt.io/EhiYm298 (TTL 24h)
```

## Super secure

Authenticated Encryption with Additional Authenticated Data (AEAD) couples confidentiality
and integrity. The most popular AEAD today, AES-GCM.

ref paper: https://eprint.iacr.org/2017/168.pdf

## Compliance

ncrypt stores the encryption keys in a `key` file, located in `$HOME/.config/ncrypt`

To comply with regulations you might need to generate encryption keys using a Hardware Security Module aka HSM. 
Ncrypt comes with a HSM security plugin for GCP and AWS. These providers offer HSM as a service. 

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
