package xoodyak

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/nixberg/blobby-go"
)

//go:embed testdata/hash.blb
var encodedHashBlobs []byte

func TestHashMode(t *testing.T) {
	blobs := blobby.MustDecode(encodedHashBlobs)

	for i := 0; i < len(blobs); i += 2 {
		msg := blobs[i+0]
		md := blobs[i+1]

		xoodyak := New()
		xoodyak.Absorb(msg)

		if !bytes.Equal(xoodyak.Squeeze(nil, len(md)), md) {
			t.Fail()
		}
	}
}
