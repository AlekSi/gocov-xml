all: fvb

prepare:
	go get -u github.com/axw/gocov/...
	gocov test -v io > io.json

fvb:
	gofmt -e -s -w .
	go vet .
	go run ./gocov-xml.go < io.json > coverage.xml
	xmllint --valid --noout coverage.xml
