package pflag

import (
	"reflect"
	"strconv"
)

// -- int8 Value
type Int8Value struct {
	Value int8
	*DefaultValues
}

func NewInt8Value(defaultValue, defaultArg interface{}) *Int8Value {
	v, dv := NewDefaultValues(reflect.TypeOf(int8(0)), defaultValue, defaultArg)
	return &Int8Value{
		Value:         v.(int8),
		DefaultValues: dv,
	}
}

func (f *Int8Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return err
	}
	f.Value = int8(v)
	return err
}

func (f *Int8Value) Type() string {
	return "int8"
}

func (f *Int8Value) String() string {
	return strconv.FormatInt(int64(f.Value), 10)
}

// Int8Var defines a int8 flag with specified name, default value, and usage string.
// The argument p point8s to a int8 variable in which to store the value of the flag.
func (f *FlagSet) Int8Var(p *Int8Value, name, usage string) *Flag {
	return f.Int8VarP(p, name, "", usage)
}

// Int8VarP is like Int8Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int8VarP(p *Int8Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// Int8Var defines a int8 flag with specified name, default value, and usage string.
// The argument p point8s to a int8 variable in which to store the value of the flag.
func Int8Var(p *Int8Value, name, usage string) *Flag {
	return Int8VarP(p, name, "", usage)
}

// Int8VarP is like Int8Var, but accepts a shorthand letter that can be used after a single dash.
func Int8VarP(p *Int8Value, name, shorthand, usage string) *Flag {
	return CommandLine.Int8VarP(p, name, shorthand, usage)
}

// Int8 defines a int8 flag with specified name, default value, and usage string.
// The return value is the address of a int8 variable that stores the value of the flag.
func (f *FlagSet) Int8(name, usage string) *Int8Value {
	return f.Int8P(name, "", usage)
}

// Int8P is like Int8, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int8P(name, shorthand, usage string) *Int8Value {
	p := NewInt8Value(nil, nil)
	f.Int8VarP(p, name, shorthand, usage)
	return p
}

// Int8 defines a int8 flag with specified name, default value, and usage string.
// The return value is the address of a int8 variable that stores the value of the flag.
func Int8(name, usage string) *Int8Value {
	return CommandLine.Int8P(name, "", usage)
}

// Int8P is like Int8, but accepts a shorthand letter that can be used after a single dash.
func Int8P(name, shorthand, usage string) *Int8Value {
	return CommandLine.Int8P(name, shorthand, usage)
}
