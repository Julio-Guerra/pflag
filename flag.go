// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package pflag is a drop-in replacement for Go's flag package, implementing
POSIX/GNU-style --flags.

pflag is compatible with the GNU extensions to the POSIX recommendations
for command-line options. See
http://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html

Usage:

pflag is a drop-in replacement of Go's native flag package. If you import
pflag under the name "flag" then all code should continue to function
with no changes.

	import flag "github.com/spf13/pflag"

There is one exception to this: if you directly instantiate the Flag struct
there is one more field "Shorthand" that you will need to set.
Most code never instantiates this struct directly, and instead uses
functions such as String(), BoolVar(), and Var(), and is therefore
unaffected.

Define flags using flag.String(), Bool(), Int(), etc.

This declares an integer flag, -flagname, stored in the pointer ip, with type *int.
	var ip = flag.Int("flagname", 1234, "help message for flagname")
If you like, you can bind the flag to a variable using the Var() functions.
	var flagvar int
	func init() {
		flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	}
Or you can create custom flags that satisfy the Value interface (with
pointer receivers) and couple them to flag parsing by
	flag.Var(&flagVal, "name", "help message for flagname")
For such flags, the default value is just the initial value of the variable.

After all flags are defined, call
	flag.Parse()
to parse the command line into the defined flags.

Flags may then be used directly. If you're using the flags themselves,
they are all pointers; if you bind to variables, they're values.
	fmt.Println("ip has value ", *ip)
	fmt.Println("flagvar has value ", flagvar)

After parsing, the arguments after the flag are available as the
slice flag.Args() or individually as flag.Arg(i).
The arguments are indexed from 0 through flag.NArg()-1.

The pflag package also defines some new functions that are not in flag,
that give one-letter shorthands for flags. You can use these by appending
'P' to the name of any function that defines a flag.
	var ip = flag.IntP("flagname", "f", 1234, "help message")
	var flagvar bool
	func init() {
		flag.BoolVarP("boolname", "b", true, "help message")
	}
	flag.VarP(&flagVar, "varname", "v", 1234, "help message")
Shorthand letters can be used with single dashes on the command line.
Boolean shorthand flags can be combined with other shorthand flags.

Command line flag syntax:
	--flag    // boolean flags only
	--flag=x

Unlike the flag package, a single dash before an option means something
different than a double dash. Single dashes signify a series of shorthand
letters for flags. All but the last shorthand letter must be boolean flags.
	// boolean flags
	-f
	-abc
	// non-boolean flags
	-n 1234
	-Ifile
	// mixed
	-abcs "hello"
	-abcn1234

Flag parsing stops after the terminator "--". Unlike the flag package,
flags can be interspersed with arguments anywhere on the command line
before this terminator.

Integer flags accept 1234, 0664, 0x1234 and may be negative.
Boolean flags (in their long form) accept 1, 0, t, f, true, false,
TRUE, FALSE, True, False.
Duration flags accept any input valid for time.ParseDuration.

The default set of command-line flags is controlled by
top-level functions.  The FlagSet type allows one to define
independent sets of flags, such as to implement subcommands
in a command-line interface. The methods of FlagSet are
analogous to the top-level functions for the command-line
flag set.
*/
package pflag

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ErrHelp is the error returned if the flag -help is invoked but no such flag is defined.
var ErrHelp = errors.New("pflag: help requested")

type getFlagValueError struct {
	name         string
	expectedType string
	actualType   string
}

func (err *getFlagValueError) Error() string {
	return fmt.Sprintf("trying to get %s value out of flag %s with type %s", err.expectedType, err.name, err.actualType)
}

// ErrorHandling defines how to handle flag parsing errors.
type ErrorHandling int

const (
	// ContinueOnError will return an err from Parse() if an error is found
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) if an error is found when parsing
	ExitOnError
	// PanicOnError will panic() if an error is found when parsing flags
	PanicOnError
)

// NormalizedName is a flag name that has been normalized according to rules
// for the FlagSet (e.g. making '-' and '_' equivalent).
type NormalizedName string

