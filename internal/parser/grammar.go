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
	Fact        *Fact        `  @@`
	SimpleQuery *SimpleQuery `| @@`
	// Queries     []*Query
}

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
	} else if e.SimpleQuery != nil {
		sb.WriteString(e.SimpleQuery.IDInFile)
		sb.WriteString(": ")
		if (len(e.SimpleQuery.IDInFile) + 2) > prettyExprIndent {
			sb.WriteString("\n")
			sb.WriteString(anchorSpace)
		} else {
			sb.WriteString(anchorSpace[:prettyExprIndent-len(e.SimpleQuery.IDInFile)-2])
		}
		sb.WriteString(e.SimpleQuery.Pretty())
	}
	return sb.String()
}

// func (f *File) Pretty() string {
// 	var sb strings.Builder
// 	for _, exp := range f.Expressions {
// 		sb.WriteString(exp.Pretty())
// 		sb.WriteString("\n")
// 	}
// 	return sb.String()
// }

// func (e *Expression) Pretty() string {
// 	var sb strings.Builder
// 	maxAnchorLen := 0

// 	for _, fact := range e.Facts {
// 		if len(fact.Anchor) > maxAnchorLen {
// 			maxAnchorLen = len(fact.Anchor)
// 		}
// 	}
// 	anchorSpace := strings.Repeat(" ", maxAnchorLen)
// 	for _, fact := range e.Facts {
// 		if fact.Anchor != "" {
// 			sb.WriteString(fact.Anchor)
// 			sb.WriteString(": ")
// 			sb.WriteString(anchorSpace[:maxAnchorLen-len(fact.Anchor)])
// 		} else {
// 			sb.WriteString(anchorSpace)
// 			sb.WriteString("  ")
// 		}
// 		sb.WriteRune('(')
// 		sb.WriteString(fact.Subject)
// 		sb.WriteString(", ")
// 		sb.WriteString(fact.Predicate)
// 		sb.WriteString(", ")
// 		sb.WriteString(fact.Object.String())
// 		sb.WriteRune(')')
// 		sb.WriteString("\n")
// 	}
// 	for _, query := range e.Queries {
// 		sb.WriteString("?:")
// 		sb.WriteString(anchorSpace[:maxAnchorLen])
// 		sb.WriteRune('(')
// 		if query.Subject != nil {
// 			sb.WriteString(*query.Subject)
// 		} else if query.SubjectVar != nil {
// 			sb.WriteString("?")
// 			sb.WriteString(*query.SubjectVar)
// 		} else if query.SubjectVarSkip != nil {
// 			sb.WriteString("!")
// 			sb.WriteString(*query.SubjectVarSkip)
// 		}
// 		sb.WriteString(", ")
// 		if query.Predicate != nil {
// 			sb.WriteString(*query.Predicate)
// 		} else if query.PredicateVar != nil {
// 			sb.WriteString("?")
// 			sb.WriteString(*query.PredicateVar)
// 		} else if query.PredicateVarSkip != nil {
// 			sb.WriteString("!")
// 			sb.WriteString(*query.PredicateVarSkip)
// 		}
// 		sb.WriteString(", ")
// 		if query.ObjectVar != nil {
// 			sb.WriteString("?")
// 			sb.WriteString(*query.ObjectVar)
// 		} else if query.ObjectVarSkip != nil {
// 			sb.WriteString("!")
// 			sb.WriteString(*query.ObjectVarSkip)
// 		} else {
// 			sb.WriteString(query.Object.String())
// 		}
// 		sb.WriteRune(')')
// 		sb.WriteString("\n")
// 	}
// 	return sb.String()
// }

// type Query struct {
// 	Subject          *string
// 	SubjectVar       *string
// 	SubjectVarSkip   *string
// 	Predicate        *string
// 	PredicateVar     *string
// 	PredicateVarSkip *string
// 	Object           Object
// 	ObjectVar        *string
// 	ObjectVarSkip    *string
// }

type Fact struct {
	Subject   string `@Ident`
	Predicate string `@Ident`
	Object    Object `@@`
	Anchor    string `[ "as" @AnchorIdent ]`
}

// func (f *Fact) ToQuery() (*Query, error) {
// 	if f.Anchor != "" {
// 		return nil, fmt.Errorf("cannot convert fact with anchor to query")
// 	}
// 	query := &Query{}
// 	if strings.HasPrefix(f.Subject, "?") {
// 		query.SubjectVar = ptrutils.StringPtr(f.Subject[1:])
// 	} else if strings.HasPrefix(f.Subject, "!") {
// 		query.SubjectVarSkip = ptrutils.StringPtr(f.Subject[1:])
// 	} else {
// 		query.Subject = ptrutils.StringPtr(f.Subject)
// 	}
// 	if strings.HasPrefix(f.Predicate, "?") {
// 		query.PredicateVar = ptrutils.StringPtr(f.Predicate[1:])
// 	} else if strings.HasPrefix(f.Predicate, "!") {
// 		query.PredicateVarSkip = ptrutils.StringPtr(f.Predicate[1:])
// 	} else {
// 		query.Predicate = ptrutils.StringPtr(f.Predicate)
// 	}
// 	s := f.Object.String()
// 	if strings.HasPrefix(s, "?") {
// 		query.ObjectVar = ptrutils.StringPtr(s[1:])
// 	} else if strings.HasPrefix(f.Object.String(), "!") {
// 		query.ObjectVarSkip = ptrutils.StringPtr(s[1:])
// 	} else {
// 		query.Object = f.Object.Copy()
// 	}
// 	return query, nil
// }

type Object interface {
	String() string
	IsSubject() bool
	Copy() Object
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

func (s SubjectObject) Copy() Object        { return SubjectObject{Value: s.Value} }
func (s StringObject) Copy() Object         { return StringObject{Value: s.Value} }
func (n NumberObject) Copy() Object         { return NumberObject{Value: n.Value} }
func (r RelationAnchorObject) Copy() Object { return RelationAnchorObject{ID: r.ID} }

type SimpleQuery struct {
	Subject      *string `( @Ident`
	SubjectVar   *string `| @QueryIdent )`
	Predicate    *string `( @Ident`
	PredicateVar *string `| @QueryIdent )`
	Object       Object  `( @@`
	ObjectVar    *string `| @QueryIdent )`
	IDInFile     string
}

func (s *SimpleQuery) Pretty() string {
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
	}

	sb.WriteRune(')')
	return sb.String()
}
