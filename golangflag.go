// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"reflect"
	"strings"
)

// flagValueWrapper implements pflag.Value around a flag.Value.  The main
// difference here is the addition of the Type method that returns a string
// name of the type.  As this is generally unknown, we approximate that with
// reflection.
type flagValueWrapper struct {
	inner    goflag.Value
	flagType string
	*DefaultValues
}

// A goflag is a bool flag when it implements a IsBoolFlag() bool method
// returning true.
type goBoolFlag interface {
	goflag.Value
	IsBoolFlag() bool
}

func wrapFlagValue(v goflag.Value, defVal string) *flagValueWrapper {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Interface, reflect.Ptr:
		t = t.Elem()
	}

	var (
		defaultArg, defaultValue interface{}
		typ                      reflect.Type
	)
	// a goflag is a bool flag when it implement a IsBoolFlag() method...
	if fv, ok := v.(goBoolFlag); ok && fv.IsBoolFlag() {
		defaultArg = true
		typ = reflect.TypeOf(true)
	} else if defVal != "" {
		typ = reflect.TypeOf("")
		defaultValue = defVal
	}
	_, dv := NewDefaultValues(typ, defaultValue, defaultArg)
	if dft := dv.DefaultValue(); dft != "" {
		v.Set(dft)
	}

	return &flagValueWrapper{
		inner:         v,
		flagType:      strings.TrimSuffix(t.Name(), "Value"),
		DefaultValues: dv,
	}
}

func (v *flagValueWrapper) String() string {
	return v.inner.String()
}

func (v *flagValueWrapper) Set(s string) error {
	return v.inner.Set(s)
}

func (v *flagValueWrapper) Type() string {
	return v.flagType
}

// PFlagFromGoFlag will return a *pflag.Flag given a *flag.Flag
// If the *flag.Flag.Name was a single character (ex: `v`) it will be accessiblei
// with both `-v` and `--v` in flags. If the golang flag was more than a single
// character (ex: `verbose`) it will only be accessible via `--verbose`
func PFlagFromGoFlag(goflag *goflag.Flag) *Flag {
	wrapped := wrapFlagValue(goflag.Value, goflag.DefValue)
	// Remember the default value as a string; it won't change.
	flag := &Flag{
		Name:      goflag.Name,
		Usage:     goflag.Usage,
		Value:     wrapped,
		ExpectArg: true,
	}
	// Ex: if the golang flag was -v, allow both -v and --v to work
	if len(flag.Name) == 1 {
		flag.Shorthand = flag.Name
	}
	return flag
}

// AddGoFlag will add the given *flag.Flag to the pflag.FlagSet
func (f *FlagSet) AddGoFlag(goflag *goflag.Flag) {
	if f.Lookup(goflag.Name) != nil {
		return
	}
	newflag := PFlagFromGoFlag(goflag)
	f.AddFlag(newflag)
}

// AddGoFlagSet will add the given *flag.FlagSet to the pflag.FlagSet
func (f *FlagSet) AddGoFlagSet(newSet *goflag.FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(goflag *goflag.Flag) {
		f.AddGoFlag(goflag)
	})
}
