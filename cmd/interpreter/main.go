package main

import (
	"flag"
	"fmt"
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
	fmt.Printf("EBNF:\n%s", parser.Ebnf())
	file, err := parser.ParseFile(*sxQLFile)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}
	log.Printf("Parsed expressions:\n\n%s", file.Pretty())
}
