// Package cli implements high level function for running the Parser from command-line.
package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/qba73/gov"
)

func Run() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [files...]\n", os.Args[0])
		fmt.Println("Parses Go dependencies from files or standard input and returns JSON output.\nInput is an output from the Go command: `go version -v -m <go-binary>`")
		flag.PrintDefaults()
	}
	flag.Parse()
	p, err := gov.NewParser(
		gov.WithInputFromArgs(flag.Args()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	s, err := p.ToJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s\n", s)
	os.Exit(0)
}
