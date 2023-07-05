package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozansz/semantix/internal/interpreter"
	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/internal/store/filestore"
)

var (
	sxQLFile = flag.String("f", "", "Path to sxQL file")
	debug    = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()
	parser := parser.New()

	store, err := filestore.New(filestore.WithPersistentFile("store.db"), filestore.WithDebug())
	if err != nil {
		log.Fatalf("Error creating store: %v", err)
	}

	intOps := []interpreter.InterpreterOption{}
	if *debug {
		intOps = append(intOps, interpreter.WithDebug())
	}
	interpreter := interpreter.New(parser, store, intOps...)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("Received SIGTERM, exiting...")
		interpreter.Quit()
		store.Close()
		os.Exit(1)
	}()

	if *sxQLFile != "" {
		file, err := parser.ParseFile(*sxQLFile)
		if err != nil {
			log.Fatalf("Error parsing file: %v", err)
		}
		interpreter.ExecuteBatch(file.Expressions)
	}
	interpreter.ExecuteREPL()
}
