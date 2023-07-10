package fact

import (
	"github.com/oklog/ulid/v2"
	"github.com/ozansz/semantix/pkg/byteutils"
)

type Subject [16]byte

func (s Subject) MinString() string {
	return byteutils.TrimNullString(s[:])
}

func (s Subject) IDRef() ulid.ULID {
	return ulid.ULID(s)
}

func (s Subject) StringRef() ulid.ULID {
	return ulid.ULID(s)
}

func NewSubjectFromMinString(s string) Subject {
	var subject Subject
	copy(subject[:], s)
	return subject

}

func NewSubjectFromIDRef(id ulid.ULID) Subject {
	return Subject(id)
}

func NewSubjectFromStringRef(id ulid.ULID) Subject {
	return Subject(id)
}
