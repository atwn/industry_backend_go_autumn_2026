package main

func rotateRunes(s string, shift int) string {
	runes := []rune(s)
	length := len(runes)
	if length == 0 {
		return ""
	}

	offset := shift % length
	if offset < 0 {
		// нормализация остатка при отрицательном сдвиге
		// (прибавляем длину, чтобы не делать сдвиг в обратную сторону по модулю)
		offset = offset + length
	}

	return string(append(runes[offset:], runes[:offset]...))
}
