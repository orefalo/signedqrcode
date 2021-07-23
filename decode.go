package qrcode

import (
	"bytes"
	"compress/zlib"
	"crypto"
	"github.com/makiuchi-d/gozxing"
	zxingqrcode "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
	"image"
	"io"
	"os"
	"snapcore.com/qrcode/cose"
)

// Decode
// Scan QRCode -> unBase45 -> unLZIB -> unCBOR+COSE -> JSON
func decode(qrcodeFile string, publickey crypto.PublicKey, alg *cose.Algorithm) (output []byte, err error) {

	qrcodestr, err := ocrQRCode(qrcodeFile)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot OCR qrcode")

	}

	compressed, err := Base45Decode([]byte(qrcodestr))
	if err != nil {
		return nil, errors.Wrap(err, "Cannot decode base45")
	}
	//fmt.Printf("compressed len %d - %x\n", len(compressed), compressed)

	msg, err := decompressZLIB(compressed)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot decompress zlib")

	}
	//fmt.Printf("cose len %d - %x\n", len(msg), msg)

	decoded, err := verifyCOSE(msg, publickey, alg)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot validate signature")
	}
	return decoded, nil
}

func verifyCOSE(input []byte, publickey crypto.PublicKey, alg *cose.Algorithm) (output []byte, err error) {

	verifier, err := cose.NewVerifierFromKey(alg, publickey)

	msg := cose.NewSignMessage()
	err = msg.UnmarshalCBOR(input)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshall")
	}

	err = msg.Verify(nil, []cose.Verifier{*verifier})
	if err != nil {
		return nil, errors.Wrap(err, "Cannot validate signature")
	}
	return msg.Payload, nil
}

func decompressZLIB(input []byte) ([]byte, error) {

	b := bytes.NewReader(input)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	_, err = io.Copy(&output, r)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func ocrQRCode(qrcodeFile string) (msg string, err error) {
	// open and decode image file
	file, _ := os.Open(qrcodeFile)
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := zxingqrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	return result.String(), nil
}
