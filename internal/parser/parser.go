package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	sxQLLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `(?:#|--)[^\n]*\n?`},
		{Name: `QueryIdent`, Pattern: `[?!][a-zA-Z][a-zA-Z_\d]*`},
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z_\d]*`},
		{Name: `AnchorIdent`, Pattern: `_[a-zA-Z_\d]+`},
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Number`, Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{Name: `Punct`, Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{Name: `Whitespace`, Pattern: `\s+`},
	})
)

type Parser struct {
	expParser  *participle.Parser[Expression]
	fileParser *participle.Parser[File]
}

func New() *Parser {
	return &Parser{
		expParser: participle.MustBuild[Expression](
			participle.Lexer(sxQLLexer),
			participle.Unquote("String"),
			participle.Union[Object](SubjectObject{}, StringObject{}, NumberObject{}, RelationAnchorObject{}),
			participle.Elide("Comment", "Whitespace"),
			participle.UseLookahead(2),
		),
		fileParser: participle.MustBuild[File](
			participle.Lexer(sxQLLexer),
			participle.Unquote("String"),
			participle.Union[Object](SubjectObject{}, StringObject{}, NumberObject{}, RelationAnchorObject{}),
			participle.Elide("Comment", "Whitespace"),
			participle.UseLookahead(2),
		),
	}
}

func (p *Parser) ParseLine(input string) (*Expression, error) {
	exp, err := p.expParser.ParseString("<LINE>", input)
	if err != nil {
		return nil, err
	}
	// if err := p.separateQueries(exp); err != nil {
	// 	return nil, err
	// }
	return exp, nil
}

func (p *Parser) ParseFile(path string) (*File, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	file, err := p.fileParser.Parse(path, f)
	if err != nil {
		return nil, err
	}
	// if err := p.separateQueries(exp); err != nil {
	// 	return nil, err
	// }
	i := 1
	for _, e := range file.Expressions {
		if e.SimpleQuery != nil {
			e.SimpleQuery.IDInFile = fmt.Sprintf("Q%d", i)
			i++
		}
	}
	return file, nil
}

func (p *Parser) Ebnf() string {
	var sb strings.Builder
	sb.WriteString("Expression:\n")
	sb.WriteString(p.expParser.String())
	sb.WriteString("\n\nFile:\n")
	sb.WriteString(p.fileParser.String())
	return sb.String()
}

// func (p *Parser) separateQueries(exp *Expression) error {
// 	queryIndices := []int{}
// 	exp.Queries = []*Query{}
// 	for i, fact := range exp.Facts {
// 		if strings.HasPrefix(fact.Subject, "?") || strings.HasPrefix(fact.Subject, "!") ||
// 			strings.HasPrefix(fact.Predicate, "?") || strings.HasPrefix(fact.Predicate, "!") ||
// 			fact.Object.IsSubject() && (strings.HasPrefix(fact.Object.String(), "?") || strings.HasPrefix(fact.Object.String(), "!")) {
// 			queryIndices = append(queryIndices, i)

// 			q, err := fact.ToQuery()
// 			if err != nil {
// 				return err
// 			}
// 			exp.Queries = append(exp.Queries, q)
// 		}
// 	}
// 	for i := len(queryIndices) - 1; i >= 0; i-- {
// 		exp.Facts = append(exp.Facts[:queryIndices[i]], exp.Facts[queryIndices[i]+1:]...)
// 	}
// 	return nil
// }
