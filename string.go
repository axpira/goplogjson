package goplogjson

// string file has copyright of ZeroLog
// this file is full copied from ZeroLog lib

import (
	"unicode/utf8"
)

const hex = "0123456789abcdef"

// noEscapeTable is the same as this init method
// func init() {
// 	for i := 0; i <= 0x7e; i++ {
// 		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
// 	}
// }
var noEscapeTable = [256]bool{
	0:   false,
	1:   false,
	2:   false,
	3:   false,
	4:   false,
	5:   false,
	6:   false,
	7:   false,
	8:   false,
	9:   false,
	10:  false,
	11:  false,
	12:  false,
	13:  false,
	14:  false,
	15:  false,
	16:  false,
	17:  false,
	18:  false,
	19:  false,
	20:  false,
	21:  false,
	22:  false,
	23:  false,
	24:  false,
	25:  false,
	26:  false,
	27:  false,
	28:  false,
	29:  false,
	30:  false,
	31:  false,
	32:  true,
	33:  true,
	34:  false,
	35:  true,
	36:  true,
	37:  true,
	38:  true,
	39:  true,
	40:  true,
	41:  true,
	42:  true,
	43:  true,
	44:  true,
	45:  true,
	46:  true,
	47:  true,
	48:  true,
	49:  true,
	50:  true,
	51:  true,
	52:  true,
	53:  true,
	54:  true,
	55:  true,
	56:  true,
	57:  true,
	58:  true,
	59:  true,
	60:  true,
	61:  true,
	62:  true,
	63:  true,
	64:  true,
	65:  true,
	66:  true,
	67:  true,
	68:  true,
	69:  true,
	70:  true,
	71:  true,
	72:  true,
	73:  true,
	74:  true,
	75:  true,
	76:  true,
	77:  true,
	78:  true,
	79:  true,
	80:  true,
	81:  true,
	82:  true,
	83:  true,
	84:  true,
	85:  true,
	86:  true,
	87:  true,
	88:  true,
	89:  true,
	90:  true,
	91:  true,
	92:  false,
	93:  true,
	94:  true,
	95:  true,
	96:  true,
	97:  true,
	98:  true,
	99:  true,
	100: true,
	101: true,
	102: true,
	103: true,
	104: true,
	105: true,
	106: true,
	107: true,
	108: true,
	109: true,
	110: true,
	111: true,
	112: true,
	113: true,
	114: true,
	115: true,
	116: true,
	117: true,
	118: true,
	119: true,
	120: true,
	121: true,
	122: true,
	123: true,
	124: true,
	125: true,
	126: true,
}

// AppendString encodes the input string to json and appends
// the encoded string to the input byte slice.
//
// The operation loops though each byte in the string looking
// for characters that need json or utf8 encoding. If the string
// does not need encoding, then the string is appended in it's
// entirety to the byte slice.
// If we encounter a byte that does need encoding, switch up
// the operation and perform a byte-by-byte read-encode-append.
func appendString(dst []byte, s string) []byte {
	// Start with a double quote.
	dst = append(dst, '"')
	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			dst = appendStringComplex(dst, s, i)
			return append(dst, '"')
		}
	}
	// The string has no need for encoding an therefore is directly
	// appended to the byte slice.
	dst = append(dst, s...)
	// End with a double quote
	return append(dst, '"')
}

// appendStringComplex is used by appendString to take over an in
// progress JSON string encoding that encountered a character that needs
// to be encoded.
func appendStringComplex(dst []byte, s string, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a remplacement character code
				// in place of the invalid sequence.
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}

// appendBytesComplex is a mirror of the appendStringComplex
// with []byte arg
func appendBytesComplex(dst, s []byte, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRune(s[i:])
			if r == utf8.RuneError && size == 1 {
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}
