package pflag

import (
	"reflect"
	"strconv"
)

// -- int32 Value
type Int32Value struct {
	Value int32
	*DefaultValues
}

func NewInt32Value(defaultValue, defaultArg interface{}) *Int32Value {
	v, dv := NewDefaultValues(reflect.TypeOf(int32(0)), defaultValue, defaultArg)
	return &Int32Value{
		Value:         v.(int32),
		DefaultValues: dv,
	}
}

func (f *Int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return err
	}
	f.Value = int32(v)
	return err
}

func (f *Int32Value) Type() string {
	return "int32"
}

func (f *Int32Value) String() string {
	return strconv.FormatInt(int64(f.Value), 10)
}

// Int32Var defines a int32 flag with specified name, default value, and usage string.
// The argument p point32s to a int32 variable in which to store the value of the flag.
func (f *FlagSet) Int32Var(p *Int32Value, name, usage string) *Flag {
	return Int32VarP(p, name, "", usage)
}

// Int32VarP is like Int32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int32VarP(p *Int32Value, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// Int32Var defines a int32 flag with specified name, default value, and usage string.
// The argument p point32s to a int32 variable in which to store the value of the flag.
func Int32Var(p *Int32Value, name, usage string) *Flag {
	return Int32VarP(p, name, "", usage)
}

// Int32VarP is like Int32Var, but accepts a shorthand letter that can be used after a single dash.
func Int32VarP(p *Int32Value, name, shorthand, usage string) *Flag {
	return CommandLine.Int32VarP(p, name, shorthand, usage)
}

// Int32 defines a int32 flag with specified name, default value, and usage string.
// The return value is the address of a int32 variable that stores the value of the flag.
func (f *FlagSet) Int32(name, usage string) *Int32Value {
	return f.Int32P(name, "", usage)
}

// Int32P is like Int32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int32P(name, shorthand, usage string) *Int32Value {
	p := NewInt32Value(nil, nil)
	f.Int32VarP(p, name, shorthand, usage)
	return p
}

// Int32 defines a int32 flag with specified name, default value, and usage string.
// The return value is the address of a int32 variable that stores the value of the flag.
func Int32(name, usage string) *Int32Value {
	return CommandLine.Int32P(name, "", usage)
}

// Int32P is like Int32, but accepts a shorthand letter that can be used after a single dash.
func Int32P(name, shorthand, usage string) *Int32Value {
	return CommandLine.Int32P(name, shorthand, usage)
}
