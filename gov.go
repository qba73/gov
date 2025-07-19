// Package gov provides functionality for parsing Go dependencies in JSON format.
package gov

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Dependency represents Go dependency module (package)
// listed in the output of the go version command:
//
//	go version -v -m <binary>
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Digest  string `json:"digest"`
}

// Parse takes a line from output from go version -v -m <binary>
// parses it and get package name, version and digest and returns
// Dependency struct populated with the data.
// It returns error if the line can't be parsed due to not expected syntax.
func Parse(line string) (Dependency, error) {
	var pkgname, version, digest string
	n, err := fmt.Sscanf(line, "dep\t%s\t%s\t%s", &pkgname, &version, &digest)
	if err != nil || n != 3 {
		return Dependency{}, fmt.Errorf("invalid line format, want 3 elements, got: %d, %w", n, err)
	}
	return Dependency{
		Name:    pkgname,
		Version: version,
		Digest:  digest,
	}, nil
}

// option is a functional option for configuring parser.
type option func(*parser) error

// WithInput configures Parser's input source.
func WithInput(input io.Reader) option {
	return func(p *parser) error {
		if input == nil {
			return errors.New("nil input reader")
		}
		p.input = input
		return nil
	}
}

// WithInputFromArgs configures Parser input as files listed as arguments.
// It allows Parser to read multiple files and return one stream with
// combined Go packages (dependencies).
//
// Note that parser does not filter duplicates.
func WithInputFromArgs(args []string) option {
	return func(p *parser) error {
		if len(args) < 1 {
			return nil
		}
		p.inputFiles = make([]io.Reader, len(args))
		for i, path := range args {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			p.inputFiles[i] = f
		}
		p.input = io.MultiReader(p.inputFiles...)
		return nil
	}
}

// WithOutput counfigures Parser output.
func WithOutput(output io.Writer) option {
	return func(p *parser) error {
		if output == nil {
			return errors.New("nil output writer")
		}
		p.output = output
		return nil
	}
}

type parser struct {
	input      io.Reader
	output     io.Writer
	inputFiles []io.Reader
}

// NewParser configures and returns Parser.
// It returns a default parser if no options are supplied.
// Default parser takes input from os.Stdin.
// It returns error if supplied options don't configure
// valid inputs and outputs.
func NewParser(opts ...option) (*parser, error) {
	p := &parser{
		input:  os.Stdin,
		output: os.Stdout,
	}
	for _, o := range opts {
		err := o(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// Dependencies parses input and returns a slice of Dependencies.
// Dependencies represents a list of parsed Go packages included
// in the output of go version -v -m <binary> command.
func (p *parser) Dependencies() ([]Dependency, error) {
	defer func(files []io.Reader) {
		for _, f := range files {
			f.(io.Closer).Close()
		}
	}(p.inputFiles)

	var dependencies []Dependency
	input := bufio.NewScanner(p.input)
	for input.Scan() {
		l := strings.TrimSpace(input.Text())
		if !strings.HasPrefix(l, "dep") {
			continue
		}
		dep, err := Parse(l)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, dep)
	}
	if err := input.Err(); err != nil {
		return nil, err
	}
	return dependencies, nil
}

// ToJSON returns JSON representation of the parsed dependencies.
func (p *parser) ToJSON() (string, error) {
	d, err := p.Dependencies()
	if err != nil {
		return "", err
	}
	if len(d) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ParseDependencies returns JSON with parsed Go dependencies.
// It uses default Parser configured to read from os.Stdin.
func ParseDependencies() (string, error) {
	p, err := NewParser()
	if err != nil {
		return "", err
	}
	return p.ToJSON()
}
