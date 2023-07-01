package parser

import (
	"fmt"
	"strings"
)

type Expression struct {
	Facts []*Fact `@@*`
}

func (e *Expression) Pretty() string {
	var sb strings.Builder
	maxAnchorLen := 0
	for _, fact := range e.Facts {
		if len(fact.Anchor) > maxAnchorLen {
			maxAnchorLen = len(fact.Anchor)
		}
	}
	anchorSpace := strings.Repeat(" ", maxAnchorLen)
	for _, fact := range e.Facts {
		if fact.Anchor != "" {
			sb.WriteString(fact.Anchor)
			sb.WriteString(": ")
			sb.WriteString(anchorSpace[:maxAnchorLen-len(fact.Anchor)])
		} else {
			sb.WriteString(anchorSpace)
			sb.WriteString("  ")
		}
		sb.WriteRune('(')
		sb.WriteString(fact.Subject)
		sb.WriteString(", ")
		sb.WriteString(fact.Predicate)
		sb.WriteString(", ")
		sb.WriteString(fact.Object.String())
		sb.WriteRune(')')
		sb.WriteString("\n")
	}
	return sb.String()
}

type Fact struct {
	Subject   string `@Ident`
	Predicate string `@Ident`
	Object    Object `@@`
	Anchor    string `[ "as" @AnchorIdent ]`
}

type Object interface {
	String() string
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