// A FlagSet represents a set of defined flags.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler.
	Usage func()

	// SortFlags is used to indicate, if user wants to have sorted flags in
	// help/usage messages.
	SortFlags bool

	name              string
	parsed            bool
	actual            map[NormalizedName]*Flag
	orderedActual     []*Flag
	sortedActual      []*Flag
	formal            map[NormalizedName]*Flag
	orderedFormal     []*Flag
	sortedFormal      []*Flag
	shorthands        map[byte]*Flag
	args              []string // arguments after flags
	argsLenAtDash     int      // len(args) when a '--' was located when parsing, or -1 if no --
	errorHandling     ErrorHandling
	output            io.Writer // nil means stderr; use out() accessor
	interspersed      bool      // allow interspersed option/non-option args
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName
}

// A Flag represents the state of a flag.
type Flag struct {
	Name                string              // name as it appears on command line
	Shorthand           string              // one-letter abbreviated flag
	Usage               string              // help message
	ExpectArg           bool                // the flag expects an argument
	Value               Value               // value as set
	Deprecated          string              // If this flag is deprecated, this string is the new or now thing to use
	Hidden              bool                // used by cobra.Command to allow flags to be hidden from help/usage text
	ShorthandDeprecated string              // If the shorthand of this flag is deprecated, this string is the new or now thing to use
	Present             bool                // The flag is present in the argument list
	Annotations         map[string][]string // used by cobra.Command bash autocomple code
}

// Value is the interface to the dynamic value stored in a flag.
// When its flag does expect an argument and is present, Set() is called with
// argument "true".
type Value interface {
	Set(string) error
	Type() string
	fmt.Stringer
	DefaultValue
}

// sortFlags returns the flags as a slice in lexicographical sorted order.
func sortFlags(flags map[NormalizedName]*Flag) []*Flag {
	list := make(sort.StringSlice, len(flags))
	i := 0
	for k := range flags {
		list[i] = string(k)
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list {
		result[i] = flags[NormalizedName(name)]
	}
	return result
}

// SetNormalizeFunc allows you to add a function which can translate flag names.
// Flags added to the FlagSet will be translated and then when anything tries to
// look up the flag that will also be translated. So it would be possible to create
// a flag named "getURL" and have it translated to "geturl".  A user could then pass
// "--getUrl" which may also be translated to "geturl" and everything will work.
func (f *FlagSet) SetNormalizeFunc(n func(f *FlagSet, name string) NormalizedName) {
	f.normalizeNameFunc = n
	f.sortedFormal = f.sortedFormal[:0]
	for k, v := range f.orderedFormal {
		delete(f.formal, NormalizedName(v.Name))
		nname := f.normalizeFlagName(v.Name)
		v.Name = string(nname)
		f.formal[nname] = v
		f.orderedFormal[k] = v
	}
}

// GetNormalizeFunc returns the previously set NormalizeFunc of a function which
// does no translation, if not set previously.
func (f *FlagSet) GetNormalizeFunc() func(f *FlagSet, name string) NormalizedName {
	if f.normalizeNameFunc != nil {
		return f.normalizeNameFunc
	}
	return func(f *FlagSet, name string) NormalizedName { return NormalizedName(name) }
}

func (f *FlagSet) normalizeFlagName(name string) NormalizedName {
	n := f.GetNormalizeFunc()
	return n(f, name)
}

func (f *FlagSet) out() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}

// VisitAll visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	if len(f.formal) == 0 {
		return
	}

	var flags []*Flag
	if f.SortFlags {
		if len(f.formal) != len(f.sortedFormal) {
			f.sortedFormal = sortFlags(f.formal)
		}
		flags = f.sortedFormal
	} else {
		flags = f.orderedFormal
	}

	for _, flag := range flags {
		fn(flag)
	}
}

// HasFlags returns a bool to indicate if the FlagSet has any flags definied.
func (f *FlagSet) HasFlags() bool {
	return len(f.formal) > 0
}

// HasAvailableFlags returns a bool to indicate if the FlagSet has any flags
// definied that are not hidden or deprecated.
func (f *FlagSet) HasAvailableFlags() bool {
	for _, flag := range f.formal {
		if !flag.Hidden && len(flag.Deprecated) == 0 {
			return true
		}
	}
	return false
}

// VisitAll visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}

