package fact

import (
	"testing"
)

func TestRowMeta(t *testing.T) {
	tests := []struct {
		desc   string
		active bool
		t      RowType
		st     RowSubjectType
		ot     RowObjectType
		want   RowMeta
	}{
		{
			desc: "empty",
			want: RowMeta{},
		},
		{
			desc:   "active",
			active: true,
			want:   RowMeta{0b10000000},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			got := NewRowMeta(tc.active, tc.t, tc.st, tc.ot)
			if got != tc.want {
				t.Fatalf("Got: %v\nWant: %v", got, tc.want)
			}
			if got.Active() != tc.active {
				t.Fatalf("Active: Got: %t\nWant: %t", got.Active(), tc.active)
			}
			if got.Type() != tc.t {
				t.Fatalf("Type: Got: %d\nWant: %d", got.Type(), tc.t)
			}
			if got.SubjectType() != tc.st {
				t.Fatalf("SubjectType: Got: %d\nWant: %d", got.SubjectType(), tc.st)
			}
			if got.ObjectType() != tc.ot {
				t.Fatalf("ObjectType: Got: %d\nWant: %d", got.ObjectType(), tc.ot)
			}
		})
	}
}
