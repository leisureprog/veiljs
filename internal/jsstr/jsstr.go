package jsstr

func hexVal(c byte) (int, bool) {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0'), true
	case 'a' <= c && c <= 'f':
		return int(c-'a') + 10, true
	case 'A' <= c && c <= 'F':
		return int(c-'A') + 10, true
	}
	return 0, false
}

func hexd(v byte) byte {
	if v < 10 {
		return '0' + v
	}
	return 'a' + v - 10
}

func appendRune(b []byte, r rune) []byte {
	var buf [4]byte
	n := encodeRune(buf[:], r)
	return append(b, buf[:n]...)
}

func encodeRune(p []byte, r rune) int {
	switch {
	case r < 0x80:
		p[0] = byte(r)
		return 1
	case r < 0x800:
		p[0] = 0xC0 | byte(r>>6)
		p[1] = 0x80 | byte(r)&0x3F
		return 2
	case r < 0x10000:
		p[0] = 0xE0 | byte(r>>12)
		p[1] = 0x80 | byte(r>>6)&0x3F
		p[2] = 0x80 | byte(r)&0x3F
		return 3
	default:
		p[0] = 0xF0 | byte(r>>18)
		p[1] = 0x80 | byte(r>>12)&0x3F
		p[2] = 0x80 | byte(r>>6)&0x3F
		p[3] = 0x80 | byte(r)&0x3F
		return 4
	}
}

func decodeHex(s []byte) (uint32, bool) {
	var v uint32
	for _, c := range s {
		h, ok := hexVal(c)
		if !ok {
			return 0, false
		}
		v = v<<4 | uint32(h)
	}
	return v, true
}

func Unquote(lit []byte) string {
	if len(lit) < 2 {
		return string(lit)
	}
	if lit[0] != '"' && lit[0] != '\'' && lit[0] != '`' {
		return string(lit)
	}
	s := lit[1 : len(lit)-1]
	var out []byte
	for i := 0; i < len(s); {
		c := s[i]
		if c != '\\' {
			out = append(out, c)
			i++
			continue
		}
		i++
		if i >= len(s) {
			break
		}
		switch e := s[i]; e {
		case 'n':
			out = append(out, '\n')
			i++
		case 't':
			out = append(out, '\t')
			i++
		case 'r':
			out = append(out, '\r')
			i++
		case 'b':
			out = append(out, '\b')
			i++
		case 'f':
			out = append(out, '\f')
			i++
		case 'v':
			out = append(out, '\v')
			i++
		case '0', '1', '2', '3', '4', '5', '6', '7':
			v := uint32(e - '0')
			i++
			for j := 0; j < 2 && i < len(s) && '0' <= s[i] && s[i] <= '7'; j++ {
				v = v<<3 | uint32(s[i]-'0')
				i++
			}
			out = appendRune(out, rune(v))
		case 'x':
			if i+2 < len(s) {
				if v, ok := decodeHex(s[i+1 : i+3]); ok {
					out = appendRune(out, rune(v))
				}
				i += 3
			} else {
				i++
			}
		case 'u':
			if i+1 < len(s) && s[i+1] == '{' {
				j := i + 2
				for j < len(s) && s[j] != '}' {
					j++
				}
				if j <= len(s) && j-1 > i+2 {
					if v, ok := decodeHex(s[i+2 : j]); ok {
						out = appendRune(out, rune(v))
					}
				}
				i = j + 1
			} else if i+4 < len(s) {
				v, ok := decodeHex(s[i+1 : i+5])
				i += 5
				if ok {
					r := rune(v)
					if 0xD800 <= v && v <= 0xDBFF && i+5 < len(s) && s[i] == '\\' && s[i+1] == 'u' {
						if lo, ok := decodeHex(s[i+2 : i+6]); ok && 0xDC00 <= lo && lo <= 0xDFFF {
							r = ((r - 0xD800) << 10) + (rune(lo) - 0xDC00) + 0x10000
							i += 6
						}
					}
					out = appendRune(out, r)
				}
			} else {
				i++
			}
		case '\r':
			i++
			if i < len(s) && s[i] == '\n' {
				i++
			}
		case '\n':
			i++
		default:
			out = append(out, s[i])
			i++
		}
	}
	return string(out)
}

func Quote(s string) []byte {
	out := make([]byte, 0, len(s)+2)
	out = append(out, '"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			out = append(out, '\\', '"')
		case '\\':
			out = append(out, '\\', '\\')
		case '\n':
			out = append(out, '\\', 'n')
		case '\r':
			out = append(out, '\\', 'r')
		case '\t':
			out = append(out, '\\', 't')
		case '\b':
			out = append(out, '\\', 'b')
		case '\f':
			out = append(out, '\\', 'f')
		default:
			if c < 0x20 {
				out = append(out, '\\', 'x', hexd(c>>4), hexd(c&15))
			} else {
				out = append(out, c)
			}
		}
	}
	return append(out, '"')
}

func EscapeUnicode(lit []byte) []byte {
	s := Unquote(lit)
	out := make([]byte, 0, len(s)+2)
	out = append(out, '"')
	for _, r := range s {
		switch {
		case r == '"':
			out = append(out, '\\', '"')
		case r == '\\':
			out = append(out, '\\', '\\')
		case r == '\n':
			out = append(out, '\\', 'n')
		case r == '\r':
			out = append(out, '\\', 'r')
		case r == '\t':
			out = append(out, '\\', 't')
		case r < 0x20:
			out = append(out, '\\', 'x', hexd(byte(r)>>4), hexd(byte(r)&15))
		case r < 0x80:
			out = append(out, byte(r))
		case r <= 0xFF:
			out = append(out, '\\', 'x', hexd(byte(r)>>4), hexd(byte(r)&15))
		case r <= 0xFFFF:
			out = appendU(out, uint16(r))
		default:
			v := r - 0x10000
			out = appendU(out, uint16(0xD800+(v>>10)))
			out = appendU(out, uint16(0xDC00+(v&0x3FF)))
		}
	}
	return append(out, '"')
}

func appendU(out []byte, v uint16) []byte {
	return append(out, '\\', 'u',
		hexd(byte(v>>12)&15),
		hexd(byte(v>>8)&15),
		hexd(byte(v>>4)&15),
		hexd(byte(v)&15))
}
