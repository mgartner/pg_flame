package main

import (
	"fmt"
	"os"
	"pg_flame/pkg/flame"
	"pg_flame/pkg/html"
	"pg_flame/pkg/plan"
)

func main() {
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
