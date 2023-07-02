package store

import (
	"fmt"
	"log"

	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/pkg/ptrutils"
)

type Triple struct {
	Subject   string
	Predicate string
	Object    Object
}

type Object struct {
	StringValue *string
	// IntegerValue *int64
	FloatValue *float64
	Kind       ObjectKind
}

type ObjectKind int

const (
	ObjectKindSubject ObjectKind = iota
	ObjectKindAnchor
	ObjectKindString
	// ObjectKindInteger
	ObjectKindFloat
)

type Query struct {
	SubjectFilter      *string
	PredicateFilter    *string
	ObjectFilterString *string
	// ObjectFilterInteger *int64
	ObjectFilterFloat *float64

	// TODO: Add support for these.
	// ObjectAnchor    *Triple
}

type Store interface {
	Add(*Triple) error
	Get(*Query) (map[uint32]*Triple, error)

	Close() error
}

func (t *Triple) Copy() *Triple {
	return &Triple{
		Subject:   t.Subject,
		Predicate: t.Predicate,
		Object: Object{
			StringValue: ptrutils.PtrFromPtr(t.Object.StringValue),
			// IntegerValue: ptrutils.PtrFromPtr(t.Object.IntegerValue),
			FloatValue: ptrutils.PtrFromPtr(t.Object.FloatValue),
			Kind:       t.Object.Kind,
		},
	}
}

func (o *Object) String() string {
	switch o.Kind {
	case ObjectKindSubject:
		return *o.StringValue
	case ObjectKindAnchor:
		return *o.StringValue
	case ObjectKindString:
		return fmt.Sprintf("%q", *o.StringValue)
	// case ObjectKindInteger:
	// 	return fmt.Sprintf("%d", *o.IntegerValue)
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
	return qq
}

func (q *Query) Pretty() string {
	s, p, o := "?", "?", "?"
	if q.SubjectFilter != nil {
		s = *q.SubjectFilter
	}
	if q.PredicateFilter != nil {
		p = *q.PredicateFilter
	}
	if q.ObjectFilterFloat != nil {
		o = fmt.Sprintf("%f", *q.ObjectFilterFloat)
	}
	if q.ObjectFilterString != nil {
		o = fmt.Sprintf("%q", *q.ObjectFilterString)
	}
	return fmt.Sprintf("<%s, %s, %s>", s, p, o)
}
