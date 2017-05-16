package pflag

import (
	"reflect"
	"strconv"
)

type BoolValue struct {
	*DefaultValues
	Value bool
}

// A boolean value is false by default and its option-argument is defaults to
// true.
func NewBoolValue(defaultValue, defaultArg interface{}) *BoolValue {
	value, defaults := NewDefaultValues(reflect.TypeOf(true), defaultValue, defaultArg)
	return &BoolValue{
		DefaultValues: defaults,
		Value:         value.(bool),
	}
}

func (b *BoolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	b.Value = v
	return nil
}

func (b *BoolValue) Type() string {
	return "bool"
}

func (b *BoolValue) String() string {
	return strconv.FormatBool(b.Value)
}

// BoolVar defines a bool flag with specified name, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (f *FlagSet) BoolVar(p *BoolValue, name string, usage string) {
	f.BoolVarP(p, name, "", usage)
}

// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolVarP(p *BoolValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// BoolVar defines a bool flag with specified name, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func BoolVar(p *BoolValue, name, usage string) {
	BoolVarP(p, name, "", usage)
}

// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
func BoolVarP(p *BoolValue, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func (f *FlagSet) Bool(name, usage string) (*BoolValue, *Flag) {
	return f.BoolP(name, "", usage)
}

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolP(name, shorthand, usage string) (*BoolValue, *Flag) {
	p := NewBoolValue(false, true)
	return p, f.BoolVarP(p, name, shorthand, usage)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(name, usage string) (*BoolValue, *Flag) {
	return BoolP(name, "", usage)
}

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func BoolP(name, shorthand, usage string) (*BoolValue, *Flag) {
	return CommandLine.BoolP(name, shorthand, usage)
}
