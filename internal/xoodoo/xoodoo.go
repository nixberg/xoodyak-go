package xoodoo

import (
	"encoding/binary"
	"math/bits"
)

type Xoodoo struct {
	Bytes [48]byte
}

func (x *Xoodoo) Permute() {
	s := [12]uint32{
		binary.LittleEndian.Uint32(x.Bytes[0:4]),
		binary.LittleEndian.Uint32(x.Bytes[4:8]),
		binary.LittleEndian.Uint32(x.Bytes[8:12]),
		binary.LittleEndian.Uint32(x.Bytes[12:16]),
		binary.LittleEndian.Uint32(x.Bytes[16:20]),
		binary.LittleEndian.Uint32(x.Bytes[20:24]),
		binary.LittleEndian.Uint32(x.Bytes[24:28]),
		binary.LittleEndian.Uint32(x.Bytes[28:32]),
		binary.LittleEndian.Uint32(x.Bytes[32:36]),
		binary.LittleEndian.Uint32(x.Bytes[36:40]),
		binary.LittleEndian.Uint32(x.Bytes[40:44]),
		binary.LittleEndian.Uint32(x.Bytes[44:48]),
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

	binary.LittleEndian.PutUint32(x.Bytes[0:4], s[0])
	binary.LittleEndian.PutUint32(x.Bytes[4:8], s[1])
	binary.LittleEndian.PutUint32(x.Bytes[8:12], s[2])
	binary.LittleEndian.PutUint32(x.Bytes[12:16], s[3])
	binary.LittleEndian.PutUint32(x.Bytes[16:20], s[4])
	binary.LittleEndian.PutUint32(x.Bytes[20:24], s[5])
	binary.LittleEndian.PutUint32(x.Bytes[24:28], s[6])
	binary.LittleEndian.PutUint32(x.Bytes[28:32], s[7])
	binary.LittleEndian.PutUint32(x.Bytes[32:36], s[8])
	binary.LittleEndian.PutUint32(x.Bytes[36:40], s[9])
	binary.LittleEndian.PutUint32(x.Bytes[40:44], s[10])
	binary.LittleEndian.PutUint32(x.Bytes[44:48], s[11])
}
