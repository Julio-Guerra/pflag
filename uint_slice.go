package pflag

type UIntSliceValue struct {
	*SliceValue
}

func NewUIntSliceValue(defaultValue, defaultArg []uint) *UIntSliceValue {
	return &UIntSliceValue{
		SliceValue: NewSliceValue(defaultValue, defaultArg, func(d interface{}) Value {
			return NewUIntValue(d, nil)
		}),
	}
}

func (f *FlagSet) UIntSliceVar(p *UIntSliceValue, name, usage string) *Flag {
	return f.UIntSliceVarP(p, name, "", usage)
}

func (f *FlagSet) UIntSliceVarP(p *UIntSliceValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

func UIntSliceVar(p *UIntSliceValue, name, usage string) *Flag {
	return CommandLine.UIntSliceVar(p, name, usage)
}

func UIntSliceVarP(p *UIntSliceValue, name, shorthand, usage string) *Flag {
	return CommandLine.UIntSliceVarP(p, name, shorthand, usage)
}

func (f *FlagSet) UIntSlice(name string, usage string) *UIntSliceValue {
	return f.UIntSliceP(name, "", usage)
}

func (f *FlagSet) UIntSliceP(name, shorthand string, usage string) *UIntSliceValue {
	p := NewUIntSliceValue(nil, nil)
	f.UIntSliceVarP(p, name, shorthand, usage)
	return p
}

func UIntSlice(name string, usage string) *UIntSliceValue {
	return CommandLine.UIntSliceP(name, "", usage)
}

func UIntSliceP(name, shorthand, usage string) *UIntSliceValue {
	return CommandLine.UIntSliceP(name, shorthand, usage)
}
