package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/internal/store"
	"github.com/ozansz/semantix/pkg/ptrutils"
)

const (
	defaultPrompt = "sxQL> "
)

type Interpreter struct {
	parser *parser.Parser
	prompt string
	store  store.Store
	quit   chan struct{}
}

// New returns a new interpreter.
func New(parser *parser.Parser, store store.Store) *Interpreter {
	return &Interpreter{
		parser: parser,
		store:  store,
		prompt: defaultPrompt,
		quit:   make(chan struct{}),
	}
}

func (i *Interpreter) Quit() {
	close(i.quit)
}

// ExecuteBatch executes the given batch of expressions.
func (i *Interpreter) ExecuteBatch(exprs []*parser.Expression) {
	for _, expr := range exprs {
		i.Execute(expr)
	}
}

// Execute executes the given expression.
func (i *Interpreter) Execute(expr *parser.Expression) {
	if expr.Query != nil {
		if err := i.executeQuery(expr.Query); err != nil {
			fmt.Printf("!! %v\n", err)
		}
	} else if expr.Fact != nil {
		if err := i.executeFact(expr.Fact); err != nil {
			fmt.Printf("!! %v\n", err)
		}
	}
}

// ExecuteREPL executes the interpreter in REPL mode.
func (i *Interpreter) ExecuteREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-i.quit:
			return
		default:
		}
		fmt.Print(i.prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		expr, err := i.parser.ParseLine(line)
		if err != nil {
			fmt.Printf("!!! Error parsing line: %v\n", err)
			continue
		}
		i.Execute(expr)
	}
}

func (i *Interpreter) executeFact(f *parser.Fact) error {
	t := &store.Triple{
		Subject:   f.Subject,
		Predicate: f.Predicate,
	}
	switch f.Object.(type) {
	case parser.SubjectObject:
		t.Object.Kind = store.ObjectKindSubject
		t.Object.StringValue = ptrutils.Ptr(f.Object.(parser.SubjectObject).Value)
	case parser.StringObject:
		t.Object.Kind = store.ObjectKindString
		t.Object.StringValue = ptrutils.Ptr(f.Object.(parser.StringObject).Value)
	case parser.RelationAnchorObject:
		t.Object.Kind = store.ObjectKindAnchor
		t.Object.StringValue = ptrutils.Ptr(f.Object.(parser.RelationAnchorObject).ID)
	case parser.NumberObject:
		t.Object.Kind = store.ObjectKindFloat
		t.Object.FloatValue = ptrutils.Ptr(f.Object.(parser.NumberObject).Value)
	}
	return i.store.Add(t)
}

func (i *Interpreter) executeQuery(q *parser.Query) error {
	if q.ObjectQuery != nil {
		return fmt.Errorf("compound queries are not supported yet")
	}
	if q.LinkedQuery != nil {
		return fmt.Errorf("linked queries are not supported yet")
	}

	qq := store.QueryFromAST(q)

	fmt.Printf("DEBUG: Query: %s\n", qq.Pretty())

	triples, err := i.store.Get(qq)
	if err != nil {
		return fmt.Errorf("failed to run query %+v: %v", qq, err)
	}

	for id, t := range triples {
		fmt.Printf("%010d: (%s, %s, %s)\n", id, t.Subject, t.Predicate, t.Object.String())
	}

	return nil
}
