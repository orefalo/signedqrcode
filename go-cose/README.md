# Why forking Mozilla go-cose?

This fork of [Mozilla go-cose](https://github.com/mozilla-services/go-cose)

adds
- [Sign1](https://tools.ietf.org/html/rfc8152#section-4.2) capabilities to the parent package.
- Ability to create Verifier from public key

## Usage

### Install

```console
go get -u github.com/orefalo/go-cose
```

### Signing a message

See [example/sign.go](example/sign.go) and run it with:

```console
$ go run example/sign.go
Bit lengths of integers r and s (256 and 256) do not match the key length 255
Message signature (ES256): 043685f99421f9e80c7c3c50d0fc8266161d3d614aaa3b63d2cdf581713fca62bb5d2e34d2352dbe41424b31d0b4a11d6b2d4764c18e2af04f4520fbe494d51c
```

### Verifying a message

See [example/verify.go](example/verify.go) and run it with:

```console
$ go run example/verify.go
Bit lengths of integers r and s (256 and 254) do not match the key length 254
Message signature (ES256): 9411dc5200c1cb67ccd76424ade09ce89c4a8d8d2b66f2bbf70edf63beb2dc3cbde83250773e659b635d3715442a1efaa6b0c030ee8a2523c3e37a22ddb055fa
Message signature verified
```

### Signing a message with one signer (Sign1)

See [example/sign1.go](example/sign1.go) and run it with:

```console
$ go run example/sign1.go
2020/11/24 14:13:48 COSE Sign1 signature bytes: 424285def70f75909cc763640d7ace555a6dafa576e0ed1b50964b80efbd92746e268ee6e5fa58258ae873bc1a510b1cc964b5c3905100ea0e625253a10150df
2020/11/24 14:13:48 COSE Sign1 message: d28443a10126a10401586da70a49948f8860d13a463e8e0b530198f50a4ff6c05861c8860d13a638ea4fe2fa0ff5100306c11a5afd322e0e0314a36f416e64726f69642041707020466f6fa10e017253656375726520456c656d656e7420456174d83dd24201236d4c696e757820416e64726f6964a10e015840424285def70f75909cc763640d7ace555a6dafa576e0ed1b50964b80efbd92746e268ee6e5fa58258ae873bc1a510b1cc964b5c3905100ea0e625253a10150df
2020/11/24 14:13:48 COSE Sign1 signature verified
```

## Development

Running tests:

1. Install [rust and cargo](https://www.rustup.rs/)

1. On OSX: `brew install nss` [nss](https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS) then in `sign_verify_cose_rust_cli_test.go` add `NSS_LIB_DIR` to `cmd` or `-L /usr/local/opt/nss/lib` to RUSTFLAGS e.g. `cmd.Env = append(os.Environ(), "NSS_LIB_DIR=/usr/local/opt/nss/lib", "RUSTFLAGS=-A dead_code -A unused_imports")`

1. If you already have `dep` and `golint` commands installed, run `make install-godep install-golint`

1. Run `go test`
