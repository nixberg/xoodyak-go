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

func (x *Xoodyak) up(output []byte, count int, flag xoodyakFlag) []byte {
	x.phase = phaseUp
	if x.mode != modeHash {
		x.state.Bytes[47] ^= byte(flag)
	}
	x.state.Permute()

	for i := 0; i < count; i++ {
		output = append(output, x.state.Bytes[i])
	}
	return output
}

func (x *Xoodyak) absorbAny(input []byte, rate xoodyakRate, downFlag xoodyakFlag) {
	for {
		block := input[:min(rate, len(input))]
		input = input[len(block):]

		if x.phase != phaseUp {
			x.up(nil, 0, flagZero)
		}

		x.down(block, downFlag)
		downFlag = flagZero

		if len(input) == 0 {
			break
		}
	}
}

func (x *Xoodyak) squeezeAny(output []byte, count int, upFlag xoodyakFlag) []byte {
	blockSize := min(x.rates.squeeze, count)
	count -= blockSize

	output = x.up(output, blockSize, upFlag)

	for count > 0 {
		blockSize = min(x.rates.squeeze, count)
		count -= blockSize

		x.down(nil, flagZero)
		output = x.up(output, blockSize, flagZero)
	}

	return output
}

func (x *Xoodyak) Absorb(input []byte) {
	x.absorbAny(input, x.rates.absorb, flagAbsorb)
}

func (x *Xoodyak) Squeeze(output []byte, count int) []byte {
	return x.squeezeAny(output, count, flagSqueeze)
}

func min(rate xoodyakRate, b int) int {
	a := int(rate)
	if a < b {
		return a
	}
	return b
}
