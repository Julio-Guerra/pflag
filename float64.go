package pflag

import (
	"reflect"
	"strconv"
)

// -- float64 Value
type Float64Value struct {
	Value float64
	*DefaultValues
}

func NewFloat64Value(defaultValue, defaultArg interface{}) *Float64Value {
	v, dv := NewDefaultValues(reflect.TypeOf(float64(0)), defaultValue, defaultArg)
	return &Float64Value{
		Value:         v.(float64),
		DefaultValues: dv,
	}
}

func (f *Float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	f.Value = v
	return err
}

func (f *Float64Value) Type() string {
	return "float64"
}

func (f *Float64Value) String() string {
	return strconv.FormatFloat(f.Value, 'g', -1, 64)
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func (f *FlagSet) Float64Var(p *Float64Value, name, usage string) *Flag {
	return f.Float64VarP(p, name, "", usage)
}

// Float64VarP is like Float64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float64VarP(p *Float64Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func Float64Var(p *Float64Value, name, usage string) *Flag {
	return Float64VarP(p, name, "", usage)
}

// Float64VarP is like Float64Var, but accepts a shorthand letter that can be used after a single dash.
func Float64VarP(p *Float64Value, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func (f *FlagSet) Float64(name, usage string) *Float64Value {
	return f.Float64P(name, "", usage)
}

// Float64P is like Float64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float64P(name, shorthand, usage string) *Float64Value {
	p := NewFloat64Value(nil, nil)
	f.Float64VarP(p, name, shorthand, usage)
	return p
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func Float64(name, usage string) *Float64Value {
	return CommandLine.Float64P(name, "", usage)
}

// Float64P is like Float64, but accepts a shorthand letter that can be used after a single dash.
func Float64P(name, shorthand, usage string) *Float64Value {
	return CommandLine.Float64P(name, shorthand, usage)
}
