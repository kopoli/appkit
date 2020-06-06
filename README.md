# Appkit for Go applications

Very basic scaffolding library for Go applications.

Features:
- Options datastore mostly for command line flags.
- Version string generation.
- Printing of licenses of the dependencies.

## Usage

To use the features, have the following kind of main.go in the program

```
package main

//go:generate licrep -o licenses.go

import (
	"fmt"
	"os"
	"github.com/kopoli/appkit"
)

var (
	version     = "Undefined"
	timestamp   = "Undefined"
	buildGOOS   = "Undefined"
	buildGOARCH = "Undefined"
	progVersion = "" + version
)

// ...

func main() {

	opts := appkit.NewOptions()
	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)

	printVersion := true
	printLicenses := true

	if printVersion {
		fmt.Println(appkit.VersionString(opts))
		os.Exit(0)
	}

	if printLicenses {
		l, err := GetLicenses()
		// if err ...
		s, err := appkit.LicenseString(l)
		// if err ...
		fmt.Print(s)
		os.Exit(0)
	}
}
```

Install the dependencies:

```
$ go get github.com/kopoli/licrep
$ go get github.com/kopoli/gobu
```

[Licrep](https://github.com/kopoli/licrep) generates the licenses.go. The
reason for including licenses is that if one distributes only binaries, e.g.
MIT license requires the license text is present.

[Gobu](https://github.com/kopoli/gobu) builds the application and
automatically fills the version information from git tags.

Therefore generate the application binary with:

```
$ go generate
$ gobu version
```

