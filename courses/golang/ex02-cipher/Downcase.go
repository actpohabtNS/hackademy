package cipher

func Downcase(str string) (string, error) {
	res := ""

	for _, char := range str {
		if 'A' <= char && char <= 'Z' {
			res += string(char + 32)
		} else {
			res += string(char)
		}
	}

	return res, nil
}
