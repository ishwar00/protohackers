package utils

func IsAlphaNumeric(r rune) bool {
	return r >= 'a' && r <= 'z' ||
		r >= 'A' && r <= 'Z' ||
		r >= '0' && r <= '9'
}

func IsAlhpaNumericStr(str string) bool {
	for _, r := range str {
		if !IsAlphaNumeric(r) {
			return false
		}
	}

	return true
}

