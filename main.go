package main

import (
	"flag"
	"fmt"
	"os"

	"pg_flame/pkg/flame"
	"pg_flame/pkg/html"
	"pg_flame/pkg/plan"
)

var (
	version  = "1.0"
	hFlag    = flag.Bool("h", false, "print help info")
	helpFlag = flag.Bool("help", false, "print help info")
)

func main() {
	flag.Parse()

	if *hFlag || *helpFlag {
		printHelp()
	}

	err, p := plan.New(os.Stdin)
	if err != nil {
		handleErr(err)
	}

	f := flame.New(p)

	err = html.Generate(os.Stdout, f)
	if err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	msg := fmt.Errorf("Error: %v", err)
	fmt.Println(msg)
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