// Visit visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func (f *FlagSet) Visit(fn func(*Flag)) {
	if len(f.actual) == 0 {
		return
	}

	var flags []*Flag
	if f.SortFlags {
		if len(f.actual) != len(f.sortedActual) {
			f.sortedActual = sortFlags(f.actual)
		}
		flags = f.sortedActual
	} else {
		flags = f.orderedActual
	}

	for _, flag := range flags {
		fn(flag)
	}
}

// Visit visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func Visit(fn func(*Flag)) {
	CommandLine.Visit(fn)
}

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) Lookup(name string) *Flag {
	return f.lookup(f.normalizeFlagName(name))
}

// ShorthandLookup returns the Flag structure of the short handed flag,
// returning nil if none exists.
// It panics, if len(name) > 1.
func (f *FlagSet) ShorthandLookup(name string) *Flag {
	if name == "" {
		return nil
	}
	if len(name) > 1 {
		msg := fmt.Sprintf("can not look up shorthand which is more than one ASCII character: %q", name)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	}
	c := name[0]
	return f.shorthands[c]
}

// lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) lookup(name NormalizedName) *Flag {
	return f.formal[name]
}

func (f *FlagSet) getFlag(name string) (*Flag, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return nil, fmt.Errorf("flag accessed but not defined: %s", name)
	}
	return flag, nil
}

// ArgsLenAtDash will return the length of f.Args at the moment when a -- was
// found during arg parsing. This allows your program to know which args were
// before the -- and which came after.
func (f *FlagSet) ArgsLenAtDash() int {
	return f.argsLenAtDash
}

// MarkDeprecated indicated that a flag is deprecated in your program. It will
// continue to function but will not show up in help or usage messages. Using
// this flag will also print the given usageMessage.
func (f *FlagSet) MarkDeprecated(name string, usageMessage string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %q does not exist", name)
	}
	if usageMessage == "" {
		return fmt.Errorf("deprecated message for flag %q must be set", name)
	}
	flag.Deprecated = usageMessage
	return nil
}

// MarkShorthandDeprecated will mark the shorthand of a flag deprecated in your
// program. It will continue to function but will not show up in help or usage
// messages. Using this flag will also print the given usageMessage.
func (f *FlagSet) MarkShorthandDeprecated(name string, usageMessage string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %q does not exist", name)
	}
	if usageMessage == "" {
		return fmt.Errorf("deprecated message for flag %q must be set", name)
	}
	flag.ShorthandDeprecated = usageMessage
	return nil
}

// MarkHidden sets a flag to 'hidden' in your program. It will continue to
// function but will not show up in help or usage messages.
func (f *FlagSet) MarkHidden(name string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %q does not exist", name)
	}
	flag.Hidden = true
	return nil
}

// Lookup returns the Flag structure of the named command-line flag,
// returning nil if none exists.
func Lookup(name string) *Flag {
	return CommandLine.Lookup(name)
}

// ShorthandLookup returns the Flag structure of the short handed flag,
// returning nil if none exists.
func ShorthandLookup(name string) *Flag {
	return CommandLine.ShorthandLookup(name)
}

// Set sets the value of the named flag.
func (f *FlagSet) Set(name, value string) error {
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok {
		return fmt.Errorf("no such flag -%v", name)
	}

	err := flag.Value.Set(value)
	if err != nil {
		var flagName string
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			flagName = fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
		} else {
			flagName = fmt.Sprintf("--%s", flag.Name)
		}
		return fmt.Errorf("invalid argument %q for %q flag: %v", value, flagName, err)
	}

	if f.actual == nil {
		f.actual = make(map[NormalizedName]*Flag)
	}
	f.actual[normalName] = flag
	f.orderedActual = append(f.orderedActual, flag)

	flag.Present = true

	if flag.Deprecated != "" {
		fmt.Fprintf(f.out(), "Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated)
	}
	return nil
}

// SetAnnotation allows one to set arbitrary annotations on a flag in the FlagSet.
// This is sometimes used by spf13/cobra programs which want to generate additional
// bash completion information.
func (f *FlagSet) SetAnnotation(name, key string, values []string) error {
	normalName := f.normalizeFlagName(name)
	flag, ok := f.formal[normalName]
	if !ok {
		return fmt.Errorf("no such flag -%v", name)
	}
	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}
	flag.Annotations[key] = values
	return nil
}

