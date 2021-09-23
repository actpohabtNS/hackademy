package cipher

type Cipher interface {
	Encode(string) string
	Decode(string) string
}

func _encode(str string, shift int32) string {
	if shift < 0 {
		shift += 26
	}

	res := ""

	str, _ = Downcase(str)

	for _, char := range str {
		if char < 97 || char > 122 {
			continue
		}

		res += (string)((char-97+shift)%26 + 97)
	}

	return res
}

type caesarCipher struct{}

func (c caesarCipher) Encode(str string) string {
	return _encode(str, 3)
}

func (c caesarCipher) Decode(str string) string {
	return _encode(str, -3)
}

func NewCaesar() Cipher {
	return caesarCipher{}
}

type shiftCipher struct {
	shift int32
}

func (s shiftCipher) Encode(str string) string {
	return _encode(str, s.shift)
}

func (s shiftCipher) Decode(str string) string {
	return _encode(str, -s.shift)
}

func NewShift(shift int) Cipher {
	if shift > 25 || shift < -25 || shift == 0 {
		return nil
	}
	return shiftCipher{int32(shift)}
}

type vigenereCipher struct {
	key string
}

func _vigOperation(str string, key string, add bool) string {
	res := ""

	str, _ = Downcase(str)

	var modifier int32 = 1

	if !add {
		modifier = -1
	}

	keyI := 0
	for _, char := range str {
		if char < 97 || char > 122 {
			continue
		}

		shift := int(key[keyI]) - 97
		res += _encode(string(char), modifier*int32(shift))

		if keyI+1 == len(key) {
			keyI = 0
		} else {
			keyI++
		}
	}

	return res
}

func (v vigenereCipher) Encode(str string) string {
	return _vigOperation(str, v.key, true)
}

func (v vigenereCipher) Decode(str string) string {
	return _vigOperation(str, v.key, false)
}

func vigenereKeyIsCorrect(key string) bool {
	notAFound := false

	for _, char := range key {
		if char < 97 || char > 122 {
			return false
		}

		if char != 97 {
			notAFound = true
		}
	}

	return notAFound
}

func NewVigenere(key string) Cipher {
	if !vigenereKeyIsCorrect(key) {
		return nil
	}

	return vigenereCipher{key}
}
