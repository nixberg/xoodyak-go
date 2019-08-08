package xoodoo

import "unsafe"

var isBigEndian = false

func init() {
	var i uint16 = 0x1234
	if 0x12 == *(*byte)(unsafe.Pointer(&i)) {
		isBigEndian = true
	}
}

type Xoodoo struct {
	s [12]uint32
}

func rotate(v uint32, n uint32) uint32 {
	return (v >> n) | (v << (32 - n))
}

func (x *Xoodoo) Permute() {
	roundConstants := [12]uint32{
		0x058, 0x038, 0x3c0, 0x0d0,
		0x120, 0x014, 0x060, 0x02c,
		0x380, 0x0f0, 0x1a0, 0x012,
	}

	for _, roundConstant := range roundConstants {
		var e [4]uint32

		for i := 0; i < 4; i++ {
			e[i] = rotate(x.s[i]^x.s[i+4]^x.s[i+8], 18)
			e[i] ^= rotate(e[i], 9)
		}

		for i := 0; i < 12; i++ {
			x.s[i] ^= e[(i-1)&3]
		}

		x.s[7], x.s[4] = x.s[4], x.s[7]
		x.s[7], x.s[5] = x.s[5], x.s[7]
		x.s[7], x.s[6] = x.s[6], x.s[7]
		x.s[0] ^= roundConstant

		for i := 0; i < 4; i++ {
			a := x.s[i]
			b := x.s[i+4]
			c := rotate(x.s[i+8], 21)

			x.s[i+8] = rotate((b&^a)^c, 24)
			x.s[i+4] = rotate((a&^c)^b, 31)
			x.s[i] ^= c & ^b
		}

		x.s[8], x.s[10] = x.s[10], x.s[8]
		x.s[9], x.s[11] = x.s[11], x.s[9]
	}
}

func (x *Xoodoo) Get(index int) byte {
	if isBigEndian {
		index += 3 - 2*(index%4)
	}
	bytes := (*[48]byte)(unsafe.Pointer(&x.s))
	return bytes[index]
}

func (x *Xoodoo) XOR(index int, b byte) {
	if isBigEndian {
		index += 3 - 2*(index%4)
	}
	bytes := (*[48]byte)(unsafe.Pointer(&x.s))
	bytes[index] ^= b
}
