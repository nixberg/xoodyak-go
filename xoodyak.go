package xoodyak

import "xoodyak/xoodoo"

type Flag byte

const (
	flagZero       Flag = 0x00
	flagAbsorbKey  Flag = 0x02
	flagAbsorb     Flag = 0x03
	flagRatchet    Flag = 0x10
	flagSqueezeKey Flag = 0x20
	flagSqueeze    Flag = 0x40
	flagCrypt      Flag = 0x80
)

type Mode int

const (
	modeHash Mode = iota
	modeKeyed
)

type Rates struct {
	absorb  int
	squeeze int
}

const (
	rateHash    = 16
	rateInput   = 44
	rateOutput  = 24
	rateRatchet = 16
)

type Phase int

const (
	phaseUp Phase = iota
	phaseDown
)

type Xoodyak struct {
	mode   Mode
	rates  Rates
	phase  Phase
	xoodoo xoodoo.Xoodoo
}

func New() *Xoodyak {
	return &Xoodyak{
		mode: modeHash,
		rates: Rates{
			absorb:  rateHash,
			squeeze: rateHash,
		},
		phase: phaseUp,
	}
}

func Keyed(key, id, counter []byte) *Xoodyak {
	xoodyak := New()
	xoodyak.absorbKey(key, id, counter)
	return xoodyak
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func blocks(data []byte, rate int) (blocks [][]byte) {
	dataLen := len(data)
	if dataLen == 0 {
		blocks = append(blocks, []byte{})
		return
	}
	for start := 0; start < dataLen; start += rate {
		end := min(start+rate, dataLen)
		blocks = append(blocks, data[start:end])
	}
	return
}

func (x *Xoodyak) absorbAny(data []byte, rate int, downFlag Flag) {
	for _, block := range blocks(data, rate) {
		if x.phase != phaseUp {
			x.up(nil, 0, flagZero)
		}
		x.down(block, downFlag)
		downFlag = flagZero
	}
}

func (x *Xoodyak) absorbKey(key, id, counter []byte) {
	x.mode = modeKeyed
	x.rates = Rates{
		absorb:  rateInput,
		squeeze: rateOutput,
	}

	var buf []byte
	buf = append(buf, key...)
	buf = append(buf, id...)
	if len(buf) > rateInput-1 {
		panic("Key + ID too long!")
	}
	buf = append(buf, byte(len(id)))

	x.absorbAny(buf, x.rates.absorb, flagAbsorbKey)

	if len(counter) > 0 {
		x.absorbAny(counter, 1, flagZero)
	}
}

func (x *Xoodyak) crypt(in, out []byte, decrypt bool) []byte {
	flag := flagCrypt
	offset := len(out)
	for _, block := range blocks(in, rateOutput) {
		x.up(nil, 0, flag)
		flag = flagZero
		for i, b := range block {
			out = append(out, b^x.xoodoo.Bytes[i])
		}
		if decrypt {
			x.down(out[offset:offset+len(block)], flagZero)
		} else {
			x.down(block, flagZero)
		}
		offset += len(block)
	}
	return out
}

func (x *Xoodyak) squeezeAny(out []byte, count int, upFlag Flag) []byte {
	iLen := len(out)
	out = x.up(out, min(count, x.rates.squeeze), upFlag)
	for len(out)-iLen < count {
		x.down(nil, flagZero)
		out = x.up(out, min(count-len(out)+iLen, x.rates.squeeze), flagZero)
	}
	return out
}

func (x *Xoodyak) down(block []byte, flag Flag) {
	x.phase = phaseDown
	for i, b := range block {
		x.xoodoo.Bytes[i] ^= b
	}
	x.xoodoo.Bytes[len(block)] ^= 0x01
	if x.mode == modeHash {
		x.xoodoo.Bytes[47] ^= byte(flag) & 0x01
	} else {
		x.xoodoo.Bytes[47] ^= byte(flag)
	}
}

func (x *Xoodyak) up(out []byte, count int, flag Flag) []byte {
	x.phase = phaseUp
	if x.mode != modeHash {
		x.xoodoo.Bytes[47] ^= byte(flag)
	}
	x.xoodoo.Permute()
	for i := 0; i < count; i++ {
		out = append(out, x.xoodoo.Bytes[i])
	}
	return out
}

func (x *Xoodyak) Absorb(in []byte) {
	x.absorbAny(in, x.rates.absorb, flagAbsorb)
}

func (x *Xoodyak) Encrypt(pt, ct []byte) []byte {
	if x.mode != modeKeyed {
		panic("Xoodyak not keyed!")
	}
	return x.crypt(pt, ct, false)
}

func (x *Xoodyak) Decrypt(ct, pt []byte) []byte {
	if x.mode != modeKeyed {
		panic("Xoodyak not keyed!")
	}
	return x.crypt(ct, pt, true)
}

func (x *Xoodyak) Squeeze(out []byte, count int) []byte {
	return x.squeezeAny(out, count, flagSqueeze)
}

func (x *Xoodyak) SqueezeKey(out []byte, count int) []byte {
	if x.mode != modeKeyed {
		panic("Xoodyak not keyed!")
	}
	return x.squeezeAny(out, count, flagSqueezeKey)
}

func (x *Xoodyak) Ratchet() {
	if x.mode != modeKeyed {
		panic("Xoodyak not keyed!")
	}
	buf := x.squeezeAny(nil, rateRatchet, flagRatchet)
	x.absorbAny(buf, x.rates.absorb, flagZero)
}
