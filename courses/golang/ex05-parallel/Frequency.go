package letter

func charFrequency(str string, ch chan map[string]int) {
	freqMap := make(map[string]int)

	for _, char := range str {
		strCh, _ := Downcase(string(char))
		freqMap[strCh]++
	}

	ch <- freqMap
}

func Frequency(str string) map[string]int {
	ch := make(chan map[string]int)
	go charFrequency(str, ch)
	return <-ch
}

func ConcurrentFrequency(strings []string) map[string]int {
	ch := make(chan map[string]int)

	for i := 0; i < len(strings); i++ {
		go charFrequency(strings[i], ch)
	}

	freqMapRes := make(map[string]int)

	for i := 0; i < len(strings); i++ {
		freqMap := <-ch
		for letter, freq := range freqMap {
			freqMapRes[letter] += freq
		}
	}

	return freqMapRes
}
