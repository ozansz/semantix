package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/internal/store"
)

const (
	defaultPrompt = "sxQL> "
)

var (
	commands = map[string]func(*Interpreter){
		"quit": func(i *Interpreter) {
			i.Quit()
		},
		"exit": func(i *Interpreter) {
			i.Quit()
		},
		"fsync": func(i *Interpreter) {
			if err := i.store.Sync(); err != nil {
				fmt.Printf("!! Error syncing the store: %v\n", err)
			}
		},
	}
)

type Interpreter struct {
	parser *parser.Parser
	prompt string
	store  store.Store
	quit   chan struct{}
	debug  bool
}

type InterpreterOption func(*Interpreter)

func WithDebug() func(*Interpreter) {
	return func(i *Interpreter) {
		i.debug = true
	}
}

// New returns a new interpreter.
func New(parser *parser.Parser, store store.Store, opts ...InterpreterOption) *Interpreter {
	i := &Interpreter{
		parser: parser,
		store:  store,
		prompt: defaultPrompt,
		quit:   make(chan struct{}),
	}
	for _, o := range opts {
		o(i)
	}
	return i
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

		if handler, ok := commands[line]; ok {
			handler(i)
			continue
		}

		expr, err := i.parser.ParseLine(line)
		if err != nil {
			fmt.Printf("!!! Error parsing line: %v\n", err)
			continue
		}
		i.Execute(expr)
	}
}

func (i *Interpreter) executeFact(f *parser.Fact) error {
	return i.store.Add(f)
}

func (i *Interpreter) executeQuery(q *parser.Query) error {
	qq := store.QueryFromAST(q)

	if i.debug {
		fmt.Printf("Executing query: %s\n", qq.Pretty())
	}

	results, err := i.store.Get(qq)
	if err != nil {
		return err
	}
	fmt.Println()
	for id, result := range results {
		fmt.Printf("%010d: %s\n", id, result.Pretty())
	}
	fmt.Println()
	return nil
}
