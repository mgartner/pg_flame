package main

import (
	"flag"
	"fmt"
	"os"

	"pg_flame/pkg/html"
	"pg_flame/pkg/plan"
)

var (
	// goreleaser automatically overrides this based on the tag
	version  = "dev"
	hFlag    = flag.Bool("h", false, "print help info")
	helpFlag = flag.Bool("help", false, "print help info")
)

func main() {
	flag.Parse()

	if *hFlag || *helpFlag {
		printHelp()
	}

	p, err := plan.New(os.Stdin)
	if err != nil {
		handleErr(err)
	}

	err = html.Generate(os.Stdout, p)
	if err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v", err)
	os.Exit(1)
}

func printHelp() {
	help := `pg_flame %s

Turn Postgres query plans into flamegraphs.

Usage:

  pg_flame [options]

Without Options:

  Reads a JSON query plan from standard input and writes the
  flamegraph html to standard output.

Options:

  -h, --help	print help information
`

	fmt.Printf(help, version)
	os.Exit(0)
}
