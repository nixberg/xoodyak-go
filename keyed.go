package xoodyak

type KeyedXoodyak struct {
	xoodyak Xoodyak
}

const (
	rateKeyedInput  xoodyakRate = 44
	rateKeyedOutput xoodyakRate = 24
	rateRatchet     xoodyakRate = 16
	rateCounter     xoodyakRate = 1
)

func NewKeyed(key, id, counter []byte) *KeyedXoodyak {
	if len(key) == 0 {
		panic("xoodyak: key is empty")
	}

	x := KeyedXoodyak{
		Xoodyak{
			phase: phaseUp,
			mode:  modeKeyed,
			rates: xoodyakRates{
				absorb:  rateKeyedInput,
				squeeze: rateKeyedOutput,
			},
		},
	}

	var buffer []byte
	buffer = append(buffer, key...)
	buffer = append(buffer, id...)
	buffer = append(buffer, byte(len(id)))
	if !(len(buffer) <= int(rateKeyedInput)) {
		panic("xoodyak: length key and id exceeds 43 bytes")
	}

	x.xoodyak.absorbAny(buffer, x.xoodyak.rates.absorb, flagAbsorbKey)

	if len(counter) > 0 {
		x.xoodyak.absorbAny(counter, rateCounter, flagZero)
	}

	return &x
}

func (x *KeyedXoodyak) crypt(in, out []byte, decrypt bool) []byte {
	flag := flagCrypt
	offset := len(out)
	for _, block := range blocks(in, rateKeyedOutput) {
		x.xoodyak.up(nil, 0, flag)
		flag = flagZero
		for i, b := range block {
			out = append(out, b^x.xoodyak.state.Bytes[i])
		}
		if decrypt {
			x.xoodyak.down(out[offset:offset+len(block)], flagZero)
		} else {
			x.xoodyak.down(block, flagZero)
		}
		offset += len(block)
	}
	return out
}

func (x *KeyedXoodyak) Absorb(in []byte) {
	x.xoodyak.Absorb(in)
}

func (x *KeyedXoodyak) Encrypt(pt, ct []byte) []byte {
	return x.crypt(pt, ct, false)
}

func (x *KeyedXoodyak) Decrypt(ct, pt []byte) []byte {
	return x.crypt(ct, pt, true)
}

func (x *KeyedXoodyak) Squeeze(out []byte, count int) []byte {
	return x.xoodyak.Squeeze(out, count)
}

func (x *KeyedXoodyak) SqueezeKey(out []byte, count int) []byte {
	return x.xoodyak.squeezeAny(out, count, flagSqueezeKey)
}

func (x *KeyedXoodyak) Ratchet() {
	buffer := x.xoodyak.squeezeAny(nil, int(rateRatchet), flagRatchet)
	x.xoodyak.absorbAny(buffer, x.xoodyak.rates.absorb, flagZero)
}
