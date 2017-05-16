package pflag

import (
	"reflect"
	"strconv"
)

// -- float32 Value
type Float32Value struct {
	Value float32
	*DefaultValues
}

func NewFloat32Value(defaultValue, defaultArg interface{}) *Float32Value {
	v, dv := NewDefaultValues(reflect.TypeOf(float32(0)), defaultValue, defaultArg)
	return &Float32Value{
		Value:         v.(float32),
		DefaultValues: dv,
	}
}

func (f *Float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	f.Value = float32(v)
	return err
}

func (f *Float32Value) Type() string {
	return "float32"
}

func (f *Float32Value) String() string {
	return strconv.FormatFloat(float64(f.Value), 'g', -1, 32)
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func (f *FlagSet) Float32Var(p *Float32Value, name, usage string) *Flag {
	return f.Float32VarP(p, name, "", usage)
}

// Float32VarP is like Float32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32VarP(p *Float32Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func Float32Var(p *Float32Value, name, usage string) *Flag {
	return Float32VarP(p, name, "", usage)
}

// Float32VarP is like Float32Var, but accepts a shorthand letter that can be used after a single dash.
func Float32VarP(p *Float32Value, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func (f *FlagSet) Float32(name, usage string) *Float32Value {
	return f.Float32P(name, "", usage)
}

// Float32P is like Float32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32P(name, shorthand, usage string) *Float32Value {
	p := NewFloat32Value(nil, nil)
	f.Float32VarP(p, name, shorthand, usage)
	return p
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func Float32(name, usage string) *Float32Value {
	return CommandLine.Float32P(name, "", usage)
}

// Float32P is like Float32, but accepts a shorthand letter that can be used after a single dash.
func Float32P(name, shorthand, usage string) *Float32Value {
	return CommandLine.Float32P(name, shorthand, usage)
}
