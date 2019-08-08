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
	xoodyak := Xoodyak{
		mode: modeKeyed,
		rates: Rates{
			absorb:  rateInput,
			squeeze: rateOutput,
		},
		phase: phaseUp,
	}

	var buf []byte
	buf = append(buf, key...)
	buf = append(buf, id...)
	buf = append(buf, byte(len(id)))
	if len(buf) > rateInput {
		panic("Key + ID too long!")
	}
	xoodyak.absorbAny(buf, xoodyak.rates.absorb, flagAbsorbKey)

	if len(counter) > 0 {
		xoodyak.absorbAny(counter, 1, flagZero)
	}

	return &xoodyak
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
			x.up(flagZero)
		}
		x.down(block, downFlag)
		downFlag = flagZero
	}
}

func (x *Xoodyak) crypt(in, out []byte, decrypt bool) []byte {
	flag := flagCrypt
	offset := 0
	for _, block := range blocks(in, rateOutput) {
		x.up(flag)
		flag = flagZero
		for i, b := range block {
			out = append(out, b^x.xoodoo.Get(i))
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
	out = x.upTo(out, min(count, x.rates.squeeze), upFlag)
	for len(out)-iLen < count {
		x.down(nil, flagZero)
		out = x.upTo(out, min(count-len(out)+iLen, x.rates.squeeze), flagZero)
	}
	return out
}

func (x *Xoodyak) down(block []byte, flag Flag) {
	x.phase = phaseDown
	for i, b := range block {
		x.xoodoo.XOR(i, b)
	}
	x.xoodoo.XOR(len(block), 0x01)
	if x.mode == modeHash {
		x.xoodoo.XOR(47, byte(flag)&0x01)
	} else {
		x.xoodoo.XOR(47, byte(flag))
	}
}

func (x *Xoodyak) up(flag Flag) {
	x.phase = phaseUp
	if x.mode != modeHash {
		x.xoodoo.XOR(47, byte(flag))
	}
	x.xoodoo.Permute()
}

func (x *Xoodyak) upTo(out []byte, count int, flag Flag) []byte {
	x.up(flag)
	for i := 0; i < count; i++ {
		out = append(out, x.xoodoo.Get(i))
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

func (x *Xoodyak) Decrypt(pt, ct []byte) []byte {
	if x.mode != modeKeyed {
		panic("Xoodyak not keyed!")
	}
	return x.crypt(pt, ct, true)
}

func (x *Xoodyak) Squeeze(out []byte, count int) []byte {
	return x.squeezeAny(out, count, flagSqueeze)
}

func (x *Xoodyak) SqueezeKey(out []byte, count int) []byte {
	return x.squeezeAny(out, count, flagSqueezeKey)
}

func (x *Xoodyak) Ratchet() {
	if x.mode != modeKeyed {
		panic("Xoodyak not keyed!")
	}
	buf := x.squeezeAny(nil, rateRatchet, flagRatchet)
	x.absorbAny(buf, x.rates.absorb, flagZero)
}
