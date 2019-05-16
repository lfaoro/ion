# WIP (work in progress) 
> Experimenting with the idea and its utility at this stage. 

🧬 ncrypt - a geeky & friendly way to simply encrypt locally & share.

Consumer grade CLI-app, designed for every-user with love for the power-user.

Encryption is done on your computer, your data does not hit the cloud unencrypted.

No logs except errors are being collected from [cmd/server](https://github.com/lfaoro/ncrypt/tree/master/cmd/server) -- check it.

Your data (in its ciphered form) lives for maximum 24 hours in a GCS Bucket. The retention policy is locked -- nobody 
can change it. Ref: https://cloud.google.com/storage/docs/bucket-lock 

[![pipeline status](https://gitlab.com/lfaoro/ncrypt/badges/master/pipeline.svg)](https://gitlab.com/lfaoro/ncrypt/commits/master)
[![coverage report](https://gitlab.com/lfaoro/ncrypt/badges/master/coverage.svg)](https://gitlab.com/lfaoro/ncrypt/commits/master)

## Quick start

```bash
# macOS (WIP)
brew install lfaoro/tap/ncrypt

# linux (WIP)
curl ncryp.to/i | sh

# developers
go get -u github.com/lfaoro/ncrypt/...
cd $GOPATH/src/github.com/lfaoro/ncrypt/cmd/ncrypt
make install
ncrypt -h
make test
```

## Easy to use

```bash
$ ncrypt genesis.doc
🔒 Encrypted genesis.doc

$ ncrypt genesis.doc
🔓 Decrypted genesis.doc

$ ncrypt upload genesis.doc
⬆️ Uploaded genesis.doc
ℹ️ Expires in 24 hours
ℹ️ Download reference: 2E3fde2a-genesis.doc

$ ncrypt download 2E3fde2a-genesis.doc
⬇️ Downloaded genesis.doc

$ ncrypt -key genesis.doc
🔑 Encryption-key: xy-TdOfXeQ5otTB0kXKLHbeYwpNCo0rn
🔒 Encrypted genesis.doc


# WIP commands

$ ncrypt genesis.doc
🧮 Unable to decrypt using your local key
🔑 Decryption-key: ***********
🔓 Decrypted genesis.doc

$ ncrypt lock 
Cryptovariable: **********
Locked .config/ncrypt/key

$ ncrypt unlock 
Cryptovariable: **********
Unlocked .config/ncrypt/key
```

## Leading encryption standard

Authenticated Encryption with Additional Authenticated Data (AEAD) couples confidentiality and integrity. Using the 
most popular AEAD today, AES-GCM.

ref paper: https://eprint.iacr.org/2017/168.pdf

## Compliance (WIP)

Right now ncrypt stores the encryption keys in a `key` file, located in `$HOME/.config/ncrypt` with `0600` permission
. Ideally we'll have the keys stored in macOS keychain -- although I don't know if there's something comparable for 
Linux and Windows.

To comply with regulators you might need to generate encryption keys using a Hardware Security Module aka HSM. 

ncrypt comes with a HSM plugin for GCP and AWS. These providers offer HSM as a service. 

Configure the GCP/AWS environment variables in order to activate Cloud HSM; ref: https://.

> In progress: https://github.com/lfaoro/ncrypt/issues/1

## Contributing

> Any help, feedback and suggestions are very welcome and greatly appreciated.
> Start by opening an [issue](https://github.com/lfaoro/pkg/issues/new).

## Motivation

It's hard to find a service one can completely trust -- everybody claims they're encrypting your data, although how 
can you be sure? 

I believe the only way trust what happens to your data is to see exactly the steps that lead to its manipulation, 
encryption & storage.

ncrypt is F/OSS -- anyone can check how data is being encrypted and handled, spot eventual issues and fix insecurities.

Designed with user-friendliness in mind, aspiring to be used also by non-dev users.