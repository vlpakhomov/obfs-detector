package detector

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"obfs-detector/pkg/null"
)

type obfs2 struct {
	magicValue uint32
	seedLength int
	keyLength  int
	paddingStr string
}

var Obfs2 obfs2 = obfs2{
	magicValue: 0x2BF5CA7E,
	seedLength: 16,
	keyLength:  16,
	paddingStr: "Initiator obfuscated data",
}

func (o *obfs2) Detect(data []byte) (null.Null[bool], null.Null[string], error) {
	seed := data[:o.seedLength]
	encrypted := data[o.seedLength:]
	mac := o.paddingStr + string(seed) + o.paddingStr

	h := sha256.New()
	if _, err := h.Write([]byte(mac)); err != nil {
		return null.NewExplicit(false, false), null.NewExplicit("", false), err
	}

	hmac := h.Sum(nil)
	aesCTR128Key := hmac[:o.keyLength]
	aesCTR128IV := hmac[o.keyLength:]

	cph, err := aes.NewCipher(aesCTR128Key)
	if err != nil {
		return null.NewExplicit(false, false), null.NewExplicit("", false), err
	}

	stream := cipher.NewCTR(cph, aesCTR128IV)

	decrypted := make([]byte, len(encrypted))
	stream.XORKeyStream(decrypted, encrypted[aes.BlockSize:])

	magicValue := binary.BigEndian.Uint32(decrypted)
	if magicValue == o.magicValue {
		return null.New(true), null.New("obfs2 detected"), nil
	}

	return null.New(false), null.New("obfs2 undetected"), nil
}
