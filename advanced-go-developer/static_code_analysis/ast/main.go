package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

var src = `package main

import "fmt"

func main()  {
	fmt.Println("С 8 Марта!")
}`

func main() {

	fileSet := token.NewFileSet()

	parsedFile, err := parser.ParseFile(fileSet, "", src,
		parser.SkipObjectResolution)

	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(parsedFile, func(n ast.Node) bool {
		if n != nil {
			switch v := n.(type) {
			case *ast.CallExpr:
				for _, arg := range v.Args {
					printer.Fprint(os.Stdout, fileSet, arg)
				}
			}
			return true
		}

		return false
	})
}
