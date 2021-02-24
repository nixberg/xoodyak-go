package xoodyak

import "github.com/nixberg/xoodyak-go/internal/xoodoo"

type xoodyakPhase int

const (
	phaseUp xoodyakPhase = iota
	phaseDown
)

type xoodyakMode int

const (
	modeHash xoodyakMode = iota
	modeKeyed
)

type xoodyakRate int

const (
	rateHash xoodyakRate = 16
)

type xoodyakRates struct {
	absorb  xoodyakRate
	squeeze xoodyakRate
}

type xoodyakFlag byte

const (
	flagZero       xoodyakFlag = 0x00
	flagAbsorbKey  xoodyakFlag = 0x02
	flagAbsorb     xoodyakFlag = 0x03
	flagRatchet    xoodyakFlag = 0x10
	flagSqueezeKey xoodyakFlag = 0x20
	flagSqueeze    xoodyakFlag = 0x40
	flagCrypt      xoodyakFlag = 0x80
)

type Xoodyak struct {
	phase xoodyakPhase
	state xoodoo.Xoodoo
	mode  xoodyakMode
	rates xoodyakRates
}

func New() *Xoodyak {
	return &Xoodyak{
		phase: phaseUp,
		mode:  modeHash,
		rates: xoodyakRates{
			absorb:  rateHash,
			squeeze: rateHash,
		},
	}
}

func blocks(data []byte, rate xoodyakRate) (blocks [][]byte) {
	dataLen := len(data)
	if dataLen == 0 {
		blocks = append(blocks, []byte{})
		return
	}
	for start := 0; start < dataLen; start += int(rate) {
		end := min(start+int(rate), dataLen)
		blocks = append(blocks, data[start:end])
	}
	return
}

func (x *Xoodyak) absorbAny(data []byte, rate xoodyakRate, downFlag xoodyakFlag) {
	for _, block := range blocks(data, rate) {
		if x.phase != phaseUp {
			x.up(nil, 0, flagZero)
		}
		x.down(block, downFlag)
		downFlag = flagZero
	}
}

func (x *Xoodyak) squeezeAny(out []byte, count int, upFlag xoodyakFlag) []byte {
	iLen := len(out)
	out = x.up(out, min(count, int(x.rates.squeeze)), upFlag)
	for len(out)-iLen < count {
		x.down(nil, flagZero)
		out = x.up(out, min(count-len(out)+iLen, int(x.rates.squeeze)), flagZero)
	}
	return out
}

func (x *Xoodyak) down(block []byte, flag xoodyakFlag) {
	x.phase = phaseDown
	for i, b := range block {
		x.state.Bytes[i] ^= b
	}
	x.state.Bytes[len(block)] ^= 0x01
	if x.mode == modeHash {
		x.state.Bytes[47] ^= byte(flag) & 0x01
	} else {
		x.state.Bytes[47] ^= byte(flag)
	}
}

func (x *Xoodyak) up(out []byte, count int, flag xoodyakFlag) []byte {
	x.phase = phaseUp
	if x.mode != modeHash {
		x.state.Bytes[47] ^= byte(flag)
	}
	x.state.Permute()
	for i := 0; i < count; i++ {
		out = append(out, x.state.Bytes[i])
	}
	return out
}

func (x *Xoodyak) Absorb(in []byte) {
	x.absorbAny(in, x.rates.absorb, flagAbsorb)
}

func (x *Xoodyak) Squeeze(out []byte, count int) []byte {
	return x.squeezeAny(out, count, flagSqueeze)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
