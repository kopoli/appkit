package main

// This is an example program to showcase appkit features.

// To execute this program, run the following commands:
//
// go get https://github.com/kopoli/licrep
// go get https://github.com/kopoli/gobu
//
// go generate
// gobu
//
// ./_demo -help

//go:generate licrep -o licenses.go

import (
	"flag"
	"fmt"
	"os"

	"github.com/kopoli/appkit"
)

var (
	version     = "Undefined"
	timestamp   = "Undefined"
	buildGOOS   = "Undefined"
	buildGOARCH = "Undefined"
	progVersion = "demo-" + version
)

// ...

func main() {
	handleError := func(err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed with error: %v\n", err)
			os.Exit(1)
		}
	}
	var err error

	opts := appkit.NewOptions()
	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)

	// Define the top level command
	base := appkit.NewCommand(nil, "", "Example program using appkit")
	optVersion := base.Flags.Bool("version", false, "Display version")

	// Define some sub-commands
	_ = appkit.NewCommand(base, "hello-world hello hw",
		"Display Hello World!")
	licenses := appkit.NewCommand(base, "show-licenses licenses",
		"Show licenses of this program")
	optVerbose := licenses.Flags.Bool("verbose", false, "Enable verbose output")

	// Define deep sub-commands
	level1 := appkit.NewCommand(base, "level1", "Only a stepping stone")
	_ = appkit.NewCommand(level1, "deadend", "Nothing here")

	level2 := appkit.NewCommand(level1, "level2", "Show deeper level")
	_ = appkit.NewCommand(level2, "end", "Unnecessarily deep")
	// Display the only command available instead of just [COMMAND] in the help
	level2.SubCommandHelp = "[end]"

	// Parse the command line
	err = base.Parse(os.Args[1:], opts)
	if err == flag.ErrHelp {
		os.Exit(0)
	}

	// Print the program version with a top-level option
	if *optVersion {
		fmt.Println(appkit.VersionString(opts))
		os.Exit(0)
	}

	// Get the actual command and residual command line arguments
	cmd := opts.Get("cmdline-command", "")
	argstr := opts.Get("cmdline-args", "")
	args := appkit.SplitArguments(argstr)

	switch cmd {
	case "":
		fmt.Println("Default command invoked.")
		fmt.Printf("Arguments as a string: %s\n", argstr)
		fmt.Println("Arguments as an array", args)
	case "hello-world":
		fmt.Println("Hello World !")
	case "show-licenses":
		l, err := GetLicenses()
		handleError(err)

		var s string
		if *optVerbose {
			s, err = appkit.LicenseString(l)
		} else {
			s, err = appkit.LicenseFormatString(l, "%Lp: %Ln\n")
		}
		handleError(err)
		fmt.Print(s)
	default:
		fmt.Println("Reached the default handler for commands")
		fmt.Println("Invoked command:", cmd)
		fmt.Println("Arguments as an array", args)
	}

	os.Exit(0)
}
