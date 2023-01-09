package logl

import (
	"bytes"
	"unicode/utf8"
)

func lastByteIs(s string, b byte) bool {
	if n := len(s); n > 0 {
		return s[n-1] == b
	}
	return false
}

func quoRem(x, y int) (quo, rem int) {
	quo = x / y
	rem = x % y
	return
}

// ------------------------------------------------------------------------------
// Nibbles
func byteToNibbles(b byte) (hi, lo byte) {
	hi = b >> 4
	lo = b & 0xF
	return
}

func nibblesToByte(hi, lo byte) (b byte) {
	b |= hi << 4
	b |= lo & 0xF
	return
}

// ------------------------------------------------------------------------------
var (
	hexTableLower = []byte("0123456789abcdef")
	hexTableUpper = []byte("0123456789ABCDEF")

	hexTable = hexTableLower
)

func byteHex(x byte) string {
	hi, lo := byteToNibbles(x)
	bs := []byte{
		hexTable[hi],
		hexTable[lo],
	}
	return string(bs)
}

func writeByteHex(b *bytes.Buffer, x byte) {
	hi, lo := byteToNibbles(x)
	b.WriteByte(hexTable[hi])
	b.WriteByte(hexTable[lo])
}

// ------------------------------------------------------------------------------
var safeSet = makeSafeSet()

func makeSafeSet() (set [utf8.RuneSelf]bool) {
	for i := range set {
		x := byte(i)
		// ASCII control characters (0-31)
		if (x > 31) && (x != '\\') && (x != '"') {
			set[i] = true
		}
	}
	return
}
