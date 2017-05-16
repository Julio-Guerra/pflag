package pflag

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"strings"
)

type SliceValue struct {
	Value []Value
	*DefaultValues
	newValue NewValueFunc
	changed  bool
}

type NewValueFunc func(defaultValue interface{}) Value

func NewSliceValue(defaultValue, defaultArg interface{}, newValue NewValueFunc) *SliceValue {
	var dv interface{}
	if defaultValue != nil {
		var dvSlice []Value
		dvValue := reflect.ValueOf(defaultValue)
		for i := 0; i < dvValue.Len(); i++ {
			dvSlice = append(dvSlice, newValue(dvValue.Index(i).Interface()))
		}
	}

	var da interface{}
	if defaultArg != nil {
		var daSlice []Value
		daValue := reflect.ValueOf(defaultArg)
		for i := 0; i < daValue.Len(); i++ {
			daSlice = append(daSlice, newValue(daValue.Index(i).Interface()))
		}
	}

	v, defaults := NewDefaultValues(reflect.TypeOf([]Value{}), dv, da)
	return &SliceValue{
		Value:         v.([]Value),
		newValue:      newValue,
		DefaultValues: defaults,
	}
}

func (s *SliceValue) Set(val string) (err error) {
	// read flag arguments with CSV parser
	strSlice, err := readAsCSV(val)
	if err != nil && err != io.EOF {
		return err
	}

	out := make([]Value, len(strSlice))
	for i, str := range strSlice {
		out[i] = s.newValue(nil)
		if err := out[i].Set(strings.TrimSpace(str)); err != nil {
			return err
		}
	}

	if s.changed {
		s.Value = append(s.Value, out...)
	} else {
		s.changed = true
		s.Value = out
	}
	return nil
}

func (s *SliceValue) Type() string {
	return s.newValue(nil).Type() + "Slice"
}

func (s *SliceValue) String() string {
	if s.Value == nil {
		return "[]"
	}
	strSlice := make([]string, len(s.Value))
	for i, b := range s.Value {
		strSlice[i] = b.String()
	}

	out, _ := writeAsCSV(strSlice)
	return "[" + out + "]"
}

func (f *FlagSet) SliceVar(p *SliceValue, name, usage string) *Flag {
	return f.SliceVarP(p, name, "", usage)
}

func (f *FlagSet) SliceVarP(p *SliceValue, name, shorthand, usage string) *Flag {
	return f.VarP(p, name, shorthand, true, usage)
}

func SliceVar(p *SliceValue, name, usage string) *Flag {
	return CommandLine.SliceVarP(p, name, "", usage)
}

func SliceVarP(p *SliceValue, name, shorthand, usage string) *Flag {
	return CommandLine.VarP(p, name, shorthand, true, usage)
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func writeAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}
