package pflag

import (
	"reflect"
	"strconv"
)

// -- uint8 Value
type UInt8Value struct {
	Value uint8
	*DefaultValues
}

func NewUInt8Value(defaultValue, defaultArg interface{}) *UInt8Value {
	v, dv := NewDefaultValues(reflect.TypeOf(uint8(0)), defaultValue, defaultArg)
	return &UInt8Value{
		Value:         v.(uint8),
		DefaultValues: dv,
	}
}

func (f *UInt8Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return err
	}
	f.Value = uint8(v)
	return err
}

func (f *UInt8Value) Type() string {
	return "uint8"
}

func (f *UInt8Value) String() string {
	return strconv.FormatUint(uint64(f.Value), 10)
}

// UInt8Var defines a uint8 flag with specified name, default value, and usage string.
// The argument p pouint8s to a uint8 variable in which to store the value of the flag.
func (f *FlagSet) UInt8Var(p *UInt8Value, name, usage string) *Flag {
	return f.UInt8Var(p, name, usage)
}

// UInt8VarP is like UInt8Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt8VarP(p *UInt8Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// UInt8Var defines a uint8 flag with specified name, default value, and usage string.
// The argument p pouint8s to a uint8 variable in which to store the value of the flag.
func UInt8Var(p *UInt8Value, name, usage string) *Flag {
	return CommandLine.UInt8Var(p, name, usage)
}

// UInt8VarP is like UInt8Var, but accepts a shorthand letter that can be used after a single dash.
func UInt8VarP(p *UInt8Value, name, shorthand, usage string) *Flag {
	return CommandLine.UInt8VarP(p, name, shorthand, usage)
}

// UInt8 defines a uint8 flag with specified name, default value, and usage string.
// The return value is the address of a uint8 variable that stores the value of the flag.
func (f *FlagSet) UInt8(name, usage string) *UInt8Value {
	return f.UInt8P(name, "", usage)
}

// UInt8P is like UInt8, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt8P(name, shorthand, usage string) *UInt8Value {
	p := NewUInt8Value(nil, nil)
	f.UInt8VarP(p, name, shorthand, usage)
	return p
}

// UInt8 defines a uint8 flag with specified name, default value, and usage string.
// The return value is the address of a uint8 variable that stores the value of the flag.
func UInt8(name, usage string) *UInt8Value {
	return CommandLine.UInt8P(name, "", usage)
}

// UInt8P is like UInt8, but accepts a shorthand letter that can be used after a single dash.
func UInt8P(name, shorthand, usage string) *UInt8Value {
	return CommandLine.UInt8P(name, shorthand, usage)
}
