package gencmp

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

func parseASTFromFile(filePath string) (*ast.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v", err)
	}
	fset := token.NewFileSet() // positions are relative to fset
	parsed, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing ast %v", err)
	}
	return parsed, nil
}
