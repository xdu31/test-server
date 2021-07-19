package util

import (
	"net"
	"strings"
)

const hexStr = "0123456789abcdef"

func IsDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	startPos := 0
	if strings.HasPrefix(s, "*.") {
		startPos += 2
	}
	for i := startPos; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
}

func IsIP(v string) bool {
	return (net.ParseIP(v) != nil)
}

func CheckIPv4Addr(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip != nil {
		if ip.To4() != nil {
			return true
		}
	}

	return false
}

func CheckIPv6Addr(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip != nil {
		if ip.To4() == nil {
			return true
		}
	}

	return false
}
