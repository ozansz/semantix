package ptrutils

func Ptr[T string | int8 | int16 | int32 |
	int64 | uint8 | uint16 | uint32 |
	uint64 | float32 | float64](v T) *T {
	return &v
}

func PtrFromPtr[T string | int8 | int16 | int32 |
	int64 | uint8 | uint16 | uint32 |
	uint64 | float32 | float64](vp *T) *T {
	if vp == nil {
		return nil
	}
	return Ptr(*vp)
}
