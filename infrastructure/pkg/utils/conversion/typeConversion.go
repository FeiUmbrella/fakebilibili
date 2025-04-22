package conversion

func Bool2Int8(b bool) int8 {
	if b == true {
		return int8(1)
	}
	return int8(0)
}

func Int82Bool(b int8) bool {
	return b > 0
}
