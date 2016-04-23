package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/axw/gocov"

	"github.com/AlekSi/gocov-xml/internal/gocov-xml"
)

func main() {
	flag.Parse()

	var r struct{ Packages []gocov.Package }
	err := json.NewDecoder(os.Stdin).Decode(&r)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	tokenFiles := make(map[string]*token.File)

	// convert packages
	packages := make([]gocov_xml.Package, len(r.Packages))
	for i, gPackage := range r.Packages {
		// group functions by filename and "class" (type)
		files := make(map[string]map[string]*gocov_xml.Class)
		for _, gFunction := range gPackage.Functions {
			classes := files[gFunction.File]
			if classes == nil {
				// group functions by "class" (type) in a File
				classes = make(map[string]*gocov_xml.Class)
				files[gFunction.File] = classes
			}

			s := strings.Split("-."+gFunction.Name, ".") // className is "-" for package-level functions
			className, methodName := s[len(s)-2], s[len(s)-1]
			class := classes[className]
			if class == nil {
				class = &gocov_xml.Class{Name: className, Filename: gFunction.File, Methods: []gocov_xml.Method{}, Lines: []gocov_xml.Line{}}
				classes[className] = class
			}

			// from github.com/axw/gocov /gocov/annotate.go#printFunctionSource
			// Load the file for line information. Probably overkill, maybe
			// just compute the lines from offsets in here.
			setContent := false
			tokenFile := tokenFiles[gFunction.File]
			if tokenFile == nil {
				info, err := os.Stat(gFunction.File)
				if err != nil {
					panic(err)
				}
				tokenFile = fset.AddFile(gFunction.File, fset.Base(), int(info.Size()))
				setContent = true
			}

			tokenData, err := ioutil.ReadFile(gFunction.File)
			if err != nil {
				panic(err)
			}
			if setContent {
				// This processes the content and records line number info.
				tokenFile.SetLinesForContent(tokenData)
			}

			// convert statements to lines
			lines := make([]gocov_xml.Line, len(gFunction.Statements))
			for i, s := range gFunction.Statements {
				lineno := tokenFile.Line(tokenFile.Pos(s.Start))
				line := gocov_xml.Line{Number: lineno, Hits: s.Reached}
				lines[i] = line
				class.Lines = append(class.Lines, line)
			}

			class.Methods = append(class.Methods, gocov_xml.Method{Name: methodName, Lines: lines})
		}

		// fill package with "classes"
		p := gocov_xml.Package{Name: gPackage.Name, Classes: []gocov_xml.Class{}}
		for _, classes := range files {
			for _, class := range classes {
				p.Classes = append(p.Classes, *class)
			}
		}
		packages[i] = p
	}

	coverage := gocov_xml.Coverage{Packages: packages, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)}

	fmt.Printf(xml.Header)
	fmt.Printf("<!DOCTYPE coverage SYSTEM \"http://cobertura.sourceforge.net/xml/coverage-03.dtd\">\n")

	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "\t")
	err = encoder.Encode(coverage)
	if err != nil {
		panic(err)
	}

	fmt.Println()
}
