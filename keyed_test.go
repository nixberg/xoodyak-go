package xoodyak

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"testing"
)

//go:embed test/aead.json
var aeadKats []byte

func TestKeyedMode(t *testing.T) {
	var kats []struct {
		Key   string `json:"key"`
		Nonce string `json:"nonce"`
		PT    string `json:"pt"`
		AD    string `json:"ad"`
		CT    string `json:"ct"`
	}
	json.Unmarshal(aeadKats, &kats)

	for i, kat := range kats {
		key, _ := hex.DecodeString(kat.Key)
		nonce, _ := hex.DecodeString(kat.Nonce)
		pt, _ := hex.DecodeString(kat.PT)
		ad, _ := hex.DecodeString(kat.AD)
		ct, _ := hex.DecodeString(kat.CT)
		tag := ct[len(pt):]

		encryptor := NewKeyed(key, nil, nil)
		encryptor.Absorb(nonce)
		encryptor.Absorb(ad)
		decryptor := *encryptor

		newCT := encryptor.Encrypt(pt, []byte("prefix"))
		newCT = encryptor.Squeeze(newCT, len(tag))

		if !bytes.Equal(ct, newCT[6:]) {
			t.Errorf("kats/aead: %d", i)
		}

		newPT := decryptor.Decrypt(ct[:len(pt)], []byte("prefix"))
		newTag := decryptor.Squeeze([]byte("prefix"), len(tag))

		if !bytes.Equal(pt, newPT[6:]) {
			t.Errorf("kats/aead: %d", i)
		}
		if !bytes.Equal(tag, newTag[6:]) {
			t.Errorf("kats/aead: %d", i)
		}
	}
}
