package pflag

import "strconv"

// -- count Value
type CountValue struct {
	*IntValue
	changed bool
}

func NewCountValue(defaultValue, defaultArg int) *CountValue {
	return &CountValue{
		IntValue: NewIntValue(defaultValue, defaultArg),
		changed:  false,
	}
}

func NewStandardCountValue() *CountValue {
	return NewCountValue(0, 1)
}

func (n *CountValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return err
	}
	vword := int(v)
	if n.changed {
		vword += n.IntValue.Value
	} else {
		n.changed = true
	}

	n.IntValue.Value = vword
	return err
}

func (_ *CountValue) Type() string {
	return "count"
}

// CountVar defines a count flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
// A count flag will add 1 to its value evey time it is found on the command line
func (f *FlagSet) CountVar(p *CountValue, name string, usage string) {
	f.CountVarP(p, name, "", usage)
}

// CountVarP is like CountVar only take a shorthand for the flag name.
func (f *FlagSet) CountVarP(p *CountValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// CountVar like CountVar only the flag is placed on the CommandLine instead of a given flag set
func CountVar(p *CountValue, name, usage string) {
	CommandLine.CountVar(p, name, usage)
}

// CountVarP is like CountVar only take a shorthand for the flag name.
func CountVarP(p *CountValue, name, shorthand, usage string) {
	CommandLine.CountVarP(p, name, shorthand, usage)
}

// Count defines a count flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value evey time it is found on the command line
func (f *FlagSet) Count(name string, usage string) *CountValue {
	return f.CountP(name, "", usage)
}

// CountP is like Count only takes a shorthand for the flag name.
func (f *FlagSet) CountP(name, shorthand, usage string) *CountValue {
	p := NewStandardCountValue()
	f.CountVarP(p, name, shorthand, usage)
	return p
}

// Count defines a count flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value evey time it is found on the command line
func Count(name, usage string) *CountValue {
	return CommandLine.CountP(name, "", usage)
}

// CountP is like Count only takes a shorthand for the flag name.
func CountP(name, shorthand string, usage string) *CountValue {
	return CommandLine.CountP(name, shorthand, usage)
}
