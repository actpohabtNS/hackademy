package cipher

type Cipher interface {
	Encode(string) string
	Decode(string) string
}

func encode(str string, shift int32) string {
	if shift < 0 {
		shift += 26
	}

	res := ""

	str, _ = Downcase(str)

	for _, char := range str {
		if char < 'a' || char > 'z' {
			continue
		}

		res += (string)((char-'a'+shift)%26 + 'a')
	}

	return res
}

type caesarCipher struct{}

func (c caesarCipher) Encode(str string) string {
	return encode(str, 3)
}

func (c caesarCipher) Decode(str string) string {
	return encode(str, -3)
}

func NewCaesar() Cipher {
	return caesarCipher{}
}

type shiftCipher struct {
	shift int32
}

func (s shiftCipher) Encode(str string) string {
	return encode(str, s.shift)
}

func (s shiftCipher) Decode(str string) string {
	return encode(str, -s.shift)
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

func vigenereEncode(str string, key string, add bool) string {
	res := ""

	str, _ = Downcase(str)

	var modifier int32 = 1

	if !add {
		modifier = -1
	}

	keyI := 0
	for _, char := range str {
		if char < 'a' || char > 'z' {
			continue
		}

		shift := int(key[keyI]) - 'a'
		res += encode(string(char), modifier*int32(shift))

		keyI = (keyI + 1) % len(key)
	}

	return res
}

func (v vigenereCipher) Encode(str string) string {
	return vigenereEncode(str, v.key, true)
}

func (v vigenereCipher) Decode(str string) string {
	return vigenereEncode(str, v.key, false)
}

func vigenereKeyIsCorrect(key string) bool {
	notAFound := false

	for _, char := range key {
		if char < 'a' || char > 'z' {
			return false
		}

		if char != 'a' {
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
