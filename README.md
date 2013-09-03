gocov XML (Cobertura) export
============================

This is a simple helper tool for generating XML output in [Cobertura](http://cobertura.sourceforge.net/) format
for CIs like [Jenkins](https://wiki.jenkins-ci.org/display/JENKINS/Cobertura+Plugin) and others
from [github.com/axw/gocov](https://github.com/axw/gocov) output.

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