// Changed returns true if the flag was explicitly set during Parse() and false
// otherwise
func (f *FlagSet) Present(name string) bool {
	flag := f.Lookup(name)
	// If a flag doesn't exist, it wasn't changed....
	if flag == nil {
		return false
	}
	return flag.Present
}

// Set sets the value of the named command-line flag.
func Set(name, value string) error {
	return CommandLine.Set(name, value)
}

// PrintDefaults prints, to standard error unless configured
// otherwise, the default values of all defined flags in the set.
func (f *FlagSet) PrintDefaults() {
	usages := f.FlagUsages()
	fmt.Fprint(f.out(), usages)
}

// Usage extracts the usage string for a flag and returns it. The name is an
// educated guess of the type of the flag's value, or the empty string if the
// flag is boolean.
func ExtractUsage(flag *Flag) (name string, usage string) {
	usage = flag.Usage
	name = flag.Value.Type()
	switch name {
	case "bool":
		name = ""
	case "float64":
		name = "float"
	case "int64":
		name = "int"
	case "uint64":
		name = "uint"
	}

	return
}

// Splits the string `s` on whitespace into an initial substring up to
// `i` runes in length and the remainder. Will go `slop` over `i` if
// that encompasses the entire string (which allows the caller to
// avoid short orphan words on the final line).
func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t")
	if w <= 0 {
		return s, ""
	}

	return s[:w], s[w+1:]
}

// Wraps the string `s` to a maximum width `w` with leading indent
// `i`. The first line is not indented (this is assumed to be done by
// caller). Pass `w` == 0 to do no wrapping
func wrap(i, w int, s string) string {
	if w == 0 {
		return s
	}

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	// If still not enough space then don't even try to wrap.
	if wrap < 24 {
		return s
	}

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap = wrap - slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r = r + l

	// Now wrap the rest
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + t
	}

	return r

}

// FlagUsagesWrapped returns a string containing the usage information
// for all flags in the FlagSet. Wrapped to `cols` columns (0 for no
// wrapping)
func (f *FlagSet) FlagUsagesWrapped(cols int) string {
	buf := new(bytes.Buffer)

	lines := make([]string, 0, len(f.formal))

	maxlen := 0
	f.VisitAll(func(flag *Flag) {
		if flag.Deprecated != "" || flag.Hidden {
			return
		}

		line := ""
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			line = fmt.Sprintf("  -%s, --%s", flag.Shorthand, flag.Name)
		} else {
			line = fmt.Sprintf("      --%s", flag.Name)
		}

		varname, usage := ExtractUsage(flag)
		if defVal := flag.Value.DefaultArg(); defVal != "" {
			switch flag.Value.Type() {
			case "string":
				line += fmt.Sprintf("[=\"%s\"]", defVal)
			case "bool":
				if defVal != "true" {
					line += fmt.Sprintf("[=%s]", defVal)
				}
			default:
				line += fmt.Sprintf("[=%s]", defVal)
			}
		} else if varname != "" {
			line += "=" + varname
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}

		line += usage
		if noFlagVal := flag.Value.DefaultValue(); flag.ExpectArg && noFlagVal != "" {
			line += fmt.Sprintf(" (defaults to \"%s\")", noFlagVal)
		}

		lines = append(lines, line)
	})

	for _, line := range lines {
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
	}

	return buf.String()
}

// FlagUsages returns a string containing the usage information for all flags in
// the FlagSet
func (f *FlagSet) FlagUsages() string {
	return f.FlagUsagesWrapped(0)
}

// PrintDefaults prints to standard error the default values of all defined command-line flags.
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// defaultUsage is the default function to print a usage message.
func defaultUsage(f *FlagSet) {
	fmt.Fprintf(f.out(), "Usage of %s:\n", f.name)
	f.PrintDefaults()
}

// NOTE: Usage is not just defaultUsage(CommandLine)
// because it serves (via godoc flag Usage) as the example
// for how to write your own usage function.

// Usage prints to standard error a usage message documenting all defined command-line flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls PrintDefaults; for details about the
// format of the output and how to control it, see the documentation for PrintDefaults.
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}

// NFlag returns the number of flags that have been set.
func (f *FlagSet) NFlag() int { return len(f.actual) }

// NFlag returns the number of command-line flags that have been set.
func NFlag() int { return len(CommandLine.actual) }

