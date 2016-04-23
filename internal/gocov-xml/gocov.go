package gocov_xml

import (
	"go/token"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/axw/gocov"
)

// ConvertGocov takes a slice of gocov packages and converts them to a slice go gocov-xml ones.
func ConvertGocov(packages []gocov.Package) []Package {
	fset := token.NewFileSet()
	tokenFiles := make(map[string]*token.File)

	// convert packages
	res := make([]Package, len(packages))
	for i, gPackage := range packages {
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
				class = &Class{Name: className, Filename: gFunction.File, Methods: []Method{}, Lines: []Line{}}
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

		// sort "classes" by filename
		sort.Sort(p.Classes)
		res[i] = p
	}

	return res
}
