package fact

// NOTE: The second byte is reserved for future use
type RowMeta [2]byte

func (r *RowMeta) Active() bool {
	return (r[0]&0b10000000)>>7 == 1
}

func (r *RowMeta) Type() RowType {
	return RowType((r[0] & 0b01100000) >> 5)
}

func (r *RowMeta) SubjectType() RowSubjectType {
	return RowSubjectType((r[0] & 0b00011000) >> 3)
}

func (r *RowMeta) ObjectType() RowObjectType {
	return RowObjectType(r[0] & 0b00000111)
}

func (r *RowMeta) SetActive(active bool) {
	if active {
		r[0] |= 0b10000000
	} else {
		r[0] &= 0b01111111
	}
}

func (r *RowMeta) SetType(t RowType) {
	r[0] &= 0b10011111
	r[0] |= byte(t) << 5
}

func (r *RowMeta) SetSubjectType(t RowSubjectType) {
	r[0] &= 0b11100111
	r[0] |= byte(t) << 3
}

func (r *RowMeta) SetObjectType(t RowObjectType) {
	r[0] &= 0b11111000
	r[0] |= byte(t)
}

func NewRowMeta(active bool, t RowType, st RowSubjectType, ot RowObjectType) RowMeta {
	r := RowMeta{}
	r.SetActive(active)
	r.SetType(t)
	r.SetSubjectType(st)
	r.SetObjectType(ot)
	return r
}
