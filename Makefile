all: test

prepare:
	go get -u github.com/axw/gocov/...
	go install -v ./...

test:
	gocov test -v github.com/AlekSi/gocov-xml/internal/package1 > package1.json
	gocov-xml < package1.json > package1.xml
	xmllint --valid --noout package1.xml
	go test ./...
