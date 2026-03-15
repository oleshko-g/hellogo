package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	_ "os"
)

func main() {
	src := `package main

func main() {
     ids := 77
     id := ids + 1
     fmt.Println("id равно:", id/2 )
}`

	fset := token.NewFileSet()

	parsedFile, _ := parser.ParseFile(fset, "", src, parser.SkipObjectResolution)

	ast.Inspect(parsedFile, func(n ast.Node) bool {
		if v, ok := n.(*ast.Ident); ok {
			if v.Name == "id" {
				v.Name = "Ident"
			}
		}
		return true
	})

	printer.Fprint(os.Stdout, fset, parsedFile)
}
