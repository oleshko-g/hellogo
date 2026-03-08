package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var src = `package main

import "fmt"

func main()  {
	fmt.Println("Hello, world!")
}

`

func main() {

	fileSet := token.NewFileSet()

	parsedFile, err := parser.ParseFile(fileSet, "", src, parser.SkipObjectResolution)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(parsedFile, func(n ast.Node) bool {
		if n != nil {
			if v, ok := n.(*ast.CallExpr); ok {
				for _, arg := range v.Args {
					fmt.Println(arg)
				}
			}
			return true
		}
		return false
	})
}
