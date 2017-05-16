package pflag

import (
	"reflect"
	"strconv"
)

// -- int Value
type IntValue struct {
	Value int
	*DefaultValues
}

func NewIntValue(defaultValue, defaultArg interface{}) *IntValue {
	v, dv := NewDefaultValues(reflect.TypeOf(int(0)), defaultValue, defaultArg)
	return &IntValue{
		Value:         v.(int),
		DefaultValues: dv,
	}
}

func (f *IntValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return err
	}
	vword := int(v)
	f.Value = vword
	return err
}

func (f *IntValue) Type() string {
	return "int"
}

func (f *IntValue) String() string {
	return strconv.FormatInt(int64(f.Value), 10)
}

// IntVar defines a int flag with specified name, default value, and usage string.
// The argument p points to a int variable in which to store the value of the flag.
func (f *FlagSet) IntVar(p *IntValue, name, usage string) *Flag {
	return f.IntVarP(p, name, "", usage)
}

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntVarP(p *IntValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// IntVar defines a int flag with specified name, default value, and usage string.
// The argument p points to a int variable in which to store the value of the flag.
func IntVar(p *IntValue, name, usage string) *Flag {
	return IntVarP(p, name, "", usage)
}

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func IntVarP(p *IntValue, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

// Int defines a int flag with specified name, default value, and usage string.
// The return value is the address of a int variable that stores the value of the flag.
func (f *FlagSet) Int(name, usage string) *IntValue {
	return f.IntP(name, "", usage)
}

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntP(name, shorthand, usage string) *IntValue {
	p := NewIntValue(nil, nil)
	f.IntVarP(p, name, shorthand, usage)
	return p
}

// Int defines a int flag with specified name, default value, and usage string.
// The return value is the address of a int variable that stores the value of the flag.
func Int(name, usage string) *IntValue {
	return CommandLine.IntP(name, "", usage)
}

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func IntP(name, shorthand, usage string) *IntValue {
	return CommandLine.IntP(name, shorthand, usage)
}