// Arg returns the i'th argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}

// Arg returns the i'th command-line argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func Arg(i int) string {
	return CommandLine.Arg(i)
}

// NArg is the number of arguments remaining after flags have been processed.
func (f *FlagSet) NArg() int { return len(f.args) }

// NArg is the number of arguments remaining after flags have been processed.
func NArg() int { return len(CommandLine.args) }

// Args returns the non-flag arguments.
func (f *FlagSet) Args() []string { return f.args }

// Args returns the non-flag command-line arguments.
func Args() []string { return CommandLine.args }

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(value Value, name string, expectArg bool, usage string) *Flag {
	return f.VarP(value, name, "", expectArg, usage)
}

// VarPF is like VarP, but returns the flag created
func (f *FlagSet) VarP(value Value, name, shorthand string, expectArg bool, usage string) *Flag {
	// Remember the default value as a string; it won't change.
	flag := &Flag{
		Name:      name,
		Shorthand: shorthand,
		Usage:     usage,
		Value:     value,
		ExpectArg: expectArg,
	}
	f.AddFlag(flag)
	return flag
}

// AddFlag will add the flag to the FlagSet
func (f *FlagSet) AddFlag(flag *Flag) {
	normalizedFlagName := f.normalizeFlagName(flag.Name)

	_, alreadyThere := f.formal[normalizedFlagName]
	if alreadyThere {
		msg := fmt.Sprintf("%s flag redefined: %s", f.name, flag.Name)
		fmt.Fprintln(f.out(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[NormalizedName]*Flag)
	}

	flag.Name = string(normalizedFlagName)
	f.formal[normalizedFlagName] = flag
	f.orderedFormal = append(f.orderedFormal, flag)

	if flag.Shorthand == "" {
		return
	}
	if len(flag.Shorthand) > 1 {
		msg := fmt.Sprintf("%q shorthand is more than one ASCII character", flag.Shorthand)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	}
	if f.shorthands == nil {
		f.shorthands = make(map[byte]*Flag)
	}
	c := flag.Shorthand[0]
	used, alreadyThere := f.shorthands[c]
	if alreadyThere {
		msg := fmt.Sprintf("unable to redefine %q shorthand in %q flagset: it's already used for %q flag", c, f.name, used.Name)
		fmt.Fprintf(f.out(), msg)
		panic(msg)
	}
	f.shorthands[c] = flag
}

// AddFlagSet adds one FlagSet to another. If a flag is already present in f
// the flag from newSet will be ignored.
func (f *FlagSet) AddFlagSet(newSet *FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(flag *Flag) {
		if f.Lookup(flag.Name) == nil {
			f.AddFlag(flag)
		}
	})
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func Var(value Value, name string, expectArg bool, usage string) *Flag {
	return VarP(value, name, "", expectArg, usage)
}

// VarP is like Var, but accepts a shorthand letter that can be used after a single dash.
func VarP(value Value, name, shorthand string, expectArg bool, usage string) *Flag {
	return CommandLine.VarP(value, name, shorthand, expectArg, usage)
}

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (f *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.out(), err)
	f.usage()
	return err
}

// usage calls the Usage method for the flag set, or the usage function if
// the flag set is CommandLine.
func (f *FlagSet) usage() {
	if f == CommandLine {
		Usage()
	} else if f.Usage == nil {
		defaultUsage(f)
	} else {
		f.Usage()
	}
}

func (f *FlagSet) parseLongOpt(opt string, args []string, fn ParseFunc) ([]string, error) {
	o := opt[2:]
	if o[0] == '-' || o[0] == '=' {
		return args, f.failf("bad flag syntax: %s", opt)
	}

	split := strings.SplitN(o, "=", 2)
	o = split[0]
	flag, exists := f.formal[f.normalizeFlagName(o)]
	if !exists {
		if o == "help" { // special case for nice help message.
			f.usage()
			return args, ErrHelp
		}
		return args, f.failf("unknown flag: --%s", o)
	}

	var optArg string
	if !flag.ExpectArg {
		optArg = "true"
	} else {
		if len(split) == 2 { // '--option=arg'
			optArg = split[1]
		} else if len(args) > 0 && (args[0] == "" || args[0][0] != '-') { // '--option' 'arg'
			optArg = args[0]
			args = args[1:]
		} else if da := flag.Value.DefaultArg(); da != "" { // '--option'
			optArg = da
		} else {
			return args, f.failf("flag needs an argument: %s", o, opt)
		}
	}

	return args, fn(flag, optArg)
}

