package brackets

func Bracket(str string) (bool, error) {
	brackets := Stack{}

	for _, bracketCode := range str {
		if bracketCode == 40 || bracketCode == 91 || bracketCode == 123 {
			brackets.Push(int(bracketCode))
			continue
		}

		lastOpenBracket := string(rune(brackets.Pop()))

		switch string(bracketCode) {
		case ")":
			if lastOpenBracket != "(" {
				return false, nil
			}

		case "]":
			if lastOpenBracket != "[" {
				return false, nil
			}

		case "}":
			if lastOpenBracket != "{" {
				return false, nil
			}
		}
	}

	return brackets.IsEmpty(), nil
}
