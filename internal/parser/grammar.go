package parser

import (
	"fmt"
	"strings"

	"github.com/ozansz/semantix/pkg/ptrutils"
)

const (
	prettyExprIndent = 10
)

type File struct {
	Expressions []*Expression `@@*`
}

type Expression struct {
	Fact  *Fact  `  @@`
	Query *Query `| @@`
}

type Fact struct {
	Subject     *string `"(" ( @Ident`
	SubjectFact *Fact   `    | @@ )`
	Predicate   string  `"," @Ident`
	Object      Object  `"," ( @@`
	ObjectFact  *Fact   `    | @@ ) ")"`
}

type QueryKind int

const (
	QueryKindSimple QueryKind = iota
	QueryKindCompound
	QueryKindLinked
	QueryKindLinkedCompound
)

type Query struct {
	Subject      *string `"(" ( @Ident`
	SubjectVar   *string `    | @QueryIdent`
	SubjectQuery *Query  `    | @@ )`
	Predicate    *string `"," ( @Ident`
	PredicateVar *string `    | @QueryIdent )`
	Object       Object  `"," ( @@`
	ObjectVar    *string `    | @QueryIdent`
	ObjectQuery  *Query  `    | @@ ) ")"`
	LinkedQuery  *Query  `[ "-" ">" @@ ]`
	IDInFile     string
	Kind         QueryKind
}

type ObjectKind int

const (
	ObjectKindSubject ObjectKind = iota
	ObjectKindString
	ObjectKindNumber
)

type Object interface {
	String() string
	IsSubject() bool
	IsNumber() bool
	Copy() Object
	InnerValue() any
	Kind() ObjectKind
}

type SubjectObject struct {
	Value string `@Ident`
}

type StringObject struct {
	Value string `@String`
}

type NumberObject struct {
	Value float64 `@Number`
}

func (s SubjectObject) String() string { return s.Value }
func (s StringObject) String() string  { return fmt.Sprintf("%q", s.Value) }
func (n NumberObject) String() string  { return fmt.Sprintf("%f", n.Value) }

func (s SubjectObject) IsSubject() bool { return true }
func (s StringObject) IsSubject() bool  { return false }
func (n NumberObject) IsSubject() bool  { return false }

func (s SubjectObject) IsNumber() bool { return false }
func (s StringObject) IsNumber() bool  { return false }
func (n NumberObject) IsNumber() bool  { return true }

func (s SubjectObject) Copy() Object { return SubjectObject{Value: s.Value} }
func (s StringObject) Copy() Object  { return StringObject{Value: s.Value} }
func (n NumberObject) Copy() Object  { return NumberObject{Value: n.Value} }

func (s SubjectObject) Kind() ObjectKind { return ObjectKindSubject }
func (s StringObject) Kind() ObjectKind  { return ObjectKindString }
func (n NumberObject) Kind() ObjectKind  { return ObjectKindNumber }

func (s SubjectObject) InnerValue() any { return s.Value }
func (s StringObject) InnerValue() any  { return s.Value }
func (n NumberObject) InnerValue() any  { return n.Value }

func (f *File) Pretty() string {
	var sb strings.Builder
	for _, exp := range f.Expressions {
		sb.WriteString(exp.Pretty())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *Expression) Pretty() string {
	space := strings.Repeat(" ", prettyExprIndent)
	var sb strings.Builder
	if e.Fact != nil {
		sb.WriteString(space)
		sb.WriteString(e.Fact.Pretty())
	} else if e.Query != nil {
		sb.WriteString(e.Query.IDInFile)
		sb.WriteString(": ")
		sb.WriteString(space[:len(space)-len(e.Query.IDInFile)])
		sb.WriteString(e.Query.Pretty())
	}
	return sb.String()
}

func (s *Query) Pretty() string {
	var sb strings.Builder
	sb.WriteRune('(')

	if s.Subject != nil {
		sb.WriteString(*s.Subject)
	} else if s.SubjectVar != nil {
		sb.WriteString(*s.SubjectVar)
	}
	sb.WriteString(", ")
	if s.Predicate != nil {
		sb.WriteString(*s.Predicate)
	} else if s.PredicateVar != nil {
		sb.WriteString(*s.PredicateVar)
	}
	sb.WriteString(", ")
	if s.Object != nil {
		sb.WriteString(s.Object.String())
	} else if s.ObjectVar != nil {
		sb.WriteString(*s.ObjectVar)
	} else if s.ObjectQuery != nil {
		sb.WriteString(s.ObjectQuery.Pretty())
	}

	sb.WriteRune(')')

	if s.LinkedQuery != nil {
		sb.WriteString(" -> ")
		sb.WriteString(s.LinkedQuery.Pretty())
	}

	return sb.String()
}

func (f *Fact) Pretty() string {
	var sb strings.Builder
	sb.WriteRune('(')

	if f.Subject != nil {
		sb.WriteString(*f.Subject)
	} else if f.SubjectFact != nil {
		sb.WriteString(f.SubjectFact.Pretty())
	}
	sb.WriteString(", ")
	sb.WriteString(f.Predicate)
	sb.WriteString(", ")
	if f.Object != nil {
		sb.WriteString(f.Object.String())
	} else if f.ObjectFact != nil {
		sb.WriteString(f.ObjectFact.Pretty())
	}

	sb.WriteRune(')')

	return sb.String()
}

func (q *Query) IsLinkedCompound() bool {
	if q.LinkedQuery == nil {
		return false
	}
	if q.ObjectQuery != nil {
		return true
	}
	currQ := q
	for {
		if currQ == nil {
			break
		}
		if currQ.ObjectQuery != nil {
			return true
		}
		currQ = currQ.LinkedQuery
	}
	return false
}

func (f *Fact) Copy() *Fact {
	newF := &Fact{
		Subject:   ptrutils.PtrFromPtr(f.Subject),
		Predicate: f.Predicate,
	}
	if f.SubjectFact != nil {
		newF.SubjectFact = f.SubjectFact.Copy()
	}
	if f.ObjectFact != nil {
		newF.ObjectFact = f.ObjectFact.Copy()
	}
	if f.Object != nil {
		newF.Object = f.Object.Copy()
	}
	return newF
}
