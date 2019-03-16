package main

import (
	"fmt"
	"os"

	"github.com/kwo/exodus"
)

func main() {
	cwd := "statements"
	out := "statements.go"
	pkg := "schema"
	args := []string{"."}
	if err := exodus.Generate(pkg, out, cwd, args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
