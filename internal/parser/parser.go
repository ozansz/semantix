package parser

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	sxQLLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z_\d]+`},
		{Name: `AnchorIdent`, Pattern: `_[a-zA-Z_\d]+`},
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Number`, Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{Name: `whitespace`, Pattern: `\s+`},
		{Name: `Punct`, Pattern: `[][=]`},
		{Name: `comment`, Pattern: `[#;][^\n]*`},
	})
)

type Parser struct {
	parser *participle.Parser[Expression]
}

func New() *Parser {
	sxQLParser := participle.MustBuild[Expression](
		participle.Lexer(sxQLLexer),
		participle.Unquote("String"),
		participle.Union[Object](SubjectObject{}, StringObject{}, NumberObject{}, RelationAnchorObject{}),
	)
	return &Parser{parser: sxQLParser}
}

func (p *Parser) Parse(input string) (*Expression, error) {
	return p.parser.ParseString("", input)
}

func (p *Parser) ParseFile(path string) (*Expression, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return p.parser.Parse("", f)
}
