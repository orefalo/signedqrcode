package main

import (
	"bytes"
	"compress/zlib"
	"crypto"
	cose "example.com/main/go-cose"
	"fmt"
	"github.com/fxamacker/cbor/v2"

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

	output, err = verifyCOSE(msg, publickey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}


func verifyCOSE(input []byte, publickey crypto.PublicKey) (output string, err error) {

	verifier, err := cose.NewVerifierFromKey(cose.ES384, &publickey)
	var msg cose.SignMessage
	err = cbor.Unmarshal(input, &msg)

	external := []byte("")

	err = msg.Verify(external, []cose.Verifier{*verifier})
	if err == nil {
		fmt.Println("Message signature verified")
		return "", nil
	} else {
		fmt.Println(fmt.Sprintf("Error verifying the message %+v", err))
		return "", err
	}

}

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

