gocov XML (Cobertura) export
============================

This is a simple helper tool for generating XML output in [Cobertura](http://cobertura.sourceforge.net/) format
for CIs like [Jenkins](https://wiki.jenkins-ci.org/display/JENKINS/Cobertura+Plugin) and others
from [github.com/axw/gocov](https://github.com/axw/gocov) output.

UPDATES: enhanced AlekSi's tool to make it support cobertura coverage-04.dtd format. And make it show correct coverage data on vsts 

Installation
------------

Just type the following to install the program and its dependencies:

    $ go get github.com/axw/gocov/...
    $ go get github.com/yaoyaozong/gocov-xml

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
