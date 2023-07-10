package byteutils

import (
	"fmt"
	"strings"
)

type ByteArray []byte

func (b ByteArray) Print() {
	fmt.Println(b.Sprint())
}

func (b ByteArray) Sprint() string {
	var sb strings.Builder
	row := 0
	for i, c := range b {
		if i%16 == 0 {
			sb.WriteString(fmt.Sprintf("0x%02x: ", row))
		}
		sb.WriteString(fmt.Sprintf("%02x ", c))
		if i%16 == 15 {
			row++
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
