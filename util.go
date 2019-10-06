package logl

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

// func quoRem(x, y int) (quo, rem int) {
// 	quo = x / y
// 	rem = x - quo*y
// 	return
// }
