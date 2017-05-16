package pflag

// StdBool defines a stdbool flag with specified name, default value, and usage string.
// The return value is the address of a stdbool variable that stores the value of the flag.
func (f *FlagSet) StdBool(name, usage string) (*BoolValue, *Flag) {
	return f.StdBoolP(name, "", usage)
}

// StdBoolP is like StdBool, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StdBoolP(name, shorthand, usage string) (*BoolValue, *Flag) {
	p := NewBoolValue(false, nil)
	return p, f.VarP(p, name, shorthand, false, usage)
}

// StdBool defines a stdbool flag with specified name, default value, and usage string.
// The return value is the address of a stdbool variable that stores the value of the flag.
func StdBool(name, usage string) (*BoolValue, *Flag) {
	return StdBoolP(name, "", usage)
}

// StdBoolP is like StdBool, but accepts a shorthand letter that can be used after a single dash.
func StdBoolP(name, shorthand, usage string) (*BoolValue, *Flag) {
	return CommandLine.StdBoolP(name, shorthand, usage)
}
