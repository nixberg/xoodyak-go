package xoodyak

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"testing"
)

//go:embed test/hash.json
var hashKats []byte

func TestHashMode(t *testing.T) {
	var kats []struct {
		Msg string `json:"msg"`
		MD  string `json:"md"`
	}
	json.Unmarshal(hashKats, &kats)

	for i, kat := range kats {
		msg, _ := hex.DecodeString(kat.Msg)
		md, _ := hex.DecodeString(kat.MD)

		xoodyak := New()
		xoodyak.Absorb(msg)
		newMD := xoodyak.Squeeze(nil, len(md))

		if !bytes.Equal(md, newMD) {
			t.Errorf("kats/hash: %d", i)
		}

		xoodyak = New()
		xoodyak.Absorb(msg)
		newMD = xoodyak.Squeeze([]byte("prefix"), len(md))

		if !bytes.Equal(md, newMD[6:]) {
			t.Errorf("kats/hash: %d", i)
		}
	}
}
