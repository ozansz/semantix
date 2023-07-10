package fact

import "github.com/ozansz/semantix/pkg/byteutils"

type Predicate [16]byte

func (p Predicate) MinString() string {
	return byteutils.TrimNullString(p[:])
}

func PredicateFromMinString(s string) Predicate {
	var predicate Predicate
	copy(predicate[:], s)
	return predicate

}
