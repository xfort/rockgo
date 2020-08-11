package rockgo

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	hasher := md5.New()
	_, err := hasher.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(hex.EncodeToString(hasher.Sum(nil)))
}
