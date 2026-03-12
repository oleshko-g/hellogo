package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	src := `/* Тестовый пакет */
package main

// Double умножает значение на 2.
func Double(i int) int {
    return i*2
}

func main() {
   // умножаем в цикле
   for i := 1; i < 5; i++ {
      fmt.Println(Double(i))
   }
}`
	fileSet := token.NewFileSet()

	parsedFile, err := parser.ParseFile(fileSet, "", src, parser.ParseComments)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, cc := range parsedFile.Comments {
		for _, c := range cc.List {
			fmt.Println(fileSet.Position(c.Slash), c.Text)
		}
	}
}
