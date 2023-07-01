package main

import (
	"flag"
	"log"

	"github.com/ozansz/semantix/internal/parser"
)

var (
	sxQLFile = flag.String("f", "", "Path to sxQL file")
)

func main() {
	flag.Parse()
	if *sxQLFile == "" {
		flag.Usage()
		return
	}
	parser := parser.New()
	expr, err := parser.ParseFile(*sxQLFile)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}
	log.Printf("Parsed expression:\n\n%s", expr.Pretty())
}
