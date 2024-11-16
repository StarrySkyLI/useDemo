package runeStr

import "unicode/utf8"

func GetUtf8Str(data string, start int, end int) string {
	lenData := utf8.RuneCountInString(data)

	if lenData < start || lenData < end {
		return ""
	}

	// 将字符串转换为rune切片
	runes := []rune(data)

	startIndex := utf8.RuneCountInString(string(runes[:start]))
	endIndex := utf8.RuneCountInString(string(runes[:end]))

	return string(runes[startIndex:endIndex])
}
