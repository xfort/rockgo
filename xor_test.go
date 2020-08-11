package rockgo

import (
	"bytes"
	"log"
	"testing"
)

func TestXOR(t *testing.T) {
	key := "com.quantumfiture.MuscleMen"
	//keyBytes := make([]byte, base64.URLEncoding.EncodedLen(len(key)))
	//base64.URLEncoding.Encode(keyBytes, key)
	dataStr := "abc123汉语"
	data := []byte(dataStr)
	log.Println(data)

	resStr := EncodeXORToString(data, key)

	resBytes, err := DecodeXOR(resStr, key)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(resBytes, bytes.NewBuffer(resBytes).String())
}
