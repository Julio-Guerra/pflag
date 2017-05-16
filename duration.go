package pflag

import (
	"reflect"
	"time"
)

// -- time.Duration Value
type DurationValue struct {
	Value time.Duration
	*DefaultValues
}

func NewDurationValue(defaultValue, defaultArg interface{}) *DurationValue {
	v, dv := NewDefaultValues(reflect.TypeOf(time.Duration(0)), defaultValue, defaultArg)
	return &DurationValue{
		Value:         v.(time.Duration),
		DefaultValues: dv,
	}
}

func (d *DurationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	d.Value = v
	return nil
}

func (d *DurationValue) Type() string {
	return "duration"
}

func (d *DurationValue) String() string {
	return d.Value.String()
}

func durationConv(sval string) (interface{}, error) {
	return time.ParseDuration(sval)
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func (f *FlagSet) DurationVar(p *DurationValue, name, usage string) *Flag {
	return f.VarP(p, name, "", true, usage)
}

// DurationVarP is like DurationVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) DurationVarP(p *DurationValue, name, shorthand, usage string) {
	f.VarP(p, name, shorthand, true, usage)
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func DurationVar(p *DurationValue, name string, value time.Duration, usage string) *Flag {
	return DurationVarP(p, name, "", usage)
}

// DurationVarP is like DurationVar, but accepts a shorthand letter that can be used after a single dash.
func DurationVarP(p *DurationValue, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func (f *FlagSet) Duration(name, usage string) *DurationValue {
	return f.DurationP(name, "", usage)
}

// DurationP is like Duration, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) DurationP(name, shorthand, usage string) *DurationValue {
	p := NewDurationValue(nil, nil)
	f.DurationVarP(p, name, shorthand, usage)
	return p
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func Duration(name, usage string) *DurationValue {
	return CommandLine.DurationP(name, "", usage)
}

// DurationP is like Duration, but accepts a shorthand letter that can be used after a single dash.
func DurationP(name, shorthand, usage string) *DurationValue {
	return CommandLine.DurationP(name, shorthand, usage)
}
