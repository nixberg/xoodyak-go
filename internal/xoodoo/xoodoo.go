package xoodoo

import (
	"encoding/binary"
	"math/bits"
)

type State [48]byte

func (state *State) Permute() {
	s0 := binary.LittleEndian.Uint32(state[0:4])
	s1 := binary.LittleEndian.Uint32(state[4:8])
	s2 := binary.LittleEndian.Uint32(state[8:12])
	s3 := binary.LittleEndian.Uint32(state[12:16])
	s4 := binary.LittleEndian.Uint32(state[16:20])
	s5 := binary.LittleEndian.Uint32(state[20:24])
	s6 := binary.LittleEndian.Uint32(state[24:28])
	s7 := binary.LittleEndian.Uint32(state[28:32])
	s8 := binary.LittleEndian.Uint32(state[32:36])
	s9 := binary.LittleEndian.Uint32(state[36:40])
	s10 := binary.LittleEndian.Uint32(state[40:44])
	s11 := binary.LittleEndian.Uint32(state[44:48])

	for _, roundConstant := range [12]uint32{
		0x058, 0x038, 0x3c0, 0x0d0,
		0x120, 0x014, 0x060, 0x02c,
		0x380, 0x0f0, 0x1a0, 0x012,
	} {
		p0 := s0 ^ s4 ^ s8
		e0 := bits.RotateLeft32(p0, 5) ^ bits.RotateLeft32(p0, 14)

		p1 := s1 ^ s5 ^ s9
		e1 := bits.RotateLeft32(p1, 5) ^ bits.RotateLeft32(p1, 14)

		p2 := s2 ^ s6 ^ s10
		e2 := bits.RotateLeft32(p2, 5) ^ bits.RotateLeft32(p2, 14)

		p3 := s3 ^ s7 ^ s11
		e3 := bits.RotateLeft32(p3, 5) ^ bits.RotateLeft32(p3, 14)

		s0 ^= e3
		s1 ^= e0
		s2 ^= e1
		s3 ^= e2
		s4 ^= e3
		s5 ^= e0
		s6 ^= e1
		s7 ^= e2
		s8 ^= e3
		s9 ^= e0
		s10 ^= e1
		s11 ^= e2

		s7, s4 = s4, s7
		s7, s5 = s5, s7
		s7, s6 = s6, s7

		s0 ^= roundConstant

		{
			a := s0
			b := s4
			c := bits.RotateLeft32(s8, 11)

			s8 = bits.RotateLeft32((b&^a)^c, 8)
			s4 = bits.RotateLeft32((a&^c)^b, 1)
			s0 ^= c & ^b
		}

		{
			a := s1
			b := s5
			c := bits.RotateLeft32(s9, 11)

			s9 = bits.RotateLeft32((b&^a)^c, 8)
			s5 = bits.RotateLeft32((a&^c)^b, 1)
			s1 ^= c & ^b
		}

		{
			a := s2
			b := s6
			c := bits.RotateLeft32(s10, 11)

			s10 = bits.RotateLeft32((b&^a)^c, 8)
			s6 = bits.RotateLeft32((a&^c)^b, 1)
			s2 ^= c & ^b
		}

		{
			a := s3
			b := s7
			c := bits.RotateLeft32(s11, 11)

			s11 = bits.RotateLeft32((b&^a)^c, 8)
			s7 = bits.RotateLeft32((a&^c)^b, 1)
			s3 ^= c & ^b
		}

		s8, s10 = s10, s8
		s9, s11 = s11, s9
	}

	binary.LittleEndian.PutUint32(state[0:4], s0)
	binary.LittleEndian.PutUint32(state[4:8], s1)
	binary.LittleEndian.PutUint32(state[8:12], s2)
	binary.LittleEndian.PutUint32(state[12:16], s3)
	binary.LittleEndian.PutUint32(state[16:20], s4)
	binary.LittleEndian.PutUint32(state[20:24], s5)
	binary.LittleEndian.PutUint32(state[24:28], s6)
	binary.LittleEndian.PutUint32(state[28:32], s7)
	binary.LittleEndian.PutUint32(state[32:36], s8)
	binary.LittleEndian.PutUint32(state[36:40], s9)
	binary.LittleEndian.PutUint32(state[40:44], s10)
	binary.LittleEndian.PutUint32(state[44:48], s11)
}
