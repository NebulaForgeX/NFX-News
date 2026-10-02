package getters

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func decodeGBK(raw []byte) string {
	out, err := io.ReadAll(transform.NewReader(bytes.NewReader(raw), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return string(raw)
	}
	return string(out)
}

func extractJSONObject(s, marker string) ([]byte, error) {
	i := strings.Index(s, marker)
	if i < 0 {
		return nil, fmt.Errorf("%s not found", marker)
	}
	i += len(marker)
	for i < len(s) && (s[i] == ' ' || s[i] == '\n' || s[i] == '\r' || s[i] == '\t' || s[i] == '=') {
		i++
	}
	if i >= len(s) || s[i] != '{' {
		return nil, fmt.Errorf("%s is not a JSON object", marker)
	}
	depth, inStr, esc := 0, false, false
	for j := i; j < len(s); j++ {
		c := s[j]
		if inStr {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return []byte(s[i : j+1]), nil
			}
		}
	}
	return nil, fmt.Errorf("%s JSON object unclosed", marker)
}
