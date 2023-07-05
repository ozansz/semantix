package store

import (
	"fmt"
	"log"
	"strings"

	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/pkg/ptrutils"
)

type Object struct {
	StringValue *string
	FloatValue  *float64
	Kind        ObjectKind
}

type ObjectKind int

const (
	ObjectKindSubject ObjectKind = iota
	ObjectKindString
	ObjectKindFloat
)

type Query struct {
	SubjectFilter      *string
	SubjectFilterQuery *Query
	PredicateFilter    *string
	ObjectFilterString *string
	ObjectFilterFloat  *float64
	ObjectFilterQuery  *Query
	LinkedQuery        *Query
}

type Store interface {
	Add(*parser.Fact) error
	Get(*Query) (map[uint32]*parser.Fact, error)

	Sync() error

	Close() error
}

func (o *Object) String() string {
	switch o.Kind {
	case ObjectKindSubject:
		return *o.StringValue
	case ObjectKindString:
		return fmt.Sprintf("%q", *o.StringValue)
	case ObjectKindFloat:
		return fmt.Sprintf("%f", *o.FloatValue)
	}
	log.Panicf("Unreachable, Object has an unexpected kind: %v", o.Kind)
	return ""
}

func QueryFromAST(q *parser.Query) *Query {
	qq := &Query{}
	if q.Subject != nil {
		qq.SubjectFilter = ptrutils.PtrFromPtr(q.Subject)
	}
	if q.SubjectQuery != nil {
		qq.SubjectFilterQuery = QueryFromAST(q.SubjectQuery)
	}
	if q.Predicate != nil {
		qq.PredicateFilter = ptrutils.PtrFromPtr(q.Predicate)
	}
	if q.Object != nil {
		if q.Object.IsNumber() {
			qq.ObjectFilterFloat = ptrutils.Ptr(q.Object.InnerValue().(float64))
		} else {
			qq.ObjectFilterString = ptrutils.Ptr(q.Object.InnerValue().(string))
		}
	}
	if q.ObjectQuery != nil {
		qq.ObjectFilterQuery = QueryFromAST(q.ObjectQuery)
	}
	return qq
}

func (q *Query) Matches(t *parser.Fact) bool {
	if q.LinkedQuery != nil {
		// TODO: Implement linked queries
		log.Panicf("Linked queries are not supported yet")
	}
	if q.SubjectFilter != nil && ((t.Subject != nil && *t.Subject != *q.SubjectFilter) || (t.Subject == nil)) {
		return false
	}
	if q.SubjectFilterQuery != nil && (t.SubjectFact != nil && !q.SubjectFilterQuery.Matches(t.SubjectFact) || t.SubjectFact == nil) {
		return false
	}
	if q.PredicateFilter != nil && t.Predicate != *q.PredicateFilter {
		return false
	}
	if q.ObjectFilterString != nil && (t.Object != nil && (!t.Object.IsNumber() && t.Object.InnerValue().(string) != *q.ObjectFilterString || t.Object.IsNumber()) || t.Object == nil) {
		return false
	}
	if q.ObjectFilterFloat != nil && (t.Object != nil && (t.Object.IsNumber() && t.Object.InnerValue().(float64) != *q.ObjectFilterFloat || !t.Object.IsNumber()) || t.Object == nil) {
		return false
	}
	if q.ObjectFilterQuery != nil && (t.ObjectFact != nil && !q.ObjectFilterQuery.Matches(t.ObjectFact) || t.ObjectFact == nil) {
		return false
	}
	return true
}

func (q *Query) Pretty() string {
	var sb strings.Builder
	sb.WriteRune('(')
	if q.SubjectFilter != nil {
		sb.WriteString(*q.SubjectFilter)
	} else if q.SubjectFilterQuery != nil {
		sb.WriteString(q.SubjectFilterQuery.Pretty())
	} else {
		sb.WriteString("*")
	}
	sb.WriteString(", ")
	if q.PredicateFilter != nil {
		sb.WriteString(*q.PredicateFilter)
	} else {
		sb.WriteString("*")
	}
	sb.WriteString(", ")
	if q.ObjectFilterString != nil {
		sb.WriteString(*q.ObjectFilterString)
	} else if q.ObjectFilterFloat != nil {
		sb.WriteString(fmt.Sprintf("%f", *q.ObjectFilterFloat))
	} else if q.ObjectFilterQuery != nil {
		sb.WriteString(q.ObjectFilterQuery.Pretty())
	} else {
		sb.WriteString("*")
	}
	sb.WriteRune(')')
	return sb.String()
}
