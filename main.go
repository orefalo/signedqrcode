package qrcode

import (
	"encoding/json"
	"fmt"
	cose "snapcore.com/qrcode/cose"
)

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func generateSigner() *cose.Signer {
	// create a signer with a new private key
	signer, err := cose.NewSigner(cose.ES384, nil)
	if err != nil {
		panic(fmt.Sprintf(fmt.Sprintf("Error creating signer %s", err)))
	}
	return signer
}


//func ParseRsaPrivateKeyFromPemStr(privPEM string) (*crypto.PrivateKey, error) {
//	block, _ := pem.Decode([]byte(privPEM))
//	if block == nil {
//		return nil, errors.New("failed to parse PEM block containing the key")
//	}
//
//	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return nil, err
//	}
//
//	return priv, nil
//}

func main() {

	signer := generateSigner()

	err:=encodeToFile(signer, []byte(jsonEscape(`{
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
	if err != nil {
		fmt.Println(err)
	} else {

		verifier := signer.Verifier()
		decoded, err := decode("./qrcode.png", verifier.PublicKey, verifier.Alg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(decoded))
		}
	}
}
