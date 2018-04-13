package pflag

import (
	"reflect"
	"strconv"
)

// -- uint64 Value
type UInt64Value struct {
	Value uint64
	*DefaultValues
}

func NewUInt64Value(defaultValue, defaultArg interface{}) *UInt64Value {
	v, dv := NewDefaultValues(reflect.TypeOf(uint64(0)), defaultValue, defaultArg)
	return &UInt64Value{
		Value:         v.(uint64),
		DefaultValues: dv,
	}
}

func (f *UInt64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	f.Value = uint64(v)
	return err
}

func (f *UInt64Value) Type() string {
	return "uint64"
}

func (f *UInt64Value) String() string {
	return strconv.FormatUint(f.Value, 10)
}

// UInt64Var defines a uint64 flag with specified name, default value, and usage string.
// The argument p pouint64s to a uint64 variable in which to store the value of the flag.
func (f *FlagSet) UInt64Var(p *UInt64Value, name, usage string) *Flag {
	return f.UInt64Var(p, name, usage)
}

// UInt64VarP is like UInt64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt64VarP(p *UInt64Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// UInt64Var defines a uint64 flag with specified name, default value, and usage string.
// The argument p pouint64s to a uint64 variable in which to store the value of the flag.
func UInt64Var(p *UInt64Value, name, usage string) *Flag {
	return CommandLine.UInt64Var(p, name, usage)
}

// UInt64VarP is like UInt64Var, but accepts a shorthand letter that can be used after a single dash.
func UInt64VarP(p *UInt64Value, name, shorthand, usage string) *Flag {
	return CommandLine.UInt64VarP(p, name, shorthand, usage)
}

// UInt64 defines a uint64 flag with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the flag.
func (f *FlagSet) UInt64(name, usage string) *UInt64Value {
	return f.UInt64P(name, "", usage)
}

// UInt64P is like UInt64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt64P(name, shorthand, usage string) *UInt64Value {
	p := NewUInt64Value(nil, nil)
	f.UInt64VarP(p, name, shorthand, usage)
	return p
}

// UInt64 defines a uint64 flag with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the flag.
func UInt64(name, usage string) *UInt64Value {
	return CommandLine.UInt64P(name, "", usage)
}

// UInt64P is like UInt64, but accepts a shorthand letter that can be used after a single dash.
func UInt64P(name, shorthand, usage string) *UInt64Value {
	return CommandLine.UInt64P(name, shorthand, usage)
}
