package shenma

func deutf8(str string) string {
	if str == "" {
		return ""
	}

	var buf []rune
	ss := []rune(str)
	for e := 0; e < len(ss); e++ {
		buf = append(buf, 1^ss[e])
	}
	return string(buf)
}
