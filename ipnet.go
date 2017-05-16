package pflag

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

// -- net.IPNet value
type IPNetValue struct {
	Value *net.IPNet
	*DefaultValues
}

func NewIPNetValue(defaultValue, defaultArg interface{}) *IPNetValue {
	v, dv := NewDefaultValues(reflect.TypeOf((*net.IPNet)(nil)), defaultValue, defaultArg)
	return &IPNetValue{
		Value:         v.(*net.IPNet),
		DefaultValues: dv,
	}
}

func (i *IPNetValue) String() string {
	return i.Value.String()
}

func (i *IPNetValue) Set(s string) error {
	_, ipnet, err := net.ParseCIDR(strings.TrimSpace(s))
	if err != nil {
		return fmt.Errorf("failed to parse IPNet: %q", s)
	}
	i.Value = ipnet
	return nil
}

func (i *IPNetValue) Type() string {
	return "ipNetwork"
}

// IPNetVar defines an net.IPNet flag with specified name, default value, and usage string.
// The argument p points to an net.IPNet variable in which to store the value of the flag.
func (f *FlagSet) IPNetVar(p *IPNetValue, name string, usage string) *Flag {
	return f.IPNetVarP(p, name, "", usage)
}

// IPNetVarP is like IPNetVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPNetVarP(p *IPNetValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// IPNetVar defines an IPNetValue flag with specified name, default value, and usage string.
// The argument p points to an IPNetValue variable in which to store the value of the flag.
func IPNetVar(p *IPNetValue, name, usage string) *Flag {
	return CommandLine.IPNetVarP(p, name, "", usage)
}

// IPNetVarP is like IPNetVar, but accepts a shorthand letter that can be used after a single dash.
func IPNetVarP(p *IPNetValue, name, shorthand, usage string) *Flag {
	return CommandLine.IPNetVarP(p, name, shorthand, usage)
}

// IPNet defines an IPNetValue flag with specified name, default value, and usage string.
// The return value is the address of an IPNetValue variable that stores the value of the flag.
func (f *FlagSet) IPNet(name string, usage string) *IPNetValue {
	return f.IPNetP(name, "", usage)
}

// IPNetP is like IPNet, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPNetP(name, shorthand, usage string) *IPNetValue {
	p := NewIPNetValue(nil, nil)
	f.IPNetVarP(p, name, shorthand, usage)
	return p
}

// IPNet defines an IPNetValue flag with specified name, default value, and usage string.
// The return value is the address of an IPNetValue variable that stores the value of the flag.
func IPNet(name, usage string) *IPNetValue {
	return CommandLine.IPNetP(name, "", usage)
}

// IPNetP is like IPNet, but accepts a shorthand letter that can be used after a single dash.
func IPNetP(name, shorthand string, value IPNetValue, usage string) *IPNetValue {
	return CommandLine.IPNetP(name, shorthand, usage)
}
