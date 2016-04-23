all: test

prepare:
	go get -u -v github.com/stretchr/testify/...
	go get -u -v github.com/axw/gocov/...

test:
	rm -f *.json *.xml
	go install -v ./...
	gocov test -v github.com/AlekSi/gocov-xml/internal/test/package1 > package1.gocov.json
	gocov test -v github.com/AlekSi/gocov-xml/internal/test/... > package12.gocov.json
	gocov-xml < package1.gocov.json > package1.gocov.xml
	gocov-xml < package12.gocov.json > package12.gocov.xml
	xmllint --valid --noout --load-trace package1.gocov.xml
	go test ./...
