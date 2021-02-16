package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/javiercbk/gencmp"
)

func main() {
	var file string
	var out string
	var templateFile string
	flag.StringVar(&file, "file", "", "the file name to read the struct from")
	flag.StringVar(&out, "out", "", "the out file path to write the comparison function")
	flag.StringVar(&templateFile, "template-file", "", "the template file to use")
	flag.Parse()
	file = strings.TrimSpace(file)
	templateFile = strings.TrimSpace(templateFile)
	if file == "" {
		log.Fatal("no file given")
	}
	if templateFile == "" {
		log.Fatal("no template file given")
	}
	outFile, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalf("could not create out file: %s", err)
	}
	defer outFile.Close()
	tmplt, err := template.ParseFiles(templateFile)
	if err != nil || tmplt == nil {
		log.Fatalf("could not read template file: %s", err)
	}
	err = gencmp.Generate(file, tmplt, outFile)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
