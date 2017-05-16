package pflag

type BoolSliceValue struct {
	*SliceValue
}

func NewBoolSliceValue(defaultValue, defaultArg []bool) *BoolSliceValue {
	return &BoolSliceValue{
		SliceValue: NewSliceValue(defaultValue, defaultArg, func(d interface{}) Value {
			return NewBoolValue(d, nil)
		}),
	}
}

func (f *FlagSet) BoolSliceVar(p *BoolSliceValue, name, usage string) *Flag {
	return f.BoolSliceVarP(p, name, "", usage)
}

func (f *FlagSet) BoolSliceVarP(p *BoolSliceValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

func BoolSliceVar(p *BoolSliceValue, name, usage string) *Flag {
	return BoolSliceVarP(p, name, "", usage)
}

func BoolSliceVarP(p *BoolSliceValue, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

func (f *FlagSet) BoolSlice(name string, usage string) *BoolSliceValue {
	return f.BoolSliceP(name, "", usage)
}

func (f *FlagSet) BoolSliceP(name, shorthand string, usage string) *BoolSliceValue {
	p := NewBoolSliceValue(nil, nil)
	f.BoolSliceVarP(p, name, shorthand, usage)
	return p
}

func BoolSlice(name string, usage string) *BoolSliceValue {
	return CommandLine.BoolSliceP(name, "", usage)
}

func BoolSliceP(name, shorthand, usage string) *BoolSliceValue {
	return CommandLine.BoolSliceP(name, shorthand, usage)
}
