package parser

import (
	"fmt"
	"os"

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
		{Name: `Whitespace`, Pattern: `[ \t\n\r]+`},
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
	exprs := []*Expression{exp}
	p.postProcessQueries(exprs)
	return exprs[0], nil
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
	p.postProcessQueries(file.Expressions)
	return file, nil
}

func (p *Parser) Ebnf() string {
	return p.fileParser.String()
}

func (p *Parser) postProcessQueries(exprs []*Expression) {
	basicIndx, compoundIndx, linkedIndx, linkedCompIndx := 1, 1, 1, 1
	for _, e := range exprs {
		if e.Query != nil {
			if e.Query.IsLinkedCompound() {
				e.Query.IDInFile = fmt.Sprintf("LCQ%d", linkedCompIndx)
				e.Query.Kind = QueryKindLinkedCompound
				linkedCompIndx++
			} else if e.Query.ObjectQuery != nil {
				e.Query.IDInFile = fmt.Sprintf("CQ%d", compoundIndx)
				e.Query.Kind = QueryKindCompound
				compoundIndx++
			} else if e.Query.LinkedQuery != nil {
				e.Query.IDInFile = fmt.Sprintf("LQ%d", linkedIndx)
				e.Query.Kind = QueryKindLinked
				linkedIndx++
			} else {
				e.Query.IDInFile = fmt.Sprintf("Q%d", basicIndx)
				e.Query.Kind = QueryKindSimple
				basicIndx++
			}
		}
	}
}
