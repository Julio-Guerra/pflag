package pflag

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

// -- net.IP value
type IPValue struct {
	Value net.IP
	*DefaultValues
}

func NewIPValue(defaultValue, defaultArg interface{}) *IPValue {
	v, dv := NewDefaultValues(reflect.TypeOf(net.IP(nil)), defaultValue, defaultArg)
	return &IPValue{
		Value:         v.(net.IP),
		DefaultValues: dv,
	}
}

func (i *IPValue) String() string {
	if i.Value == nil {
		return "<nil>"
	}
	return i.Value.String()
}

func (i *IPValue) Set(s string) error {
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil {
		return fmt.Errorf("failed to parse IP: %q", s)
	}
	i.Value = ip
	return nil
}

func (i *IPValue) Type() string {
	return "ip"
}

// IPVar defines an net.IP flag with specified name, default value, and usage string.
// The argument p points to an net.IP variable in which to store the value of the flag.
func (f *FlagSet) IPVar(p *IPValue, name string, usage string) *Flag {
	return f.IPVarP(p, name, "", usage)
}

// IPVarP is like IPVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPVarP(p *IPValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

// IPVar defines an IPValue flag with specified name, default value, and usage string.
// The argument p points to an IPValue variable in which to store the value of the flag.
func IPVar(p *IPValue, name, usage string) *Flag {
	return IPVarP(p, name, "", usage)
}

// IPVarP is like IPVar, but accepts a shorthand letter that can be used after a single dash.
func IPVarP(p *IPValue, name, shorthand, usage string) *Flag {
	return CommandLine.IPVarP(p, name, shorthand, usage)
}

// IP defines an IPValue flag with specified name, default value, and usage string.
// The return value is the address of an IPValue variable that stores the value of the flag.
func (f *FlagSet) IP(name string, usage string) *IPValue {
	return f.IPP(name, "", usage)
}

// IPP is like IP, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPP(name, shorthand, usage string) *IPValue {
	p := NewIPValue(nil, nil)
	f.IPVarP(p, name, shorthand, usage)
	return p
}

// IP defines an IPValue flag with specified name, default value, and usage string.
// The return value is the address of an IPValue variable that stores the value of the flag.
func IP(name, usage string) *IPValue {
	return CommandLine.IPP(name, "", usage)
}

// IPP is like IP, but accepts a shorthand letter that can be used after a single dash.
func IPP(name, shorthand string, value IPValue, usage string) *IPValue {
	return CommandLine.IPP(name, shorthand, usage)
}
