package pflag

import (
	"reflect"
	"strconv"
)

// -- uint Value
type UIntValue struct {
	Value uint
	*DefaultValues
}

func NewUIntValue(defaultValue, defaultArg interface{}) *UIntValue {
	v, dv := NewDefaultValues(reflect.TypeOf(uint(0)), defaultValue, defaultArg)
	return &UIntValue{
		Value:         v.(uint),
		DefaultValues: dv,
	}
}

func (f *UIntValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 0)
	if err != nil {
		return err
	}
	f.Value = uint(v)
	return err
}

func (f *UIntValue) Type() string {
	return "uint"
}

func (f *UIntValue) String() string {
	return strconv.FormatUint(uint64(f.Value), 10)
}

// UIntVar defines a uint flag with specified name, default value, and usage string.
// The argument p pouints to a uint variable in which to store the value of the flag.
func (f *FlagSet) UIntVar(p *UIntValue, name, usage string) *Flag {
	return f.UIntVarP(p, name, "", usage)
}

// UIntVarP is like UIntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UIntVarP(p *UIntValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// UIntVar defines a uint flag with specified name, default value, and usage string.
// The argument p pouints to a uint variable in which to store the value of the flag.
func UIntVar(p *UIntValue, name, usage string) *Flag {
	return CommandLine.UIntVar(p, name, usage)
}

// UIntVarP is like UIntVar, but accepts a shorthand letter that can be used after a single dash.
func UIntVarP(p *UIntValue, name, shorthand, usage string) *Flag {
	return CommandLine.UIntVarP(p, name, shorthand, usage)
}

// UInt defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint variable that stores the value of the flag.
func (f *FlagSet) UInt(name, usage string) *UIntValue {
	return f.UIntP(name, "", usage)
}

// UIntP is like UInt, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UIntP(name, shorthand, usage string) *UIntValue {
	p := NewUIntValue(nil, nil)
	f.UIntVarP(p, name, shorthand, usage)
	return p
}

// UInt defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint variable that stores the value of the flag.
func UInt(name, usage string) *UIntValue {
	return CommandLine.UIntP(name, "", usage)
}

// UIntP is like UInt, but accepts a shorthand letter that can be used after a single dash.
func UIntP(name, shorthand, usage string) *UIntValue {
	return CommandLine.UIntP(name, shorthand, usage)
}
