package downcase

func Downcase(str string) (string, error) {
    res := ""

    for _, char := range str {
        if 65 <= char && char <= 90 {
            res += string(char + 32)
        } else {
            res += string(char)
        }
    }

    return res, nil
}
