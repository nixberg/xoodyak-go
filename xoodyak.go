package xoodyak

import "github.com/nixberg/xoodyak-go/internal/xoodoo"

const (
	flagZero       byte = 0x00
	flagAbsorbKey  byte = 0x02
	flagAbsorb     byte = 0x03
	flagRatchet    byte = 0x10
	flagSqueezeKey byte = 0x20
	flagSqueeze    byte = 0x40
	flagCrypt      byte = 0x80

	rateHash        int = 16
	rateKeyedInput  int = 44
	rateKeyedOutput int = 24
	rateRatchet     int = 16
	rateCounter     int = 1
)

type Xoodyak struct {
	state xoodoo.State
	rates struct {
		absorb  int
		squeeze int
	}
	isPhaseUp  bool
	isModeHash bool
}

func New() *Xoodyak {
	return &Xoodyak{
		rates: struct {
			absorb  int
			squeeze int
		}{
			absorb:  rateHash,
			squeeze: rateHash,
		},
		isPhaseUp:  true,
		isModeHash: true,
	}
}

func (x *Xoodyak) down(block []byte, flag byte) {
	x.isPhaseUp = false

	for i, b := range block {
		x.state[i] ^= b
	}

	x.state[len(block)] ^= 0x01
	if x.isModeHash {
		x.state[47] ^= flag & 0x01
	} else {
		x.state[47] ^= flag
	}
}

func (x *Xoodyak) up(output []byte, count int, flag byte) []byte {
	x.isPhaseUp = true
	if !x.isModeHash {
		x.state[47] ^= flag
	}
	x.state.Permute()

	for i := 0; i < count; i++ {
		output = append(output, x.state[i])
	}
	return output
}

func (x *Xoodyak) absorbAny(input []byte, rate int, downFlag byte) {
	for {
		block := input[:min(rate, len(input))]
		input = input[len(block):]

		if !x.isPhaseUp {
			x.up(nil, 0, flagZero)
		}

		x.down(block, downFlag)
		downFlag = flagZero

		if len(input) == 0 {
			break
		}
	}
}

func (x *Xoodyak) squeezeAny(output []byte, count int, upFlag byte) []byte {
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
