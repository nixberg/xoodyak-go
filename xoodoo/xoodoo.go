package xoodoo

import "math/bits"

type Xoodoo struct {
	Bytes [48]byte
}

func (x *Xoodoo) Permute() {
	var s [12]uint32
	for i, j := 0, 0; i < 12; i++ {
		s[i] |= (uint32(x.Bytes[j]) << 0)
		j++
		s[i] |= (uint32(x.Bytes[j]) << 8)
		j++
		s[i] |= (uint32(x.Bytes[j]) << 16)
		j++
		s[i] |= (uint32(x.Bytes[j]) << 24)
		j++
	}

	roundConstants := [12]uint32{
		0x058, 0x038, 0x3c0, 0x0d0,
		0x120, 0x014, 0x060, 0x02c,
		0x380, 0x0f0, 0x1a0, 0x012,
	}

	for _, roundConstant := range roundConstants {
		var e [4]uint32

		for i := 0; i < 4; i++ {
			e[i] = bits.RotateLeft32(s[i]^s[i+4]^s[i+8], -18)
			e[i] ^= bits.RotateLeft32(e[i], -9)
		}

		for i := 0; i < 12; i++ {
			s[i] ^= e[(i-1)&3]
		}

		s[7], s[4] = s[4], s[7]
		s[7], s[5] = s[5], s[7]
		s[7], s[6] = s[6], s[7]
		s[0] ^= roundConstant

		for i := 0; i < 4; i++ {
			a := s[i]
			b := s[i+4]
			c := bits.RotateLeft32(s[i+8], -21)

			s[i+8] = bits.RotateLeft32((b&^a)^c, -24)
			s[i+4] = bits.RotateLeft32((a&^c)^b, -31)
			s[i] ^= c & ^b
		}

		s[8], s[10] = s[10], s[8]
		s[9], s[11] = s[11], s[9]
	}

	for i, j := 0, 0; i < 12; i++ {
		x.Bytes[j] = byte(s[i] >> 0)
		j++
		x.Bytes[j] = byte(s[i] >> 8)
		j++
		x.Bytes[j] = byte(s[i] >> 16)
		j++
		x.Bytes[j] = byte(s[i] >> 24)
		j++
	}
}
