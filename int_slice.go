package pflag

type IntSliceValue struct {
	*SliceValue
}

func NewIntSliceValue(defaultValue, defaultArg interface{}) *IntSliceValue {
	return &IntSliceValue{
		SliceValue: NewSliceValue(defaultValue, defaultArg, func(d interface{}) Value {
			return NewIntValue(d, nil)
		}),
	}
}

func (f *FlagSet) IntSliceVar(p *IntSliceValue, name, usage string) *Flag {
	return f.IntSliceVarP(p, name, "", usage)
}

func (f *FlagSet) IntSliceVarP(p *IntSliceValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

func IntSliceVar(p *IntSliceValue, name, usage string) *Flag {
	return IntSliceVarP(p, name, "", usage)
}

func IntSliceVarP(p *IntSliceValue, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

func (f *FlagSet) IntSlice(name string, usage string) *IntSliceValue {
	return f.IntSliceP(name, "", usage)
}

func (f *FlagSet) IntSliceP(name, shorthand string, usage string) *IntSliceValue {
	p := NewIntSliceValue(nil, nil)
	f.IntSliceVarP(p, name, shorthand, usage)
	return p
}

func IntSlice(name string, usage string) *IntSliceValue {
	return CommandLine.IntSliceP(name, "", usage)
}

func IntSliceP(name, shorthand, usage string) *IntSliceValue {
	return CommandLine.IntSliceP(name, shorthand, usage)
}
