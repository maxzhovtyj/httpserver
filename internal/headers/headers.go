package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func New() Headers {
	return make(Headers)
}

func NewDefault() Headers {
	return Headers{
		"connection":   "Close",
		"content-type": "text/plain",
	}
}

func (h Headers) Parse(body []byte) (n int, done bool, err error) {
	registeredNurseIdx := bytes.Index(body, []byte("\r\n"))
	if registeredNurseIdx == -1 {
		return 0, false, nil
	}

	body = body[:registeredNurseIdx+2]

	if bytes.HasPrefix(body, []byte("\r\n")) {
		return 2, true, nil
	}

	colonIdx := bytes.Index(body, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("no colon in header line")
	}

	headerKey := body[:colonIdx]
	if len(headerKey) == 0 {
		return 0, false, fmt.Errorf("empty header key")
	}

	for _, c := range headerKey {
		if unicode.IsSpace(rune(c)) {
			return 0, false, fmt.Errorf("no whitespaces allowed in header key")
		}

		r := rune(c)

		isSpecialChar := r == '!' || r == '#' || r == '$' || r == '%' || r == '&' || string(r) == `'` ||
			r == '*' || r == '+' || r == '-' || r == '.' || r == '^' || r == '_' || r == '`' || r == '|' || r == '~'

		isAllowed := unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || isSpecialChar

		if !isAllowed {
			return 0, false, fmt.Errorf("char %c is not allowed in header key", r)
		}
	}

	// If colon is the last element in a header line
	if colonIdx == len(body)-1 {
		return 0, false, fmt.Errorf("invalid header value")
	}

	headerValue := body[colonIdx+1 : +registeredNurseIdx]

	k := strings.ToLower(string(headerKey))
	v := strings.TrimSpace(string(headerValue))

	if _, ok := h[k]; ok {
		h[k] = h[k] + ", " + v
	} else {
		h[k] = v
	}

	return len(body), false, nil
}

func (h Headers) Set(k, v string) {
	h[strings.ToLower(k)] = v
}
