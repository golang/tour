# Go Tour

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/tour.svg)](https://pkg.go.dev/golang.org/x/tour)

A Tour of Go is an introduction to the Go programming language. Visit
https://tour.golang.org to start the tour.

## Download/Install

To run the tour locally, you'll first need to install
[Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) and
[Go](https://golang.org/doc/install) and then run:

	$ go get golang.org/x/tour

This will place a `tour` binary in your
[workspace](https://golang.org/doc/code.html#Workspaces)'s `bin` directory.
The tour program can be run offline.

## Contributing

Contributions should follow the same procedure as for the Go project:
https://golang.org/doc/contribute.html

To run the tour server locally:

```sh
go run .
```

Your browser should now open. If not, please visit [http://localhost:3999/](http://localhost:3999).


## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The issue tracker for the tour's code is located at https://github.com/golang/go/issues.
Prefix your issue with "x/tour:" in the subject line, so it is easy to find.

Issues with the tour's content itself should be reported in the issue tracker
at https://github.com/golang/tour/issues.

## Deploying

1.	To deploy tour.golang.org, run:

	```
	GO111MODULE=on gcloud --project=golang-org app deploy --no-promote app.yaml
	```

	This will create a new version, which can be viewed within the
	[golang-org GCP project](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=tour).

2.	Check that the deployed version looks OK (click the version link in GCP).

3.	If all is well, click "Migrate Traffic" to move 100% of the tour.golang.org
	traffic to the new version.

4.	You're done.

## License

Unless otherwise noted, the go-tour source files are distributed
under the BSD-style license found in the LICENSE file.
