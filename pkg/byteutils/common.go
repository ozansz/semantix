package byteutils

func TrimNull(b []byte) []byte {
	bb := make([]byte, 0, len(b))
	for _, c := range b {
		if c != 0 {
			bb = append(bb, c)
		}
	}
	return bb
}

func TrimNullString(b []byte) string {
	return string(TrimNull(b))
}
