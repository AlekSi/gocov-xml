package gocov_xml

import (
	"encoding/xml"
	"sort"
)

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
	Classes    Classes `xml:"classes>class"`
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

type Classes []Class

func (c Classes) Len() int           { return len(c) }
func (c Classes) Less(i, j int) bool { return c[i].Filename < c[j].Filename }
func (c Classes) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// check interfaces
var (
	_ sort.Interface = Classes{}
)
