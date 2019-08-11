package xoodyak

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestHashMode(t *testing.T) {
	type KAT struct {
		Msg string `json:"msg"`
		MD  string `json:"md"`
	}

	in, err := ioutil.ReadFile("kats/hash.json")
	if err != nil {
		panic(err)
	}

	kats := make([]KAT, 0)
	json.Unmarshal(in, &kats)

	for i, kat := range kats {
		msg, _ := hex.DecodeString(kat.Msg)
		md, _ := hex.DecodeString(kat.MD)

		xoodyak := New()
		xoodyak.Absorb(msg)
		newMD := xoodyak.Squeeze([]byte("header"), len(md))

		if !bytes.Equal(md, newMD[6:]) {
			t.Errorf("kats/hash: %d", i)
		}
	}
}

func TestKeyedMode(t *testing.T) {
	type KAT struct {
		Key   string `json:"key"`
		Nonce string `json:"nonce"`
		PT    string `json:"pt"`
		AD    string `json:"ad"`
		CT    string `json:"ct"`
	}

	in, err := ioutil.ReadFile("kats/aead.json")
	if err != nil {
		panic(err)
	}

	kats := make([]KAT, 0)
	json.Unmarshal(in, &kats)

	for i, kat := range kats {
		key, _ := hex.DecodeString(kat.Key)
		nonce, _ := hex.DecodeString(kat.Nonce)
		pt, _ := hex.DecodeString(kat.PT)
		ad, _ := hex.DecodeString(kat.AD)
		ct, _ := hex.DecodeString(kat.CT)
		tag := ct[len(pt):]

		xoodyak := Keyed(key, nil, nil)
		xoodyak.Absorb(nonce)
		xoodyak.Absorb(ad)
		newCT := xoodyak.Encrypt(pt, []byte("header"))
		newCT = xoodyak.Squeeze(newCT, len(tag))

		if !bytes.Equal(ct, newCT[6:]) {
			t.Errorf("kats/aead: %d", i)
		}

		xoodyak = Keyed(key, nil, nil)
		xoodyak.Absorb(nonce)
		xoodyak.Absorb(ad)
		newPT := xoodyak.Decrypt(ct[:len(pt)], []byte("header"))
		newTag := xoodyak.Squeeze([]byte("header"), len(tag))

		if !bytes.Equal(pt, newPT[6:]) {
			t.Errorf("kats/aead: %d", i)
		}

		if !bytes.Equal(tag, newTag[6:]) {
			t.Errorf("kats/aead: %d", i)
		}
	}
}
