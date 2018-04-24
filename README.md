gocov XML (Cobertura) export
============================

This is a simple helper tool for generating XML output in [Cobertura](http://cobertura.sourceforge.net/) format
for CIs like [Jenkins](https://wiki.jenkins-ci.org/display/JENKINS/Cobertura+Plugin), [vsts](https://www.visualstudio.com/team-services) and others
from [github.com/axw/gocov](https://github.com/axw/gocov) output.
The generated XML output is in the latest [coverage-04.dtd](http://cobertura.sourceforge.net/xml/coverage-04.dtd) schema

Installation
------------

Just type the following to install the program and its dependencies:

    $ go get github.com/axw/gocov/...
    $ go get github.com/AlekSi/gocov-xml

Usage
-----

`gocov-xml` reads from the standard input:

    $ gocov test github.com/gorilla/mux | gocov-xml > coverage.xml

Authors
-------

* [Alexey Palazhchenko (AlekSi)](https://github.com/AlekSi)
* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
* [Marin Bek (marinbek)](https://github.com/marinbek)
* [Alex Castle (acastle)](https://github.com/acastle)
* [Billy Yao (yaoyaozong)](https://github.com/yaoyaozong)
