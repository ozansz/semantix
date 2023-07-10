package fact

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/oklog/ulid/v2"
)

type Row struct {
	ID        ulid.ULID
	Meta      RowMeta
	Reserved  [14]byte
	Subject   Subject
	Predicate Predicate
	Object    Object
}

func NewRow(r io.Reader) (*Row, error) {
	f := &Row{}
	if err := binary.Read(r, binary.BigEndian, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Row) Encode(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, r)
}

func (r *Row) Diff(other *Row) string {
	var sb strings.Builder
	if bytes.Compare(r.ID[:], other.ID[:]) != 0 {
		sb.WriteString("+ ID: ")
		sb.WriteString(fmt.Sprintf("%v", r.ID))
		sb.WriteString("\n- ID: ")
		sb.WriteString(fmt.Sprintf("%v\n", r.ID))
	}
	if r.Meta.Active() != other.Meta.Active() {
		sb.WriteString("+ Active: ")
		sb.WriteString(fmt.Sprintf("%t", r.Meta.Active()))
		sb.WriteString("\n- Active: ")
		sb.WriteString(fmt.Sprintf("%t\n", other.Meta.Active()))
	}
	if r.Meta.Type() != other.Meta.Type() {
		sb.WriteString("+ Type: ")
		sb.WriteString(fmt.Sprintf("%v", r.Meta.Type()))
		sb.WriteString("\n- Type: ")
		sb.WriteString(fmt.Sprintf("%v\n", other.Meta.Type()))
	}
	if r.Meta.SubjectType() != other.Meta.SubjectType() {
		sb.WriteString("+ SubjectType: ")
		sb.WriteString(fmt.Sprintf("%v", r.Meta.SubjectType()))
		sb.WriteString("\n- SubjectType: ")
		sb.WriteString(fmt.Sprintf("%v\n", other.Meta.SubjectType()))
	}
	if r.Meta.ObjectType() != other.Meta.ObjectType() {
		sb.WriteString("+ ObjectType: ")
		sb.WriteString(fmt.Sprintf("%v", r.Meta.ObjectType()))
		sb.WriteString("\n- ObjectType: ")
		sb.WriteString(fmt.Sprintf("%v\n", other.Meta.ObjectType()))
	}
	switch r.Meta.SubjectType() {
	case RowSubjectTypeMinString:
		if r.Subject.MinString() != other.Subject.MinString() {
			sb.WriteString("+ Subject: ")
			sb.WriteString(fmt.Sprintf("%q", r.Subject.MinString()))
			sb.WriteString("\n- Subject: ")
			sb.WriteString(fmt.Sprintf("%q\n", other.Subject.MinString()))
		}
	case RowSubjectTypeStringRef:
		ref := r.Subject.StringRef()
		otherRef := other.Subject.StringRef()
		if bytes.Compare(ref[:], otherRef[:]) != 0 {
			sb.WriteString("+ Subject: ")
			sb.WriteString(fmt.Sprintf("%v", r.Subject.StringRef()))
			sb.WriteString("\n- Subject: ")
			sb.WriteString(fmt.Sprintf("%v\n", other.Subject.StringRef()))
		}
	case RowSubjectTypeIDRef:
		ref := r.Subject.IDRef()
		otherRef := other.Subject.IDRef()
		if bytes.Compare(ref[:], otherRef[:]) != 0 {
			sb.WriteString("+ Subject: ")
			sb.WriteString(fmt.Sprintf("%v", r.Subject.IDRef()))
			sb.WriteString("\n- Subject: ")
			sb.WriteString(fmt.Sprintf("%v\n", other.Subject.IDRef()))
		}
	}
	if r.Predicate.MinString() != other.Predicate.MinString() {
		sb.WriteString("+ Predicate: ")
		sb.WriteString(fmt.Sprintf("%q", r.Predicate.MinString()))
		sb.WriteString("\n- Predicate: ")
		sb.WriteString(fmt.Sprintf("%q\n", other.Predicate.MinString()))
	}
	switch r.Meta.ObjectType() {
	case RowObjectTypeMinString:
		if r.Object.MinString() != other.Object.MinString() {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%q", r.Object.MinString()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%q\n", other.Object.MinString()))
		}
	case RowObjectTypeStringRef:
		ref := r.Object.StringRef()
		otherRef := other.Object.StringRef()
		if bytes.Compare(ref[:], otherRef[:]) != 0 {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%v", r.Object.StringRef()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%v\n", other.Object.StringRef()))
		}
	case RowObjectTypeIDRef:
		ref := r.Object.IDRef()
		otherRef := other.Object.IDRef()
		if bytes.Compare(ref[:], otherRef[:]) != 0 {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%v", r.Object.IDRef()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%v\n", other.Object.IDRef()))
		}
	case RowObjectTypeUInt64:
		if r.Object.UInt64() != other.Object.UInt64() {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%d", r.Object.UInt64()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%d\n", other.Object.UInt64()))
		}
	case RowObjectTypeFloat64:
		if r.Object.Float64() != other.Object.Float64() {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%f", r.Object.Float64()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%f\n", other.Object.Float64()))
		}
	case RowObjectTypeSubjectMinString:
		if r.Object.SubjectMinString() != other.Object.SubjectMinString() {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%q", r.Object.SubjectMinString()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%q\n", other.Object.SubjectMinString()))
		}
	case RowObjectTypeSubjectStringRef:
		ref := r.Object.SubjectStringRef()
		otherRef := other.Object.SubjectStringRef()
		if bytes.Compare(ref[:], otherRef[:]) != 0 {
			sb.WriteString("+ Object: ")
			sb.WriteString(fmt.Sprintf("%v", r.Object.SubjectStringRef()))
			sb.WriteString("\n- Object: ")
			sb.WriteString(fmt.Sprintf("%v\n", other.Object.SubjectStringRef()))
		}
	}
	return sb.String()
}
