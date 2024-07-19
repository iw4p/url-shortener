package base62

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Base62 struct{}

func (b Base62) EncodeBase62(number int) string {
	if number == 0 {
		return string(base62Digits[0])
	}
	base62 := ""
	for number > 0 {
		remainder := number % 62
		base62 = string(base62Digits[remainder]) + base62
		number /= 62
	}
	for len(base62) < 7 {
		base62 = string(base62Digits[0]) + base62
	}
	return base62
}

func (b Base62) DecodeBase62(value string) int {
	var number int
	for _, char := range value {
		number = number*62 + indexOf(char)
	}
	return number
}

func indexOf(char rune) int {
	for i, c := range base62Digits {
		if c == char {
			return i
		}
	}
	return -1
}
