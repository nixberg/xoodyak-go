package xoodyak

type KeyedXoodyak struct {
	xoodyak Xoodyak
}

func NewKeyed(key, id, counter []byte) *KeyedXoodyak {
	if len(key) == 0 {
		panic("xoodyak: key is empty")
	}

	x := KeyedXoodyak{
		Xoodyak{
			rates: struct {
				absorb  int
				squeeze int
			}{
				absorb:  rateKeyedInput,
				squeeze: rateKeyedOutput,
			},
			isPhaseUp:  true,
			isModeHash: false,
		},
	}

	var buffer []byte
	buffer = append(buffer, key...)
	buffer = append(buffer, id...)
	buffer = append(buffer, byte(len(id)))
	if !(len(buffer) <= rateKeyedInput) {
		panic("xoodyak: length of key+id exceeds 43 bytes")
	}

	x.xoodyak.absorbAny(buffer, x.xoodyak.rates.absorb, flagAbsorbKey)

	if len(counter) > 0 {
		x.xoodyak.absorbAny(counter, rateCounter, flagZero)
	}

	return &x
}

func (x *KeyedXoodyak) crypt(input, output []byte, decrypt bool) []byte {
	flag := flagCrypt
	offset := len(output)

	for {
		block := input[:min(rateKeyedOutput, len(input))]
		input = input[len(block):]

		x.xoodyak.up(nil, 0, flag)
		flag = flagZero

		for i, b := range block {
			output = append(output, b^x.xoodyak.state[i])
		}

		if decrypt {
			x.xoodyak.down(output[offset:offset+len(block)], flagZero)
		} else {
			x.xoodyak.down(block, flagZero)
		}

		offset += len(block)

		if len(input) == 0 {
			break
		}
	}
	return output
}

func (x *KeyedXoodyak) Absorb(input []byte) {
	x.xoodyak.Absorb(input)
}

func (x *KeyedXoodyak) Encrypt(plaintext, ciphertext []byte) []byte {
	return x.crypt(plaintext, ciphertext, false)
}

func (x *KeyedXoodyak) Decrypt(ciphertext, plaintext []byte) []byte {
	return x.crypt(ciphertext, plaintext, true)
}

func (x *KeyedXoodyak) Squeeze(output []byte, count int) []byte {
	return x.xoodyak.Squeeze(output, count)
}

func (x *KeyedXoodyak) SqueezeKey(output []byte, count int) []byte {
	return x.xoodyak.squeezeAny(output, count, flagSqueezeKey)
}

func (x *KeyedXoodyak) Ratchet() {
	buffer := x.xoodyak.squeezeAny(nil, rateRatchet, flagRatchet)
	x.xoodyak.absorbAny(buffer, x.xoodyak.rates.absorb, flagZero)
}
