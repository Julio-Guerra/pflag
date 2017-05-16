package pflag

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

// -- net.IPMask value
type IPMaskValue struct {
	Value net.IPMask
	*DefaultValues
}

func NewIPMaskValue(defaultValue, defaultArg interface{}) *IPMaskValue {
	v, dv := NewDefaultValues(reflect.TypeOf(net.IPMask(nil)), defaultValue, defaultArg)
	return &IPMaskValue{
		Value:         v.(net.IPMask),
		DefaultValues: dv,
	}
}

func (i *IPMaskValue) String() string {
	return ((net.IP)(i.Value)).String()
}

func (i *IPMaskValue) Set(s string) error {
	ipmask := net.IPMask(net.ParseIP(strings.TrimSpace(s)))
	if ipmask == nil {
		return fmt.Errorf("failed to parse IPMask: %q", s)
	}
	i.Value = ipmask
	return nil
}

func (i *IPMaskValue) Type() string {
	return "ipmask"
}

// IPMaskVar defines an net.IPMask flag with specified name, default value, and usage string.
// The argument p points to an net.IPMask variable in which to store the value of the flag.
func (f *FlagSet) IPMaskVar(p *IPMaskValue, name string, usage string) *Flag {
	return f.IPMaskVarP(p, name, "", usage)
}

// IPMaskVarP is like IPMaskVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPMaskVarP(p *IPMaskValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// IPMaskVar defines an IPMaskValue flag with specified name, default value, and usage string.
// The argument p points to an IPMaskValue variable in which to store the value of the flag.
func IPMaskVar(p *IPMaskValue, name, usage string) *Flag {
	return CommandLine.IPMaskVarP(p, name, "", usage)
}

// IPMaskVarP is like IPMaskVar, but accepts a shorthand letter that can be used after a single dash.
func IPMaskVarP(p *IPMaskValue, name, shorthand, usage string) *Flag {
	return CommandLine.IPMaskVarP(p, name, shorthand, usage)
}

// IPMask defines an IPMaskValue flag with specified name, default value, and usage string.
// The return value is the address of an IPMaskValue variable that stores the value of the flag.
func (f *FlagSet) IPMask(name string, usage string) *IPMaskValue {
	return f.IPMaskP(name, "", usage)
}

// IPMaskP is like IPMask, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPMaskP(name, shorthand, usage string) *IPMaskValue {
	p := NewIPMaskValue(nil, nil)
	f.IPMaskVarP(p, name, shorthand, usage)
	return p
}

// IPMask defines an IPMaskValue flag with specified name, default value, and usage string.
// The return value is the address of an IPMaskValue variable that stores the value of the flag.
func IPMask(name, usage string) *IPMaskValue {
	return CommandLine.IPMaskP(name, "", usage)
}

// IPMaskP is like IPMask, but accepts a shorthand letter that can be used after a single dash.
func IPMaskP(name, shorthand string, value IPMaskValue, usage string) *IPMaskValue {
	return CommandLine.IPMaskP(name, shorthand, usage)
}
