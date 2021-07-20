<img src="/Users/orefalo/GitRepositories/qr-related/signedqrcode/qrcode.png" style="zoom:10%;" />



### Introduction

A quick library to generate and validate PKI signed QR Codes

**Encode**: JSON Document  -> CBOR+COSE -> ZLIB -> Base45 -> QRCode generation (PNG)  
**Decode**: PNG -> Scan QRCode -> unBase45 -> unLZIB -> unCBOR+COSE -> JSON Document

Inspired and compliant with the [HCERT Spec](https://github.com/ehn-dcc-development/hcert-spec) 

The project implements its own customized version of COSE along with other key QR code projects

### Testing

```shell
go test
```

Enjoy,
Olivier Refalo