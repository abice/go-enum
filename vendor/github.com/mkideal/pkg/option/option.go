package option

func Bool(dft bool, bools ...bool) bool {
	if len(bools) == 0 {
		return dft
	}
	return bools[0]
}

func Byte(dft byte, bytes ...byte) byte {
	if len(bytes) == 0 {
		return dft
	}
	return bytes[0]
}

func Rune(dft rune, runes ...rune) rune {
	if len(runes) == 0 {
		return dft
	}
	return runes[0]
}

func String(dft string, strings ...string) string {
	if len(strings) == 0 {
		return dft
	}
	return strings[0]
}

func Int(dft int, ints ...int) int {
	if len(ints) == 0 {
		return dft
	}
	return ints[0]
}

func Int8(dft int8, int8s ...int8) int8 {
	if len(int8s) == 0 {
		return dft
	}
	return int8s[0]
}

func Int16(dft int16, int16s ...int16) int16 {
	if len(int16s) == 0 {
		return dft
	}
	return int16s[0]
}

func Int32(dft int32, int32s ...int32) int32 {
	if len(int32s) == 0 {
		return dft
	}
	return int32s[0]
}

func Int64(dft int64, int64s ...int64) int64 {
	if len(int64s) == 0 {
		return dft
	}
	return int64s[0]
}

func Uint(dft uint, uints ...uint) uint {
	if len(uints) == 0 {
		return dft
	}
	return uints[0]
}

func Uint8(dft uint8, uint8s ...uint8) uint8 {
	if len(uint8s) == 0 {
		return dft
	}
	return uint8s[0]
}

func Uint16(dft uint16, uint16s ...uint16) uint16 {
	if len(uint16s) == 0 {
		return dft
	}
	return uint16s[0]
}

func Uint32(dft uint32, uint32s ...uint32) uint32 {
	if len(uint32s) == 0 {
		return dft
	}
	return uint32s[0]
}

func Uint64(dft uint64, uint64s ...uint64) uint64 {
	if len(uint64s) == 0 {
		return dft
	}
	return uint64s[0]
}

func Interface(dft interface{}, interfaces ...interface{}) interface{} {
	if len(interfaces) == 0 {
		return dft
	}
	return interfaces[0]
}
