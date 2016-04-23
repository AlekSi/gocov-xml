package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
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

	packages := gocov_xml.ConvertGocov(r.Packages)
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
