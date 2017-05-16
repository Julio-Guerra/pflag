package pflag

type StringSliceValue struct {
	*SliceValue
}

func NewStringSliceValue(defaultValue, defaultArg []string) *StringSliceValue {
	return &StringSliceValue{
		SliceValue: NewSliceValue(defaultValue, defaultArg, func(d interface{}) Value {
			return NewStringValue(d, nil)
		}),
	}
}

func (f *FlagSet) StringSliceVar(p *StringSliceValue, name, usage string) *Flag {
	return f.StringSliceVarP(p, name, "", usage)
}

func (f *FlagSet) StringSliceVarP(p *StringSliceValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

func StringSliceVar(p *StringSliceValue, name, usage string) *Flag {
	return CommandLine.StringSliceVarP(p, name, "", usage)
}

func StringSliceVarP(p *StringSliceValue, name, shorthand, usage string) *Flag {
	return CommandLine.StringSliceVarP(p, name, shorthand, usage)
}

func (f *FlagSet) StringSlice(name string, usage string) *StringSliceValue {
	return f.StringSliceP(name, "", usage)
}

func (f *FlagSet) StringSliceP(name, shorthand string, usage string) *StringSliceValue {
	p := NewStringSliceValue(nil, nil)
	f.StringSliceVarP(p, name, shorthand, usage)
	return p
}

func StringSlice(name string, usage string) *StringSliceValue {
	return CommandLine.StringSliceP(name, "", usage)
}

func StringSliceP(name, shorthand, usage string) *StringSliceValue {
	return CommandLine.StringSliceP(name, shorthand, usage)
}
