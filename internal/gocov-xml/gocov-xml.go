package gocov_xml

import (
	"encoding/xml"
	"sort"
)

// Coverage representes coverage information for Go packages.
type Coverage struct {
	XMLName    xml.Name  `xml:"coverage"`
	LineRate   float32   `xml:"line-rate,attr"`
	BranchRate float32   `xml:"branch-rate,attr"`
	Version    string    `xml:"version,attr"`
	Timestamp  int64     `xml:"timestamp,attr"`
	Packages   []Package `xml:"packages>package"`
}

// Package is a single Go package.
type Package struct {
	Name       string  `xml:"name,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`
	Classes    Classes `xml:"classes>class"`
}

// Class is a single "class", represented by single Go type.
type Class struct {
	Name       string   `xml:"name,attr"`
	Filename   string   `xml:"filename,attr"`
	LineRate   float32  `xml:"line-rate,attr"`
	BranchRate float32  `xml:"branch-rate,attr"`
	Complexity float32  `xml:"complexity,attr"`
	Methods    []Method `xml:"methods>method"`
	Lines      []Line   `xml:"lines>line"`
}

// Classes is a slice of "classes", with methods helping to sort it by filename.
type Classes []Class

func (c Classes) Len() int           { return len(c) }
func (c Classes) Less(i, j int) bool { return c[i].Filename < c[j].Filename }
func (c Classes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// Method is a single Go method or function.
type Method struct {
	Name       string  `xml:"name,attr"`
	Signature  string  `xml:"signature,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Lines      []Line  `xml:"lines>line"`
}

// Line is a single line of Go code.
type Line struct {
	Number int   `xml:"number,attr"`
	Hits   int64 `xml:"hits,attr"`
}

// check interfaces
var (
	_ sort.Interface = Classes{}
)
