package rockgo

import (
	"encoding/base64"
)

func DoXOR(source []byte, key []byte) []byte {
	keyLen := len(key)
	for index, item := range source {
		source[index] = item ^ key[index%keyLen]
	}
	return source
}

func EncodeXOR(source []byte, key string) []byte {
	keyData := []byte(key)
	keyBytes := make([]byte, base64.RawURLEncoding.EncodedLen(len(keyData)))
	base64.RawURLEncoding.Encode(keyBytes, keyData)
	return DoXOR(source, keyBytes)
}

func EncodeXORToString(source []byte, key string) string {
	resBytes := EncodeXOR(source, key)
	return base64.RawURLEncoding.EncodeToString(resBytes)
}

func DecodeXOR(data string, key string) ([]byte, error) {
	keyData := []byte(key)
	keyBytes := make([]byte, base64.RawURLEncoding.EncodedLen(len(keyData)))
	base64.RawURLEncoding.Encode(keyBytes, keyData)

	sourceData, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	resBytes := DoXOR(sourceData, keyBytes)
	return resBytes, nil
}
