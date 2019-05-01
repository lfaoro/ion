# 🔏 Encrypto - encryption helpers.

### AES / GCM

Standard library implementation of AES at 256bit w/ Galois Counter Mode (GCM) data authentication.

One of the most popular authenticated
encryption schemes today is AES-GCM due to its impressive speed.

> Paper: https://eprint.iacr.org/2015/102.pdf

### CloudKMS

Makes it easy to interact with GCP's CloudKMS service.

Assumes you have the "GOOGLE_APPLICATION_CREDENTIALS" environment
variable setup in your environment with access to the Cloud KMS service.

Authentication documentation: https://cloud.google.com/docs/authentication/getting-started
Go client library: https://cloud.google.com/kms/docs/reference/libraries#client-libraries-install-go

Remember to create a KeyRing and CryptoKey.
Documentation: https://cloud.google.com/kms/docs/creating-keys

CloudKMS pricing: https://cloud.google.com/kms/pricing

### Various helpers
- Random string generator
- Token generator
- HMAC512 signing
- SHA256 hashing