package fact

type RowType uint8
type RowSubjectType uint8
type RowObjectType uint8

// NOTE: UPDATE RowMeta AFTER UPDATING THIS
// TYPE COUNT MUST BE <= 4
const (
	RowTypeBasic RowType = iota
	RowTypeSubjectRef
	RowTypeObjectRef
	RowTypeSubjectRefObjectRef
)

// NOTE: UPDATE RowMeta AFTER UPDATING THIS
// TYPE COUNT MUST BE <= 4
const (
	RowSubjectTypeMinString RowSubjectType = iota
	RowSubjectTypeIDRef
	RowSubjectTypeStringRef
)

// NOTE: UPDATE RowMeta AFTER UPDATING THIS
// TYPE COUNT MUST BE <= 8
const (
	RowObjectTypeMinString RowObjectType = iota
	RowObjectTypeSubjectMinString
	RowObjectTypeSubjectStringRef
	RowObjectTypeUInt64
	RowObjectTypeFloat64
	RowObjectTypeStringRef
	RowObjectTypeIDRef
)
