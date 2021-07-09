package main

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"fmt"
	"github.com/pkg/errors"

	cose "example.com/main/go-cose"
	qrcode "github.com/yeqown/go-qrcode"
)

// Encode
// JSON -> CBOR -> COSE -> LZ4 -> Base45 -> QRCode generation
func encodeToFile(signer *cose.Signer, input []byte, qrFile string) error {

	msg, err := signCOSE(signer, input)
	if err != nil {
		return errors.Wrap(err, "Cannot sign")
	}

	fmt.Printf("cose len %d - %x\n", len(msg), msg)

	compressed, err := compressZLIB(msg)
	if err != nil {
		return errors.Wrap(err, "Cannot compress zlib")
	}
	fmt.Printf("compressed len %d - %x\n", len(compressed), compressed)

	qrcodebin, err := Base45Encode(compressed)
	if err != nil {
		return errors.Wrap(err, "Cannot generate qrCode")
	}

	qrcodestr := string(qrcodebin)
	fmt.Printf("qrcodebin len %d - %s\n", len(qrcodestr), qrcodestr)

	err = genQRCode2(qrcodestr, qrFile)
	if err != nil {
		return errors.Wrap(err, "Cannot generate qrCode")
	}
	return nil
}

func genQRCode2(qrcodestr string, destinationFile string) error {
	qrc, err := qrcode.New(qrcodestr)
	if err != nil {
		return errors.Wrap(err, "could not generate QRCode")
	}

	// save file
	if err = qrc.Save(destinationFile); err != nil {
		return errors.Wrap(err, "could not save image")
	}
	return nil
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
	if err != nil {
		return nil, errors.Wrap(err, "Error signing the message")
	}
	return msg.MarshalCBOR()

}

func compressZLIB(input []byte) ([]byte, error) {
	var output bytes.Buffer
	w := zlib.NewWriter(&output)
	_, err := w.Write(input)
	if err != nil {
		return nil, err
	}
	w.Close()
	return output.Bytes(), nil
}
