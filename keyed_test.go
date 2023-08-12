package xoodyak

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/nixberg/blobby-go"
)

//go:embed testdata/aead.blb
var encodedAEADBlobs []byte

func TestKeyedMode(t *testing.T) {
	blobs := blobby.MustDecode(encodedAEADBlobs)

	for i := 0; i < len(blobs); i += 6 {
		key := blobs[i+0]
		nonce := blobs[i+1]
		ad := blobs[i+2]
		pt := blobs[i+3]
		ct := blobs[i+4]
		tag := blobs[i+5]

		encryptor := NewKeyed(key, nil, nil)
		encryptor.Absorb(nonce)
		encryptor.Absorb(ad)

		decryptor := *encryptor

		if !bytes.Equal(encryptor.Encrypt(pt, nil), ct) {
			t.Fail()
		}
		if !bytes.Equal(encryptor.Squeeze(nil, len(tag)), tag) {
			t.Fail()
		}

		if !bytes.Equal(decryptor.Decrypt(ct, nil), pt) {
			t.Fail()
		}
		if !bytes.Equal(decryptor.Squeeze(nil, len(tag)), tag) {
			t.Fail()
		}
	}
}
