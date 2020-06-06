package appkit

import (
	"fmt"
	"runtime"
)

// VersionString returns a version string that should be printed with the -v
// or the --version flag. It gets the components from the following keys from
// the options:
// program-name
// program-version
// program-timestamp
func VersionString(opts Options) string {
	progName := opts.Get("program-name", "undefined")
	progVersion := opts.Get("program-version", "undefined")
	progTimestamp := opts.Get("program-timestamp", "undefined")

	return fmt.Sprintf("%s: %s\nBuilt %v with: %s/%s for %s/%s",
		progName, progVersion, progTimestamp, runtime.Compiler,
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
