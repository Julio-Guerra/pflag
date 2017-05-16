package pflag

type IPSliceValue struct {
	*SliceValue
}

func NewIPSliceValue(defaultValue, defaultArg interface{}) *IPSliceValue {
	return &IPSliceValue{
		SliceValue: NewSliceValue(defaultValue, defaultArg, func(d interface{}) Value {
			return NewIPValue(d, nil)
		}),
	}
}

func (f *FlagSet) IPSliceVar(p *IPSliceValue, name, usage string) *Flag {
	return f.IPSliceVarP(p, name, "", usage)
}

func (f *FlagSet) IPSliceVarP(p *IPSliceValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

func IPSliceVar(p *IPSliceValue, name, usage string) *Flag {
	return CommandLine.IPSliceVarP(p, name, "", usage)
}

func IPSliceVarP(p *IPSliceValue, name, shorthand, usage string) *Flag {
	return CommandLine.IPSliceVarP(p, name, shorthand, usage)
}

func (f *FlagSet) IPSlice(name string, usage string) *IPSliceValue {
	return f.IPSliceP(name, "", usage)
}

func (f *FlagSet) IPSliceP(name, shorthand string, usage string) *IPSliceValue {
	p := NewIPSliceValue(nil, nil)
	f.IPSliceVarP(p, name, shorthand, usage)
	return p
}

func IPSlice(name string, usage string) *IPSliceValue {
	return CommandLine.IPSliceP(name, "", usage)
}

func IPSliceP(name, shorthand, usage string) *IPSliceValue {
	return CommandLine.IPSliceP(name, shorthand, usage)
}
