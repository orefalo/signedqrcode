package main

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/grkuntzmd/qrcodegen"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	qrcodeocr "github.com/tuotoo/qrcode"
	qrcode "github.com/yeqown/go-qrcode"
	cose "go.mozilla.org/cose"
	"image"
	"image/png"
	"os"
	"strings"
)

// Decode
// Scan QRCode -> unBase45 -> unLZ4 -> unCBOR -> JSON

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

// Encode
//JSON -> CBOR -> COSE -> LZ4 -> Base45 -> QRCode generation
// JSON -> COSE -> CBOR -> LZ4 -> Base45

func main() {
	encode([]byte(jsonEscape(`{
    "ver": "1.2.1",
    "nam": {
        "fn": "Musterfrau-G\u00f6\u00dfinger",
        "gn": "Gabriele",
        "fnt": "MUSTERFRAU<GOESSINGER",
        "gnt": "GABRIELE"
    },
    "dob": "1998-02-26",
    "v": [
        {
            "tg": "840539006",
            "vp": "1119349007",
            "mp": "EU\/1\/20\/1528",
            "ma": "ORG-100030215",
            "dn": 1,
            "sd": 2,
            "dt": "2021-02-18",
            "co": "AT",
            "is": "Ministry of Health, Austria",
            "ci": "URN:UVCI:01:AT:10807843F94AEE0EE5093FBC254BD813#B"
        }
    ]
}`)), "./qrcode.png")

	decode("./qrcode.png")
}

func decode(qrcode string) {

	fi, err := os.Open(qrcode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer fi.Close()
	qrmatrix, err := qrcodeocr.Decode(fi)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(qrmatrix.Content)

}

func encode(input []byte, file string) {
	msg, err := signCOSE(input)
	if err != nil {
		panic("cose error")
	}
	fmt.Printf("msg len %d - %x\n", len(msg), msg)
	compressed := compressZLIB(msg)

	fmt.Printf("compressed len %d - %x\n", len(compressed), compressed)

	qrcodebin, err := Base45Encode(compressed)
	if err != nil {
		fmt.Printf("could not generate QRCode: %s", err)
	}

	qrcodestr := string(qrcodebin)
	fmt.Printf("qrcodebin len %d - %s\n", len(qrcodestr), qrcodestr)

	genQRCode1(qrcodestr, file)
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

func genQRCode1(qrcodestr string, destinationFile string) {

	segs := []*qrcodegen.QRSegment{
		qrcodegen.MakeAlphanumeric(qrcodestr),
		//qrcodegen.MakeNumeric("007020004930000600600300000000000050200010008006900400003700900020050001000008000"),
	}
	qrCode, err := qrcodegen.EncodeSegments(segs, qrcodegen.Quartile, qrcodegen.WithAutoMask())
	if err != nil {
		// Handle this.
	}
	svg, _ := qrCode.ToSVGString(4, true)

		// save svg
	//f, err := os.Create("./qr.svg")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//l, err := f.WriteString(svg)
	//if err != nil {
	//	fmt.Println(err)
	//	f.Close()
	//	return
	//}
	//fmt.Println(l, "bytes written successfully")
	//err = f.Close()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	svgToPng(svg, destinationFile)

}

func signCOSE(input []byte) ([]byte, error) {
	// create a signer with a new private key
	signer, err := cose.NewSigner(cose.ES384, nil)
	if err != nil {
		panic(fmt.Sprintf(fmt.Sprintf("Error creating signer %s", err)))
	}

	// create a signature
	sig := cose.NewSignature()
	//sig.Headers.Unprotected["kid"] = 1
	sig.Headers.Protected["alg"] = "ES384"

	// create a message
	external := []byte("") // optional external data see https://tools.ietf.org/html/rfc8152#section-4.3

	msg := cose.NewSignMessage()
	msg.Payload = input
	msg.AddSignature(sig)

	err = msg.Sign(rand.Reader, external, []cose.Signer{*signer})
	if err == nil {
		return msg.MarshalCBOR()
		//fmt.Println(fmt.Sprintf("Message signature (ES256): %x", msg.Signatures[0].SignatureBytes))
	} else {
		panic(fmt.Sprintf("Error signing the message -> %+v", err))
	}
	return nil, nil
}

func verifyCOSE() {
	// create a signer with a new private key
	signer, err := cose.NewSigner(cose.ES384, nil)
	if err != nil {
		panic(fmt.Sprintf(fmt.Sprintf("Error creating signer -> %s", err)))
	}

	// create a signature
	sig := cose.NewSignature()
	sig.Headers.Unprotected["kid"] = 1
	sig.Headers.Protected["alg"] = "ES384"

	// create a message
	external := []byte("") // optional external data see https://tools.ietf.org/html/rfc8152#section-4.3

	msg := cose.NewSignMessage()
	msg.Payload = []byte("payload to sign")
	msg.AddSignature(sig)

	err = msg.Sign(rand.Reader, external, []cose.Signer{*signer})
	if err == nil {
		fmt.Println(fmt.Sprintf("Message signature (ES256): %x", msg.Signatures[0].SignatureBytes))
	} else {
		panic(fmt.Sprintf("Error signing the message %+v", err))
	}

	// derive a verifier using the signer's public key and COSE algorithm
	verifier := signer.Verifier()

	// Verify
	err = msg.Verify(external, []cose.Verifier{*verifier})
	if err == nil {
		fmt.Println("Message signature verified")
	} else {
		fmt.Println(fmt.Sprintf("Error verifying the message %+v", err))
	}
}

func compressZLIB(input []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(input)
	w.Close()
	return b.Bytes()
}

func svgToPng(inputSVG string, outputPNG string) {

	icon, _ := oksvg.ReadIconStream(strings.NewReader(inputSVG))

	//w := int(icon.ViewBox.W)
	//h := int(icon.ViewBox.H)

	w, h := 512, 512

	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	out, err := os.Create(outputPNG)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		panic(err)
	}
}

//// BytesToUint16 converts a big endian array of bytes to an array of unit16s
//func BytesToUint16(bytes []byte) []uint16 {
//	values := make([]uint16, len(bytes)/2)
//
//	for i := range values {
//		values[i] = binary.BigEndian.Uint16(bytes[i*2 : (i+1)*2])
//	}
//	return values
//}
//
//// Uint16ToBytes converts an array of uint16s to a big endian array of bytes
//func Uint16ToBytes(values []uint16) []byte {
//	bytes := make([]byte, len(values)*2)
//
//	for i, value := range values {
//		binary.BigEndian.PutUint16(bytes[i*2:(i+1)*2], value)
//	}
//	return bytes
//}

//func generateKeyPay() (*g1pubs.PublicKey, *g1pubs.SecretKey) {
//	var src cryptoSource
//	rnd := rand.New(src)
//	r := NewXORShift(rnd.Uint64())
//	privateKey, _ := g1pubs.RandKey(r)
//	publicKey := g1pubs.PrivToPub(privateKey)
//	return publicKey, privateKey
//}