func (f *FlagSet) parseShortOpt(opt string, args []string, fn ParseFunc) ([]string, error) {
	if strings.HasPrefix(opt, "test.") {
		return args, nil
	}

	for i := 1; i < len(opt); i++ {
		o := opt[i]
		flag, exists := f.shorthands[o]
		if !exists {
			if o == 'h' { // special help message case
				f.usage()
				return args, ErrHelp
			}
			return args, f.failf("unknown shorthand flag: %c in %s", o, opt)
		}

		if flag.ShorthandDeprecated != "" {
			fmt.Fprintf(f.out(), "Flag shorthand -%s has been deprecated, %s\n", flag.Shorthand, flag.ShorthandDeprecated)
		}

		var optArg string
		if !flag.ExpectArg { // '-XoX' 'X'
			optArg = "true"
		} else {
			if len(opt[i+1:]) == 0 { // '-Xo'
				if len(args) > 0 && (args[0] == "" || args[0][0] != '-') { // '-Xo' 'arg'
					optArg = args[0]
					args = args[1:]
				} else if da := flag.Value.DefaultArg(); da != "" { // '-Xo' '-X'
					optArg = da
				} else {
					return args, f.failf("flag needs an argument: -%c in %s", o, opt)
				}
			} else { // '-Xoarg'
				optArg = opt[i+1:]
				i = len(opt)
			}
		}

		if err := fn(flag, optArg); err != nil {
			return args, err
		}
	}
	return args, nil
}

func (f *FlagSet) parseArgs(args []string, fn ParseFunc) (err error) {
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if len(s) <= 1 || s[0] != '-' {
			if !f.interspersed {
				f.args = append(f.args, s)
				f.args = append(f.args, args...)
				return nil
			}
			f.args = append(f.args, s)
			continue
		}

		if len(s) == 2 && s == "--" { // "--" terminates the flags
			f.argsLenAtDash = len(f.args)
			f.args = append(f.args, args...)
			break
		}

		switch s[0:2] {
		case "--":
			args, err = f.parseLongOpt(s, args, fn)
		default:
			args, err = f.parseShortOpt(s, args, fn)
		}
		if err != nil {
			return
		}
	}
	return
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	return f.ParseAll(arguments, func(flag *Flag, value string) error {
		return f.Set(flag.Name, strings.Trim(value, " "))
	})
}

type ParseFunc func(flag *Flag, value string) error

// ParseAll parses flag definitions from the argument list, which should not
// include the command name. The arguments for fn are flag and value. Must be
// called after all flags in the FlagSet are defined and before flags are
// accessed by the program. The return value will be ErrHelp if -help was set
// but not defined.
func (f *FlagSet) ParseAll(arguments []string, fn ParseFunc) error {
	f.parsed = true
	f.args = make([]string, 0, len(arguments))

	err := f.parseArgs(arguments, fn)
	if err != nil {
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// Parsed reports whether f.Parse has been called.
func (f *FlagSet) Parsed() bool {
	return f.parsed
}

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
}

// ParseAll parses the command-line flags from os.Args[1:] and called fn for each.
// The arguments for fn are flag and value. Must be called after all flags are
// defined and before flags are accessed by the program.
func ParseAll(fn func(flag *Flag, value string) error) {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.ParseAll(os.Args[1:], fn)
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func SetInterspersed(interspersed bool) {
	CommandLine.SetInterspersed(interspersed)
}

// Parsed returns true if the command-line flags have been parsed.
func Parsed() bool {
	return CommandLine.Parsed()
}

// CommandLine is the default set of command-line flags, parsed from os.Args.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// NewFlagSet returns a new, empty flag set with the specified name,
// error handling property and SortFlags set to true.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
		argsLenAtDash: -1,
		interspersed:  true,
		SortFlags:     true,
	}
	return f
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func (f *FlagSet) SetInterspersed(interspersed bool) {
	f.interspersed = interspersed
}

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name
	f.errorHandling = errorHandling
	f.argsLenAtDash = -1
}
