package xoodoo

import "testing"

func TestXoodoo(t *testing.T) {
	var xoodoo Xoodoo

	for i := 0; i < 384; i++ {
		xoodoo.Permute()
	}

	expected := [12]uint32{
		0xfe04fab0, 0x42d5d8ce, 0x29c62ee7, 0x2a7ae5cf,
		0xea36eba3, 0x14649e0a, 0xfe12521b, 0xfe2eff69,
		0xf1826ca5, 0xfc4c41e0, 0x1597394f, 0xeb092faf,
	}

	if xoodoo.s != expected {
		t.Fail()
	}
}

func TestGetSet(t *testing.T) {
	var xoodoo Xoodoo

	for i := 0; i < 48; i++ {
		xoodoo.XOR(i, byte(i))
		if xoodoo.Get(i) != byte(i) {
			t.Fail()
		}
	}
}
