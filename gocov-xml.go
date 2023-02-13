package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlekSi/gocov-xml/convert"
)

// Coverage information
type Coverage struct {
	XMLName         xml.Name  `xml:"coverage"`
	LineRate        float32   `xml:"line-rate,attr"`
	BranchRate      float32   `xml:"branch-rate,attr"`
	LinesCovered    float32   `xml:"lines-covered,attr"`
	LinesValid      int64     `xml:"lines-valid,attr"`
	BranchesCovered int64     `xml:"branches-covered,attr"`
	BranchesValid   int64     `xml:"branches-valid,attr"`
	Complexity      float32   `xml:"complexity,attr"`
	Version         string    `xml:"version,attr"`
	Timestamp       int64     `xml:"timestamp,attr"`
	Packages        []Package `xml:"packages>package"`
	Sources         []string  `xml:"sources>source"`
}

// Package information
type Package struct {
	Name       string  `xml:"name,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`
	Classes    []Class `xml:"classes>class"`
	LineCount  int64   `xml:"line-count,attr"`
	LineHits   int64   `xml:"line-hits,attr"`
}

// Class information
type Class struct {
	Name       string   `xml:"name,attr"`
	Filename   string   `xml:"filename,attr"`
	LineRate   float32  `xml:"line-rate,attr"`
	BranchRate float32  `xml:"branch-rate,attr"`
	Complexity float32  `xml:"complexity,attr"`
	Methods    []Method `xml:"methods>method"`
	Lines      []Line   `xml:"lines>line"`
	LineCount  int64    `xml:"line-count,attr"`
	LineHits   int64    `xml:"line-hits,attr"`
}

// Method information
type Method struct {
	Name       string  `xml:"name,attr"`
	Signature  string  `xml:"signature,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`
	Lines      []Line  `xml:"lines>line"`
	LineCount  int64   `xml:"line-count,attr"`
	LineHits   int64   `xml:"line-hits,attr"`
}

// Line information
type Line struct {
	Number int   `xml:"number,attr"`
	Hits   int64 `xml:"hits,attr"`
}

func main() {
	sourcePathPtr := flag.String(
		"source",
		"",
		"Absolute path to source. Defaults to current working directory.",
	)

	flag.Parse()

	var sourcePath string

	// Parse the commandline arguments.
	var err error
	if *sourcePathPtr != "" {
		sourcePath = *sourcePathPtr
		if !filepath.IsAbs(sourcePath) {
			panic(fmt.Errorf("Source path is a relative path: %s", sourcePath))
		}
	} else {
		sourcePath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	err = convert.Convert(sourcePath, os.Stdout)
	if err != nil {
		panic(err)
	}
}
