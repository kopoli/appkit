package appkit

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

var cmdParseRegex = regexp.MustCompile(`[A-Za-z0-9][-A-Za-z0-9]*`)

type Command struct {
	Cmd         []string
	Help        string
	Flags       *flag.FlagSet
	subCommands []*Command
}

// HasFlags returns true if a flag.FlagSet has any flags.
func HasFlags(fs *flag.FlagSet) bool {
	ret := false
	fs.VisitAll(func(f *flag.Flag) {
		ret = true
	})
	return ret
}

// SplitCommand splits a command string to a command and its synonyms. See
// NewCommand for more information.
func SplitCommand(cmdstr string) []string {
	return strings.Fields(cmdstr)
}

// SplitArguments splits the arguments from the option "cmdline-args". See
// Parse for details.
func SplitArguments(argstr string) []string {
	return strings.Split(argstr, "\000")
}

// JoinArguments is the counter operation for SplitArguments. See that for
// more details.
func JoinArguments(args []string) string {
	return strings.Join(args, "\000")
}

// NewCommand creates a recursive command line argument with flags.
//
// Parent of the top-level command should be nil. The cmd string can contain
// multiple space-separated commands that are regarded as synonyms of the
// command. The help string is displayed if help option is given.
//
// If parent == nil, then the Usage function prints out all sub-commands with
// helps. This can be overridden by re-defining the Flags.Usage function.
//
// Example:
//
//   opts := appkit.NewOptions()
//   base := appkit.NewCommand(nil, "", "")
//   optVersion := base.Flags.Bool("version", false, "Display version")
//   add := appkit.NewCommand(base, "add a", "Adding stuff")
//   _ = appkit.NewCommand(add, "package p", "Add package")
//   _ = appkit.NewCommand(add, "dependency d", "Add dependency")
//   del := appkit.NewCommand(base, "delete del d", "Deleting stuff")
//   _ = appkit.NewCommand(del, "package p", "Delete package")
//   _ = appkit.NewCommand(del, "dependency d", "Delete dependency")
//   optRecurse := del.Flags.Bool("recurse", false, "Delete recursively")
//
//   err = base.Parse(os.Args[1:], opts)
//   if err == flag.ErrHelp {
//      os.Exit(0)
//   }
//
//   if *optVersion {
//   fmt.Println(appkit.VersionString(opts))
//      os.Exit(0)
//   }
//   cmd := opts.Get("cmdline-command", "")
//   switch cmd {
//   case "add package":
//   ...
//   case "delete package":
//   ...
//   }
func NewCommand(parent *Command, cmd string, help string) *Command {
	cmds := []string{""}
	if len(cmd) > 0 {
		cmds = SplitCommand(cmd)

		for i := range cmds {
			if !cmdParseRegex.MatchString(cmds[i]) {
				s := fmt.Sprintf("Error: Could not parse command: %s", cmds[i])
				panic(s)
			}
		}
	}
	flags := flag.NewFlagSet(cmds[0], flag.ContinueOnError)

	ret := &Command{
		Cmd:   cmds,
		Help:  help,
		Flags: flags,
	}

	flags.Usage = func() {
		out := flags.Output()
		if parent == nil {
			fmt.Fprintf(out,
				"%s: %s\n\nCommands:\n",
				os.Args[0], ret.Help)
			ret.CommandList(out)
		} else {
			fmt.Fprintf(out, "Command: %s\n\n%s\n",
				strings.Join(ret.Cmd, ", "), ret.Help)
		}
		if HasFlags(flags) {
			fmt.Fprintf(out, "\nOptions:\n")
			flags.PrintDefaults()
		}
	}

	if parent != nil {
		parent.subCommands = append(parent.subCommands, ret)
	}
	return ret
}

// Parse the command line arguments according to the recursive command
// structure.
//
// The actual command will be set to the option "cmdline-command" inside the
// opts structure. Additional positional arguments after the command are in
// the "cmdline-args" option.
//
// Recursive commands in "cmdline-command" are space separated. If there are
// multiple synonyms defined for a command, the first one is listed.
//
// The positional arguments in "cmdline-args" are NUL separated inside the
// string. They can be split to an array using SplitArguments.
func (c *Command) Parse(args []string, opts Options) error {
	var err error
	if c.Flags != nil {
		err = c.Flags.Parse(args)
		if err != nil {
			return err
		}
	}

	args = c.Flags.Args()

	cmd := ""
	if c.Cmd[0] != "" {
		cmd = opts.Get("cmdline-command", "") + " " + c.Cmd[0]
		cmd = strings.TrimLeft(cmd, " ")
	}

	opts.Set("cmdline-command", cmd)
	opts.Set("cmdline-args", JoinArguments(args))

	if len(args) == 0 {
		return nil
	}

	for _, sc := range c.subCommands {
		for i := range sc.Cmd {
			if sc.Cmd[i] == args[0] {
				err = sc.Parse(args[1:], opts)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CommandList prints out a recursive tree of sub-commands to the given
// io.Writer.
func (c *Command) CommandList(out io.Writer) {
	wr := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	var printall func(pfx string, c *Command)
	printall = func(pfx string, c *Command) {
		if c.Cmd[0] != "" {
			fmt.Fprintf(wr, "%s%s\t-\t%s\n", pfx, strings.Join(c.Cmd, ", "), c.Help)
		}
		for _, sc := range c.subCommands {
			printall(pfx+"  ", sc)
		}
	}

	printall("", c)
	wr.Flush()
}
