package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go/build"
	"go/token"
	"io"
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
	Sources    []Source  `xml:"sources>source"`
	Packages   []Package `xml:"packages>package"`
}

type Source struct {
	Path string `xml:",chardata"`
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
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.Writer) {
	var r struct{ Packages []gocov.Package }
	err := json.NewDecoder(in).Decode(&r)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	tokenFiles := make(map[string]*token.File)

	srcDirs := build.Default.SrcDirs()
	sources := make([]Source, len(srcDirs))
	for i, dir := range srcDirs {
		sources[i] = Source{dir}
	}

	// convert packages
	packages := make([]Package, len(r.Packages))
	for i, gPackage := range r.Packages {
		// group functions by filename and "class" (type)
		files := make(map[string]map[string]*Class)
		for _, gFunction := range gPackage.Functions {
			classes := files[gFunction.File]
			if classes == nil {
				// group functions by "class" (type) in a File
				classes = make(map[string]*Class)
				files[gFunction.File] = classes
			}

			s := strings.Split("-."+gFunction.Name, ".") // className is "-" for package-level functions
			className, methodName := s[len(s)-2], s[len(s)-1]
			class := classes[className]
			if class == nil {
				fileName := stripKnownSources(sources, gFunction.File)
				class = &Class{Name: className, Filename: fileName, Methods: []Method{}, Lines: []Line{}}
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
			lines := make([]Line, len(gFunction.Statements))
			for i, s := range gFunction.Statements {
				lineno := tokenFile.Line(tokenFile.Pos(s.Start))
				line := Line{Number: lineno, Hits: s.Reached}
				lines[i] = line
				class.Lines = append(class.Lines, line)
			}

			class.Methods = append(class.Methods, Method{Name: methodName, Lines: lines})
		}

		// fill package with "classes"
		p := Package{Name: gPackage.Name, Classes: []Class{}}
		for _, classes := range files {
			for _, class := range classes {
				p.Classes = append(p.Classes, *class)
			}
		}
		packages[i] = p
	}

	coverage := Coverage{Sources: sources, Packages: packages, Timestamp: time.Now().UnixNano() / int64(time.Millisecond)}

	fmt.Fprintf(out, xml.Header)
	fmt.Fprintf(out, "<!DOCTYPE coverage SYSTEM \"http://cobertura.sourceforge.net/xml/coverage-03.dtd\">\n")

	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	err = encoder.Encode(coverage)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(out)
}

func stripKnownSources(sources []Source, fileName string) string {
	for _, source := range sources {
		prefix := source.Path
		prefix = strings.TrimSuffix(prefix, string(os.PathSeparator)) + string(os.PathSeparator)
		if strings.HasPrefix(fileName, prefix) {
			return strings.TrimPrefix(fileName, prefix)
		}
	}
	return fileName
}
