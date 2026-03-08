package main

import (
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

	err = ast.Print(fileSet, parsedFile)
	if err != nil {
		log.Fatal(err)
	}

}
