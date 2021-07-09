package main

import (
	"bytes"
	"compress/zlib"
	"crypto"
	"fmt"

	"github.com/makiuchi-d/gozxing"
	zxingqrcode "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
	"image"
	"io"
	"os"
)


// Decode
// Scan QRCode -> unBase45 -> unLZ4 -> unCBOR -> JSON
func decode(qrcodeFile string, publickey crypto.PublicKey) {

	fmt.Println("==================================================")

	qrcodestr, err := ocrQRCode2(qrcodeFile)
	if err != nil {
		errors.Errorf("%s", err)
	}

	compressed, err := Base45Decode([]byte(qrcodestr))
	if err != nil {
		fmt.Printf("could not decode base45: %s", err)
	}
	fmt.Printf("compressed len %d - %x\n", len(compressed), compressed)

	msg := decompressZLIB(compressed)
	fmt.Printf("cose len %d - %x\n", len(msg), msg)

	//err = verifyCOSE(msg)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
}

func ocrQRCode2(qrcodeFile string) (msg string, err error) {
	// open and decode image file
	file, _ := os.Open(qrcodeFile)
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := zxingqrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	output := result.String()
	fmt.Println("read from qr %s" + output)
	return output, nil
}

//func verifyCOSE([]byte) error {
//
//
//	var msg SignMessage
//	err := cbor.Unmarshal(testCase.bytes, &msg)
//
//
//
//
//	// create a signer with a new private key
//	signer, err := cose.NewSigner(cose.ES384, nil)
//	if err != nil {
//		panic(fmt.Sprintf(fmt.Sprintf("Error creating signer -> %s", err)))
//	}
//
//	// create a signature
//	sig := cose.NewSignature()
//	sig.Headers.Unprotected["kid"] = 1
//	sig.Headers.Protected["alg"] = "ES384"
//
//	// create a message
//	external := []byte("") // optional external data see https://tools.ietf.org/html/rfc8152#section-4.3
//
//	msg := cose.NewSignMessage()
//	msg.Payload = []byte("payload to sign")
//	msg.AddSignature(sig)
//
//	err = msg.Sign(rand.Reader, external, []cose.Signer{*signer})
//	if err == nil {
//		fmt.Println(fmt.Sprintf("Message signature (ES256): %x", msg.Signatures[0].SignatureBytes))
//	} else {
//		panic(fmt.Sprintf("Error signing the message %+v", err))
//	}
//
//	// derive a verifier using the signer's public key and COSE algorithm
//	verifier := signer.Verifier()
//	verifier := signer.Verifier()
//	cose.NewSignerFromKey()
//
//	// Verify
//	err = msg.Verify(external, []cose.Verifier{*verifier})
//	if err == nil {
//		fmt.Println("Message signature verified")
//	} else {
//		fmt.Println(fmt.Sprintf("Error verifying the message %+v", err))
//	}
//	return nil
//}

func decompressZLIB(input []byte) []byte {

	b := bytes.NewReader(input)
	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	var output bytes.Buffer
	io.Copy(&output, r)
	r.Close()
	return output.Bytes()
}

//func svgToPng(inputSVG string, outputPNG string) {
//
//	icon, _ := oksvg.ReadIconStream(strings.NewReader(inputSVG))
//
//	//w := int(icon.ViewBox.W)
//	//h := int(icon.ViewBox.H)
//
//	w, h := 512, 512
//
//	icon.SetTarget(0, 0, float64(w), float64(h))
//	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
//	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
//
//	out, err := os.Create(outputPNG)
//	if err != nil {
//		panic(err)
//	}
//	defer out.Close()
//
//	err = png.Encode(out, rgba)
//	if err != nil {
//		panic(err)
//	}
//}

