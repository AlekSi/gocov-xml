package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/axw/gocov"
)

type Coverage struct {
	XMLName    xml.Name  `xml:"coverage"`
	LineRate   float32   `xml:"line-rate,attr"`
	BranchRate float32   `xml:"branch-rate,attr"`
	Version    string    `xml:"version,attr"`
	Timestamp  int64     `xml:"timestamp,attr"`
	Packages   []Package `xml:"packages>package"`
}

type Package struct {
	Name       string  `xml:"name,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`
	Classes    []Class `xml:"classes>class"`
}

type Class struct {
	Name       string   `xml:"name,attr"`
	Filename   string   `xml:"filename,attr"`
	LineRate   float32  `xml:"line-rate,attr"`
	BranchRate float32  `xml:"branch-rate,attr"`
	Complexity float32  `xml:"complexity,attr"`
	Methods    []Method `xml:"methods>method"`
	Lines      []Line   `xml:"lines>line"`
}

type Method struct {
	Name       string  `xml:"name,attr"`
	Signature  string  `xml:"signature,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Lines      []Line  `xml:"lines>line"`
}

type Line struct {
	Number int   `xml:"number,attr"`
	Hits   int64 `xml:"hits,attr"`
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var r struct{ Packages []gocov.Package }
	err = json.Unmarshal(data, &r)
	if err != nil {
		panic(err)
	}

	// convert packages
	packages := make([]Package, len(r.Packages))
	for i, gPackage := range r.Packages {
		// group functions by "class" (type)
		classes := make(map[string]*Class)
		for _, gFunction := range gPackage.Functions {
			s := strings.Split("-."+gFunction.Name, ".")
			className, methodName := s[len(s)-2], s[len(s)-1]
			class := classes[className]
			if class == nil {
				// type's methods can be in many files, we just pick first
				class = &Class{Name: className, Filename: gFunction.File, Methods: []Method{}}
				classes[className] = class
			}

			// FIXME convert statements to lines
			lines := make([]Line, len(gFunction.Statements))
			for i, s := range gFunction.Statements {
				lines[i] = Line{Number: s.Start, Hits: s.Reached}
			}

			class.Methods = append(class.Methods, Method{Name: methodName, Lines: lines})
		}

		// fill package with "classes"
		p := Package{Name: gPackage.Name, Classes: make([]Class, len(classes))}
		j := 0
		for _, class := range classes {
			p.Classes[j] = *class
			j++
		}
		packages[i] = p
	}

	coverage := Coverage{Packages: packages, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)}

	data, err = xml.MarshalIndent(coverage, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf(xml.Header)
	fmt.Printf("<!DOCTYPE coverage SYSTEM \"http://cobertura.sourceforge.net/xml/coverage-03.dtd\">")
	fmt.Printf("%s\n", data)
}
