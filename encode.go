package qrcode

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"github.com/pkg/errors"

	qrcode "github.com/yeqown/go-qrcode"
	"snapcore.com/qrcode/cose"
)

// Encode
// JSON -> CBOR+COSE -> ZLIB -> Base45 -> QRCode generation
func encodeToFile(signer *cose.Signer, input []byte, qrFile string) error {

	msg, err := signCOSE(signer, input)
	if err != nil {
		return errors.Wrap(err, "Cannot sign")
	}

	//fmt.Printf("cose len %d - %x\n", len(msg), msg)

	compressed, err := compressZLIB(msg)
	if err != nil {
		return errors.Wrap(err, "Cannot compress zlib")
	}
	//fmt.Printf("compressed len %d - %x\n", len(compressed), compressed)

	qrcodebin, err := Base45Encode(compressed)
	if err != nil {
		return errors.Wrap(err, "Cannot generate qrCode")
	}

	qrcodestr := string(qrcodebin)
	//fmt.Printf("qrcodebin len %d - %s\n", len(qrcodestr), qrcodestr)

	err = genQRCode(qrcodestr, qrFile)
	if err != nil {
		return errors.Wrap(err, "Cannot generate qrCode")
	}
	return nil
}

func genQRCode(qrcodestr string, destinationFile string) error {

	encOpts := qrcode.DefaultConfig()
	encOpts.EcLevel = qrcode.ErrorCorrectionQuart
	encOpts.EncMode = qrcode.EncModeAuto

	qrc, err := qrcode.NewWithConfig(qrcodestr, encOpts, qrcode.WithQRWidth(20))
	if err != nil {
		return errors.Wrap(err, "could not generate QRCode")
	}

	// save file
	if err = qrc.Save(destinationFile); err != nil {
		return errors.Wrap(err, "could not save image")
	}
	return nil
}

func signCOSE(signer *cose.Signer, input []byte) ([]byte, error) {

	// create a signature
	sig := cose.NewSignature()
	sig.Headers.Unprotected["kid"] = 1
	sig.Headers.Protected["alg"] = "ES384"

	// create a message
	//external := []byte("") // optional external data see https://tools.ietf.org/html/rfc8152#section-4.3

	msg := cose.NewSignMessage()
	msg.Payload = input
	msg.AddSignature(sig)

	err := msg.Sign(rand.Reader, nil, []cose.Signer{*signer})
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
