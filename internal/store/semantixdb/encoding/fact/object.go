package fact

import (
	"encoding/binary"
	"math"

	"github.com/oklog/ulid/v2"
	"github.com/ozansz/semantix/pkg/byteutils"
)

type Object [64]byte

func (o Object) MinString() string {
	return byteutils.TrimNullString(o[:])
}

func (o Object) SubjectMinString() string {
	return byteutils.TrimNullString(o[:])
}

func (o Object) SubjectStringRef() ulid.ULID {
	var u ulid.ULID
	copy(u[:], o[:])
	return u
}

func (o Object) UInt64() uint64 {
	return binary.BigEndian.Uint64(o[:8])
}

func (o Object) Float64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(o[:8]))
}

func (o Object) IDRef() ulid.ULID {
	var u ulid.ULID
	copy(u[:], o[:])
	return u
}

func (o Object) StringRef() ulid.ULID {
	var u ulid.ULID
	copy(u[:], o[:])
	return u
}

func NewObjectFromMinString(s string) Object {
	var o Object
	copy(o[:], s)
	return o

}

func NewObjectFromSubjectMinString(subject string) Object {
	var o Object
	copy(o[:], subject)
	return o
}

func NewObjectFromSubjectStringRef(id ulid.ULID) Object {
	var o Object
	copy(o[:], id[:])
	return o
}

func NewObjectFromUInt64(val uint64) Object {
	var o Object
	binary.BigEndian.PutUint64(o[:8], val)
	return o
}

func NewObjectFromFloat64(val float64) Object {
	var o Object
	binary.BigEndian.PutUint64(o[:8], math.Float64bits(val))
	return o
}

func NewObjectFromIDRef(id ulid.ULID) Object {
	var o Object
	copy(o[:], id[:])
	return o
}

func NewObjectFromStringRef(id ulid.ULID) Object {
	var o Object
	copy(o[:], id[:])
	return o
}
