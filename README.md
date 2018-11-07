# Go Tour

A Tour of Go is an introduction to the Go programming language. Visit
https://tour.golang.org to start the tour.

## Download/Install

To install the tour from source, first
[set up a workspace](https://golang.org/doc/code.html) and then run:

	$ go get golang.org/x/tour

This will place a `tour` binary in your workspace's `bin` directory, which
can be run offline.

## Contributing

Contributions should follow the same procedure as for the Go project:
https://golang.org/doc/contribute.html

To run the tour server locally:

```sh
dev_appserver.py app.yaml
```

and then visit [http://localhost:8080/](http://localhost:8080) in your browser.

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the tour is located at
https://github.com/golang/go/issues. Prefix your issue with "tour:" in the
subject line, so it is easy to find.

## License

Unless otherwise noted, the go-tour source files are distributed
under the BSD-style license found in the LICENSE file.
