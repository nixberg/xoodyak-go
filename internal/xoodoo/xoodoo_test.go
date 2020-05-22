package xoodoo

import "testing"

func TestXoodoo(t *testing.T) {
	var xoodoo Xoodoo

	for i := 0; i < 384; i++ {
		xoodoo.Permute()
	}

	expected := [48]byte{
		0xb0, 0xfa, 0x04, 0xfe, 0xce, 0xd8, 0xd5, 0x42,
		0xe7, 0x2e, 0xc6, 0x29, 0xcf, 0xe5, 0x7a, 0x2a,
		0xa3, 0xeb, 0x36, 0xea, 0x0a, 0x9e, 0x64, 0x14,
		0x1b, 0x52, 0x12, 0xfe, 0x69, 0xff, 0x2e, 0xfe,
		0xa5, 0x6c, 0x82, 0xf1, 0xe0, 0x41, 0x4c, 0xfc,
		0x4f, 0x39, 0x97, 0x15, 0xaf, 0x2f, 0x09, 0xeb,
	}

	if xoodoo.Bytes != expected {
		t.Fail()
	}
}

var result [48]byte

func BenchmarkXoodoo(b *testing.B) {
	var xoodoo Xoodoo
	for i := 0; i < 128*1024; i++ {
		xoodoo.Permute()
	}
	result = xoodoo.Bytes
}