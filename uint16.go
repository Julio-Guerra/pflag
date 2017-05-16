package pflag

import (
	"reflect"
	"strconv"
)

// -- uint16 Value
type UInt16Value struct {
	Value uint16
	*DefaultValues
}

func NewUInt16Value(defaultValue, defaultArg interface{}) *UInt16Value {
	v, dv := NewDefaultValues(reflect.TypeOf(uint16(0)), defaultValue, defaultArg)
	return &UInt16Value{
		Value:         v.(uint16),
		DefaultValues: dv,
	}
}

func (f *UInt16Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return err
	}
	f.Value = uint16(v)
	return err
}

func (f *UInt16Value) Type() string {
	return "uint16"
}

func (f *UInt16Value) String() string {
	return strconv.FormatUint(uint64(f.Value), 10)
}

// UInt16Var defines a uint16 flag with specified name, default value, and usage string.
// The argument p pouint16s to a uint16 variable in which to store the value of the flag.
func (f *FlagSet) UInt16Var(p *UInt16Value, name, usage string) *Flag {
	return f.UInt16VarP(p, name, "", usage)
}

// UInt16VarP is like UInt16Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt16VarP(p *UInt16Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// UInt16Var defines a uint16 flag with specified name, default value, and usage string.
// The argument p pouint16s to a uint16 variable in which to store the value of the flag.
func UInt16Var(p *UInt16Value, name, usage string) *Flag {
	return CommandLine.UInt16Var(p, name, usage)
}

// UInt16VarP is like UInt16Var, but accepts a shorthand letter that can be used after a single dash.
func UInt16VarP(p *UInt16Value, name, shorthand, usage string) *Flag {
	return CommandLine.UInt16VarP(p, name, shorthand, usage)
}

// UInt16 defines a uint16 flag with specified name, default value, and usage string.
// The return value is the address of a uint16 variable that stores the value of the flag.
func (f *FlagSet) UInt16(name, usage string) *UInt16Value {
	return f.UInt16P(name, "", usage)
}

// UInt16P is like UInt16, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt16P(name, shorthand, usage string) *UInt16Value {
	p := NewUInt16Value(nil, nil)
	f.UInt16VarP(p, name, shorthand, usage)
	return p
}

// UInt16 defines a uint16 flag with specified name, default value, and usage string.
// The return value is the address of a uint16 variable that stores the value of the flag.
func UInt16(name, usage string) *UInt16Value {
	return CommandLine.UInt16P(name, "", usage)
}

// UInt16P is like UInt16, but accepts a shorthand letter that can be used after a single dash.
func UInt16P(name, shorthand, usage string) *UInt16Value {
	return CommandLine.UInt16P(name, shorthand, usage)
}
