package parser

import (
	"fmt"
	"strings"
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
	Subject   string `@Ident`
	Predicate string `@Ident`
	Object    Object `@@`
	Anchor    string `[ "as" @AnchorIdent ]`
}

type QueryKind int

const (
	QueryKindSimple QueryKind = iota
	QueryKindCompound
	QueryKindLinked
	QueryKindLinkedCompound
)

type Query struct {
	Subject      *string `( @Ident`
	SubjectVar   *string `| @QueryIdent )`
	Predicate    *string `( @Ident`
	PredicateVar *string `| @QueryIdent )`
	Object       Object  `( @@`
	ObjectVar    *string `| @QueryIdent`
	ObjectQuery  *Query  `| "(" @@ ")" )`
	LinkedQuery  *Query  `[ "-" ">" @@ ]`
	IDInFile     string
	Kind         QueryKind
}

type Object interface {
	String() string
	IsSubject() bool
	IsNumber() bool
	Copy() Object
	InnerValue() any
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

type RelationAnchorObject struct {
	ID string `@AnchorIdent`
}

func (s SubjectObject) String() string        { return s.Value }
func (s StringObject) String() string         { return fmt.Sprintf("%q", s.Value) }
func (n NumberObject) String() string         { return fmt.Sprintf("%f", n.Value) }
func (r RelationAnchorObject) String() string { return fmt.Sprintf("@%s", r.ID) }

func (s SubjectObject) IsSubject() bool        { return true }
func (s StringObject) IsSubject() bool         { return false }
func (n NumberObject) IsSubject() bool         { return false }
func (r RelationAnchorObject) IsSubject() bool { return false }

func (s SubjectObject) IsNumber() bool        { return false }
func (s StringObject) IsNumber() bool         { return false }
func (n NumberObject) IsNumber() bool         { return true }
func (r RelationAnchorObject) IsNumber() bool { return false }

func (s SubjectObject) Copy() Object        { return SubjectObject{Value: s.Value} }
func (s StringObject) Copy() Object         { return StringObject{Value: s.Value} }
func (n NumberObject) Copy() Object         { return NumberObject{Value: n.Value} }
func (r RelationAnchorObject) Copy() Object { return RelationAnchorObject{ID: r.ID} }

func (s SubjectObject) InnerValue() any        { return s.Value }
func (s StringObject) InnerValue() any         { return s.Value }
func (n NumberObject) InnerValue() any         { return n.Value }
func (r RelationAnchorObject) InnerValue() any { return r.ID }

func (f *File) Pretty() string {
	var sb strings.Builder
	for _, exp := range f.Expressions {
		sb.WriteString(exp.Pretty())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *Expression) Pretty() string {
	var sb strings.Builder
	anchorSpace := strings.Repeat(" ", prettyExprIndent)
	if e.Fact != nil {
		if e.Fact.Anchor != "" {
			sb.WriteString(e.Fact.Anchor)
			sb.WriteString(": ")
			if (len(e.Fact.Anchor) + 2) > prettyExprIndent {
				sb.WriteString("\n")
				sb.WriteString(anchorSpace)
			} else {
				sb.WriteString(anchorSpace[:prettyExprIndent-len(e.Fact.Anchor)-2])
			}
		} else {
			sb.WriteString(anchorSpace)
		}
		sb.WriteRune('(')
		sb.WriteString(e.Fact.Subject)
		sb.WriteString(", ")
		sb.WriteString(e.Fact.Predicate)
		sb.WriteString(", ")
		sb.WriteString(e.Fact.Object.String())
		sb.WriteRune(')')
	} else if e.Query != nil {
		sb.WriteString(e.Query.IDInFile)
		sb.WriteString(": ")
		if (len(e.Query.IDInFile) + 2) > prettyExprIndent {
			sb.WriteString("\n")
			sb.WriteString(anchorSpace)
		} else {
			sb.WriteString(anchorSpace[:prettyExprIndent-len(e.Query.IDInFile)-2])
		}
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
