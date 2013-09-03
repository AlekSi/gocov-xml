package main

import (
	"encoding/xml"
	"go/build"
	"io"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

type dirInfo struct {
	PkgPath string
	DirPath string
}

func TestConvertEmpty(t *testing.T) {
	data := `{ "Packages": [] }`

	pipe2rd, pipe2wr := io.Pipe()
	go convert(strings.NewReader(data), pipe2wr)

	v := Coverage{}
	dec := xml.NewDecoder(pipe2rd)
	dec.Decode(&v)

	if v.XMLName.Local != "coverage" {
		t.Fail()
	}
	if v.Packages != nil {
		t.FailNow()
	}
}

func TestConvertFunc1(t *testing.T) {
	tmpl, err := template.ParseFiles("testdata/func1_test.json")
	if err != nil {
		t.FailNow()
	}
	dirInfo := dirInfo{}
	dirInfo.PkgPath = reflect.TypeOf(Coverage{}).PkgPath()
	pkg, err := build.Import(dirInfo.PkgPath, "", build.FindOnly)
	if err != nil {
		t.FailNow()
	}
	dirInfo.DirPath = pkg.Dir

	pipe1rd, pipe1wr := io.Pipe()
	go func() {
		err := tmpl.Execute(pipe1wr, dirInfo)
		if err != nil {
			t.Fail()
			panic("tmpl.Execute failed")
		}
	}()

	pipe2rd, pipe2wr := io.Pipe()
	go convert(pipe1rd, pipe2wr)

	v := Coverage{}
	dec := xml.NewDecoder(pipe2rd)
	dec.Decode(&v)

	if v.XMLName.Local != "coverage" {
		t.Fail()
	}
	if v.Packages == nil || len(v.Packages) != 1 {
		t.FailNow()
	}

	p := v.Packages[0]
	if p.Name != dirInfo.PkgPath+"/testdata" {
		t.Fail()
	}
	if p.Classes == nil || len(p.Classes) != 1 {
		t.FailNow()
	}

	c := p.Classes[0]
	if c.Name != "-" {
		t.Fail()
	}
	if c.Filename != dirInfo.DirPath+"/testdata/func1.go" {
		t.Fail()
	}
	if c.Methods == nil || len(c.Methods) != 1 {
		t.FailNow()
	}
	if c.Lines == nil || len(c.Lines) != 2 {
		t.FailNow()
	}

	m := c.Methods[0]
	if m.Name != "Func1" {
		t.Fail()
	}
	if m.Lines == nil || len(m.Lines) != 2 {
		t.FailNow()
	}

	l1 := m.Lines[0]
	if l1.Number != 5 {
		t.Fail()
	}
	if l1.Hits != 1 {
		t.Fail()
	}

	l2 := m.Lines[1]
	if l2.Number != 6 {
		t.Fail()
	}
	if l2.Hits != 0 {
		t.Fail()
	}

	l1 = c.Lines[0]
	if l1.Number != 5 {
		t.Fail()
	}
	if l1.Hits != 1 {
		t.Fail()
	}

	l2 = c.Lines[1]
	if l2.Number != 6 {
		t.Fail()
	}
}
