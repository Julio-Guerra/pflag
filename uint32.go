package pflag

import (
	"reflect"
	"strconv"
)

// -- uint32 Value
type UInt32Value struct {
	Value uint32
	*DefaultValues
}

func NewUInt32Value(defaultValue, defaultArg interface{}) *UInt32Value {
	v, dv := NewDefaultValues(reflect.TypeOf(uint32(0)), defaultValue, defaultArg)
	return &UInt32Value{
		Value:         v.(uint32),
		DefaultValues: dv,
	}
}

func (f *UInt32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return err
	}
	f.Value = uint32(v)
	return err
}

func (f *UInt32Value) Type() string {
	return "uint32"
}

func (f *UInt32Value) String() string {
	return strconv.FormatUint(uint64(f.Value), 10)
}

// UInt32Var defines a uint32 flag with specified name, default value, and usage string.
// The argument p pouint32s to a uint32 variable in which to store the value of the flag.
func (f *FlagSet) UInt32Var(p *UInt32Value, name, usage string) *Flag {
	return f.UInt32VarP(p, "", name, usage)
}

// UInt32VarP is like UInt32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt32VarP(p *UInt32Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// UInt32Var defines a uint32 flag with specified name, default value, and usage string.
// The argument p pouint32s to a uint32 variable in which to store the value of the flag.
func UInt32Var(p *UInt32Value, name, usage string) *Flag {
	return CommandLine.UInt32Var(p, name, usage)
}

// UInt32VarP is like UInt32Var, but accepts a shorthand letter that can be used after a single dash.
func UInt32VarP(p *UInt32Value, name, shorthand, usage string) *Flag {
	return CommandLine.UInt32VarP(p, name, shorthand, usage)
}

// UInt32 defines a uint32 flag with specified name, default value, and usage string.
// The return value is the address of a uint32 variable that stores the value of the flag.
func (f *FlagSet) UInt32(name, usage string) *UInt32Value {
	return f.UInt32P(name, "", usage)
}

// UInt32P is like UInt32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UInt32P(name, shorthand, usage string) *UInt32Value {
	p := NewUInt32Value(nil, nil)
	f.UInt32VarP(p, name, shorthand, usage)
	return p
}

// UInt32 defines a uint32 flag with specified name, default value, and usage string.
// The return value is the address of a uint32 variable that stores the value of the flag.
func UInt32(name, usage string) *UInt32Value {
	return CommandLine.UInt32P(name, "", usage)
}

// UInt32P is like UInt32, but accepts a shorthand letter that can be used after a single dash.
func UInt32P(name, shorthand, usage string) *UInt32Value {
	return CommandLine.UInt32P(name, shorthand, usage)
}
