package fact

import (
	"bytes"
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestEncodeDecodeRow(t *testing.T) {
	tests := []struct {
		desc                     string
		row                      *Row
		beforeEncodeCustomChecks func(t *testing.T, row *Row)
		afterDecodeCustomChecks  func(t *testing.T, row *Row)
	}{
		{
			desc: "empty row",
			row:  &Row{},
		},
		{
			desc: "basic row, all values min string",
			row: &Row{
				ID:        useULID(),
				Meta:      NewRowMeta(true, RowTypeBasic, RowSubjectTypeMinString, RowObjectTypeMinString),
				Subject:   NewSubjectFromMinString("Ozan"),
				Predicate: PredicateFromMinString("likes"),
				Object:    NewObjectFromMinString("Pizza"),
			},
			beforeEncodeCustomChecks: checkForOzanLikesPizza,
			afterDecodeCustomChecks:  checkForOzanLikesPizza,
		},
		{
			desc: "basic row, min string and uint64",
			row: &Row{
				ID:        useULID(),
				Meta:      NewRowMeta(true, RowTypeBasic, RowSubjectTypeMinString, RowObjectTypeUInt64),
				Subject:   NewSubjectFromMinString("Ozan"),
				Predicate: PredicateFromMinString("age"),
				Object:    NewObjectFromUInt64(24),
			},
			beforeEncodeCustomChecks: checkForOzanAge24,
			afterDecodeCustomChecks:  checkForOzanAge24,
		},
		{
			desc: "basic row, min string and uint64",
			row: &Row{
				ID:        useULID(),
				Meta:      NewRowMeta(true, RowTypeBasic, RowSubjectTypeMinString, RowObjectTypeFloat64),
				Subject:   NewSubjectFromMinString("John"),
				Predicate: PredicateFromMinString("height"),
				Object:    NewObjectFromFloat64(1.83),
			},
			beforeEncodeCustomChecks: checkForJohnHeight183,
			afterDecodeCustomChecks:  checkForJohnHeight183,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			if err := tc.row.Encode(&buf); err != nil {
				t.Fatalf("Failed to encode row: %v", err)
			}
			if tc.beforeEncodeCustomChecks != nil {
				tc.beforeEncodeCustomChecks(t, tc.row)
			}
			encoded := buf.Bytes()
			if len(encoded) != 128 {
				t.Fatalf("Encoded row is not 128 bytes long: %d", len(encoded))
			}
			decoded, err := NewRow(bytes.NewReader(encoded))
			if err != nil {
				t.Fatalf("Failed to decode row: %v", err)
			}
			if tc.afterDecodeCustomChecks != nil {
				tc.afterDecodeCustomChecks(t, decoded)
			}
			if diff := decoded.Diff(tc.row); diff != "" {
				t.Fatalf("Decoded row is different from original: %s", diff)
			}
		})
	}
}

func zeros128() []byte {
	var b [128]byte
	return b[:]
}

func useULID() ulid.ULID {
	return ulid.MustParse("01H4TQSE4S3QB93S8Z899B9AWS")
}

func useULIDBytes() []byte {
	u := useULID()
	return u[:]
}

func checkForOzanLikesPizza(t *testing.T, row *Row) {
	if row.ID != useULID() {
		t.Fatalf("Decoded row has wrong ID: %s", row.ID)
	}
	if !row.Meta.Active() {
		t.Fatalf("Row is not marked as active")
	}
	if row.Meta.Type() != RowTypeBasic {
		t.Fatalf("Row is not marked as basic, got %d", row.Meta.Type())
	}
	if row.Meta.SubjectType() != RowSubjectTypeMinString {
		t.Fatalf("Row is not marked as minstring subject, got %d", row.Meta.SubjectType())
	}
	if row.Meta.ObjectType() != RowObjectTypeMinString {
		t.Fatalf("Row is not marked as minstring object, got %d", row.Meta.ObjectType())
	}
	if row.Subject.MinString() != "Ozan" {
		t.Fatalf("Row has wrong subject: %q, expected: %q", row.Subject.MinString(), "Ozan")
	}
	if row.Predicate.MinString() != "likes" {
		t.Fatalf("Row has wrong predicate: %q, expected: %q", row.Predicate.MinString(), "likes")
	}
	if row.Object.MinString() != "Pizza" {
		t.Fatalf("Row has wrong object: %q, expected: %q", row.Object.MinString(), "Pizza")
	}
}

func checkForOzanAge24(t *testing.T, row *Row) {
	if row.ID != useULID() {
		t.Fatalf("Decoded row has wrong ID: %s", row.ID)
	}
	if !row.Meta.Active() {
		t.Fatalf("Row is not marked as active")
	}
	if row.Meta.Type() != RowTypeBasic {
		t.Fatalf("Row is not marked as basic, got %d", row.Meta.Type())
	}
	if row.Meta.SubjectType() != RowSubjectTypeMinString {
		t.Fatalf("Row is not marked as minstring subject, got %d", row.Meta.SubjectType())
	}
	if row.Meta.ObjectType() != RowObjectTypeUInt64 {
		t.Fatalf("Row is not marked as uint64 object, got %d", row.Meta.ObjectType())
	}
	if row.Subject.MinString() != "Ozan" {
		t.Fatalf("Row has wrong subject: %q, expected: %q", row.Subject.MinString(), "Ozan")
	}
	if row.Predicate.MinString() != "age" {
		t.Fatalf("Row has wrong predicate: %q, expected: %q", row.Predicate.MinString(), "age")
	}
	if row.Object.UInt64() != 24 {
		t.Fatalf("Row has wrong object: %d, expected: %d", row.Object.UInt64(), 24)
	}
}

func checkForJohnHeight183(t *testing.T, row *Row) {
	if row.ID != useULID() {
		t.Fatalf("Decoded row has wrong ID: %s", row.ID)
	}
	if !row.Meta.Active() {
		t.Fatalf("Row is not marked as active")
	}
	if row.Meta.Type() != RowTypeBasic {
		t.Fatalf("Row is not marked as basic, got %d", row.Meta.Type())
	}
	if row.Meta.SubjectType() != RowSubjectTypeMinString {
		t.Fatalf("Row is not marked as minstring subject, got %d", row.Meta.SubjectType())
	}
	if row.Meta.ObjectType() != RowObjectTypeFloat64 {
		t.Fatalf("Row is not marked as float64 object, got %d", row.Meta.ObjectType())
	}
	if row.Subject.MinString() != "John" {
		t.Fatalf("Row has wrong subject: %q, expected: %q", row.Subject.MinString(), "John")
	}
	if row.Predicate.MinString() != "height" {
		t.Fatalf("Row has wrong predicate: %q, expected: %q", row.Predicate.MinString(), "height")
	}
	if row.Object.Float64() != 1.83 {
		t.Fatalf("Row has wrong object: %f, expected: %f", row.Object.Float64(), 1.83)
	}
}
