package pflag

import (
	"reflect"
)

// -- string Value
type StringValue struct {
	Value string
	*DefaultValues
}

func NewStringValue(defaultValue, defaultArg interface{}) *StringValue {
	v, dv := NewDefaultValues(reflect.TypeOf(""), defaultValue, defaultArg)
	return &StringValue{
		Value:         v.(string),
		DefaultValues: dv,
	}
}

func (f *StringValue) Set(s string) error {
	f.Value = s
	return nil
}

func (f *StringValue) Type() string {
	return "string"
}

func (f *StringValue) String() string {
	return f.Value
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p postrings to a string variable in which to store the value of the flag.
func (f *FlagSet) StringVar(p *StringValue, name, usage string) *Flag {
	return f.StringVarP(p, name, "", usage)
}

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringVarP(p *StringValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p postrings to a string variable in which to store the value of the flag.
func StringVar(p *StringValue, name, usage string) *Flag {
	return CommandLine.StringVarP(p, name, "", usage)
}

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
func StringVarP(p *StringValue, name, shorthand, usage string) *Flag {
	return CommandLine.StringVarP(p, name, shorthand, usage)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name, usage string) *StringValue {
	return f.StringP(name, "", usage)
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringP(name, shorthand, usage string) *StringValue {
	p := NewStringValue(nil, nil)
	f.StringVarP(p, name, shorthand, usage)
	return p
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name, usage string) *StringValue {
	return CommandLine.StringP(name, "", usage)
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func StringP(name, shorthand, usage string) *StringValue {
	return CommandLine.StringP(name, shorthand, usage)
}
