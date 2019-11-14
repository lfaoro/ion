# 🧬 ion - AES256/GCM encrypt and share large data objects.

> End to end encrypted, if you want...

Encryption is done on your computer, your data does not hit the cloud unencrypted.

No logs except errors are being collected from [cmd/server](https://github.com/lfaoro/ion/tree/master/cmd/server) -- check it.

Your data (in its ciphered form) lives for maximum 24h in a GCS bucket. 

[![pipeline status](https://gitlab.com/lfaoro/ion/badges/master/pipeline.svg)](https://gitlab.com/lfaoro/ion/commits/master)
[![coverage report](https://gitlab.com/lfaoro/ion/badges/master/coverage.svg)](https://gitlab.com/lfaoro/ion/commits/master)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Flfaoro%2Fion.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Flfaoro%2Fion?ref=badge_shield)

## Quick start

```bash
# macOS
brew install lfaoro/tap/ion

# linux (WIP)
curl apionic.com/ion.sh | sh

# developers
go get -u github.com/lfaoro/ion
make install
ion -h

make test
```

## Usage

```bash
$ ion upload genesis.txt
13.09 MiB / 1.14 GiB [>-----------------------------]   1.12% 1.72 MiB/s 11m11s
Download from: https:/s.apionic.com/nERuG_genesis.txt

$ ion download nERuG_genesis.txt
290.19 MiB / 1.14 GiB [=======>---------------------]  24.83% 5.52 MiB/s 02m38s
Downloaded: genesis.txt

$ ion encrypt --key genesis.txt
🔑 Encryption-key: 238dFomyjB3wEejjoSUef97Y/k1gMib6XvVS56i4Apg=
🔒 Encrypted /tmp/genesis.txt

$ ion decrypt --key genesis.txt
🔑 Encryption-key: 238dFomyjB3wEejjoSUef97Y/k1gMib6XvVS56i4Apg=
🔓 Decrypted /tmp/genesis.txt
```

# WIP beta commands

```
$ ion lock 
Password: **********
Locked .config/ion/key

$ ion unlock 
Password: **********
Unlocked .config/ion/key

$ ion genesis.txt
🧮 Unable to decrypt using your local key
🔑 Decryption-key: ***********
🔓 Decrypted genesis.txt

$ ion up genesis.txt -to hello@lsh.io
```

## Sample email
Subject: You've got data!

Body: Download your data from https://s.apionic.com/lsYuh_genesis.txt

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

ion is F/OSS -- anyone can check how data is being encrypted and handled, spot eventual issues and fix insecurities.

## Compliance (WIP)

Right now lsh stores the encryption keys in a `key` file, located in `$HOME/.config/ion` with `0600` permission
. Ideally we'll have the keys stored in the macOS keychain -- although I don't know if there's something comparable for 
Linux and Windows.

To comply with regulators you might need to generate encryption keys using a Hardware Security Module aka HSM. 

ion comes with a HSM plugin for GCP and AWS. These providers offer HSM as a service. 

Configure the GCP/AWS environment variables in order to activate Cloud HSM; ref: https://.

> In progress: https://github.com/lfaoro/ion/issues/1

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Flfaoro%2Fion.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Flfaoro%2Fion?ref=badge_large)