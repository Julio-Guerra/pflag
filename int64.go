package pflag

import (
	"reflect"
	"strconv"
)

// -- int64 Value
type Int64Value struct {
	Value int64
	*DefaultValues
}

func NewInt64Value(defaultValue, defaultArg interface{}) *Int64Value {
	v, dv := NewDefaultValues(reflect.TypeOf(int64(0)), defaultValue, defaultArg)
	return &Int64Value{
		Value:         v.(int64),
		DefaultValues: dv,
	}
}

func (f *Int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	f.Value = int64(v)
	return err
}

func (f *Int64Value) Type() string {
	return "int64"
}

func (f *Int64Value) String() string {
	return strconv.FormatInt(f.Value, 10)
}

// Int64Var defines a int64 flag with specified name, default value, and usage string.
// The argument p point64s to a int64 variable in which to store the value of the flag.
func (f *FlagSet) Int64Var(p *Int64Value, name, usage string) *Flag {
	return Int64VarP(p, name, "", usage)
}

// Int64VarP is like Int64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64VarP(p *Int64Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// Int64Var defines a int64 flag with specified name, default value, and usage string.
// The argument p point64s to a int64 variable in which to store the value of the flag.
func Int64Var(p *Int64Value, name, usage string) *Flag {
	return CommandLine.Int64VarP(p, name, "", usage)
}

// Int64VarP is like Int64Var, but accepts a shorthand letter that can be used after a single dash.
func Int64VarP(p *Int64Value, name, shorthand, usage string) *Flag {
	return CommandLine.Int64VarP(p, name, shorthand, usage)
}

// Int64 defines a int64 flag with specified name, default value, and usage string.
// The return value is the address of a int64 variable that stores the value of the flag.
func (f *FlagSet) Int64(name, usage string) *Int64Value {
	return f.Int64P(name, "", usage)
}

// Int64P is like Int64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64P(name, shorthand, usage string) *Int64Value {
	p := NewInt64Value(nil, nil)
	f.Int64VarP(p, name, shorthand, usage)
	return p
}

// Int64 defines a int64 flag with specified name, default value, and usage string.
// The return value is the address of a int64 variable that stores the value of the flag.
func Int64(name, usage string) *Int64Value {
	return CommandLine.Int64P(name, "", usage)
}

// Int64P is like Int64, but accepts a shorthand letter that can be used after a single dash.
func Int64P(name, shorthand, usage string) *Int64Value {
	return CommandLine.Int64P(name, shorthand, usage)
}
