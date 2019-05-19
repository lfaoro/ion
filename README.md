# 🧬 lsh - upload and share large data objects

> End to end encrypted, if you want...

Encryption is done on your computer, your data does not hit the cloud unencrypted.

No logs except errors are being collected from [cmd/server](https://github.com/lfaoro/lsh/tree/master/cmd/server) -- 
check it.

Your data (in its ciphered form) lives for maximum 24h in a GCS bucket. 

[![pipeline status](https://gitlab.com/lfaoro/lsh/badges/master/pipeline.svg)](https://gitlab
.com/lfaoro/lsh/commits/master)
[![coverage report](https://gitlab.com/lfaoro/lsh/badges/master/coverage.svg)](https://gitlab
.com/lfaoro/lsh/commits/master)

## Quick start

```bash
# macOS
brew install lfaoro/tap/lsh

# linux (WIP)
curl lsh.io/i | sh

# developers
go get -u github.com/lfaoro/lsh/cmd/lsh
cd $GOPATH/src/github.com/lfaoro/lsh/cmd/lsh
make install
lsh -h
make test
```

## Quick usage

```bash
lsh up genesis.txt -to hello@lsh.io
```

Subject: You've got data!
Body: Download your data from https://lsh.io/lsYuhN8h_genesis.txt

## Usage

```bash
$ lsh upload genesis.doc # short: lsh u
Success: https://lsh.io/lsYuhN8h_genesis.txt

$ lsh download 2E3fde2a-genesis.doc # short: lsh d
Downloaded genesis.doc

$ lsh encrypt genesis.doc # short: lsh en
Encrypted genesis.doc

$ lsh decrypt genesis.doc # short: lsh de
Decrypted genesis.doc

$ lsh u -p genesis.doc
Secret-key: *********
Success: https://lsh.io/lsYuhN8h_genesis.txt

$ lsh lock 
Password: **********
Locked .config/lsh/key

$ lsh unlock 
Password: **********
Unlocked .config/lsh/key

# WIP commands

$ lsh genesis.doc
🧮 Unable to decrypt using your local key
🔑 Decryption-key: ***********
🔓 Decrypted genesis.doc
```

## Leading encryption standard

Authenticated Encryption with Additional Authenticated Data (AEAD) couples confidentiality and integrity. Using the 
most popular AEAD today: AES-GCM. 

The "AES-GCM" algorithm identifier is used to perform authenticated encryption and decryption using AES in 
Galois/Counter Mode mode, as described in [NIST SP 800-38D](https://csrc.nist.gov/publications/detail/sp/800-38d/final)

ref paper: https://eprint.iacr.org/2017/168.pdf

## Contributing

> Any help, feedback and suggestions are very welcome and greatly appreciated.
> Start by opening an [issue](https://github.com/lfaoro/pkg/issues/new).

## Motivation

It's hard to find a service one can completely trust -- everybody claims they're encrypting your data, although how can you be sure? 

I believe the only way trust what happens to your data is to see exactly the steps that lead to its manipulation, 
encryption & storage.

lsh is F/OSS -- anyone can check how data is being encrypted and handled, spot eventual issues and fix insecurities.

## Compliance (WIP)

Right now lsh stores the encryption keys in a `key` file, located in `$HOME/.config/lsh` with `0600` permission
. Ideally we'll have the keys stored in macOS keychain -- although I don't know if there's something comparable for 
Linux and Windows.

To comply with regulators you might need to generate encryption keys using a Hardware Security Module aka HSM. 

lsh comes with a HSM plugin for GCP and AWS. These providers offer HSM as a service. 

Configure the GCP/AWS environment variables in order to activate Cloud HSM; ref: https://.

> In progress: https://github.com/lfaoro/lsh/issues/1