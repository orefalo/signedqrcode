package main

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"fmt"

	qrcode "github.com/yeqown/go-qrcode"
	cose "go.mozilla.org/cose"
)

// Encode
// JSON -> CBOR -> COSE -> LZ4 -> Base45 -> QRCode generation
func encode(signer *cose.Signer, input []byte, file string) {

	fmt.Println("==================================================")

	msg, err := signCOSE(signer, input)
	if err != nil {
		panic("cose error")
	}
	fmt.Printf("cose len %d - %x\n", len(msg), msg)

	compressed := compressZLIB(msg)
	fmt.Printf("compressed len %d - %x\n", len(compressed), compressed)

	qrcodebin, err := Base45Encode(compressed)
	if err != nil {
		fmt.Printf("could not generate QRCode: %s", err)
	}
	qrcodestr := string(qrcodebin)
	fmt.Printf("qrcodebin len %d - %s\n", len(qrcodestr), qrcodestr)

	genQRCode2(qrcodestr, file)
}

func genQRCode2(qrcodestr string, destinationFile string) {
	qrc, err := qrcode.New(qrcodestr)
	if err != nil {
		fmt.Printf("could not generate QRCode: %s", err)
	}

	// save file
	if err = qrc.Save(destinationFile); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}

func signCOSE(keypair *cose.Signer, input []byte) ([]byte, error) {

	// create a signature
	sig := cose.NewSignature()
	sig.Headers.Unprotected["kid"] = 1
	sig.Headers.Protected["alg"] = "ES384"

	// create a message
	external := []byte("") // optional external data see https://tools.ietf.org/html/rfc8152#section-4.3

	msg := cose.NewSignMessage()
	msg.Payload = input
	msg.AddSignature(sig)

	err := msg.Sign(rand.Reader, external, []cose.Signer{*keypair})
	if err == nil {
		return msg.MarshalCBOR()
	} else {
		panic(fmt.Sprintf("Error signing the message -> %+v", err))
	}
	return nil, nil
}

func compressZLIB(input []byte) []byte {
	var output bytes.Buffer
	w := zlib.NewWriter(&output)
	_, err := w.Write(input)
	if err != nil {
		return nil
	}
	w.Close()
	return output.Bytes()
}
