package gencmp

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"regexp"
	"strings"
	"text/template"
)

const (
	EmptyInterface   = "interface{}"
	goTypeBool       = "bool"
	goTypeString     = "string"
	goTypeInt        = "int"
	goTypeInt8       = "int8"
	goTypeInt16      = "int16"
	goTypeInt32      = "int32"
	goTypeInt64      = "int64"
	goTypeUint       = "uint"
	goTypeUint8      = "uint8"
	goTypeUint16     = "uint16"
	goTypeUint32     = "uint32"
	goTypeUint64     = "uint64"
	goTypeUintptr    = "uintptr"
	goTypeByte       = "byte"
	goTypeRune       = "rune"
	goTypeFloat32    = "float32"
	goTypeFloat64    = "float64"
	goTypeComplex64  = "complex64"
	goTypeComplex128 = "complex128"
	goTypeTime       = "time.Time"
)

var (
	jsonNameRegexp = regexp.MustCompile(`json:"(.+?)(,.*)?"`)
	lineRegexp     = regexp.MustCompile(`^[\s\t]*$`)
)

type field struct {
	Name          string
	VarName       string
	JSONName      string
	FieldType     string
	StructName    string
	IsStructSlice bool
}

type str struct {
	Name   string
	Fields []field
}

// Generate generates the compare function for a struct in the given file and
// writes it in the out writer
func Generate(file string, tmplt *template.Template, out io.Writer) error {
	astFile, err := parseASTFromFile(file)
	if err != nil {
		return err
	}
	structs := readStructs(astFile)
	var b []byte
	buf := bytes.NewBuffer(b)
	err = tmplt.Execute(buf, structs)
	if err != nil {
		return err
	}
	var outBytes []byte
	outBuf := bytes.NewBuffer(outBytes)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		if !lineRegexp.MatchString(line) {
			outBuf.WriteString(line)
			outBuf.WriteString("\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	_, err = io.Copy(out, outBuf)
	return err
}

func readStructs(astFile *ast.File) []str {
	structs := make([]str, 0)
	for _, d := range astFile.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if ok {
			for _, s := range genDecl.Specs {
				t, ok := s.(*ast.TypeSpec)
				if ok {
					st, ok := t.Type.(*ast.StructType)
					if ok {
						readStruct := str{
							Name: t.Name.Name,
						}
						fields := make([]field, len(st.Fields.List))
						for i, f := range st.Fields.List {
							name := f.Names[0].Name
							varName := strings.ToLower(string(name[0])) + name[1:]
							fi := field{
								Name:    f.Names[0].Name,
								VarName: varName,
							}
							if f.Tag != nil {
								matches := jsonNameRegexp.FindAllString(f.Tag.Value, -1)
								if len(matches) > 2 {
									fi.JSONName = matches[1]
								}
							}
							var b strings.Builder
							isGoType, isSlice, _ := typeToString(f.Type, &b)
							typeName := b.String()
							if !isGoType {
								fi.StructName = typeName
								fi.IsStructSlice = isSlice
							} else {
								fi.FieldType = typeName
							}
							fields[i] = fi
						}
						readStruct.Fields = fields
						structs = append(structs, readStruct)
					}
				}
			}
		}
	}
	return structs
}

func typeToString(expr ast.Expr, b *strings.Builder) (bool, bool, bool) {
	switch typed := expr.(type) {
	case *ast.Ident:
		b.WriteString(typed.Name)
		return isGoType(typed.Name), false, false
	case *ast.ArrayType:
		b.WriteString("[]")
		isStruct, _, _ := typeToString(typed.Elt, b)
		return isStruct, true, false
	case *ast.MapType:
		i, ok := typed.Key.(*ast.Ident)
		if ok {
			b.WriteString(fmt.Sprintf("map[%s]", i.Name))
			isStruct, _, _ := typeToString(typed.Value, b)
			return isStruct, false, true
		}
	}
	return false, false, false
}

func isGoType(t string) bool {
	switch t {
	case EmptyInterface:
		return true
	case goTypeTime:
		return true
	case goTypeBool:
		return true
	case goTypeString:
		return true
	case goTypeInt:
		return true
	case goTypeInt8:
		return true
	case goTypeInt16:
		return true
	case goTypeInt32:
		return true
	case goTypeInt64:
		return true
	case goTypeUint:
		return true
	case goTypeUint8:
		return true
	case goTypeUint16:
		return true
	case goTypeUint32:
		return true
	case goTypeUint64:
		return true
	case goTypeUintptr:
		return true
	case goTypeByte:
		return true
	case goTypeRune:
		return true
	case goTypeFloat32:
		return true
	case goTypeFloat64:
		return true
	case goTypeComplex64:
		return true
	case goTypeComplex128:
		return true
	default:
		return false
	}
}
