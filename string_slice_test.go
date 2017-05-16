// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strings"
	"testing"
)

func setUpSSFlagSet(ssp *StringSliceValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringSliceVar(ssp, "ss", "Comma-separated list!")
	return f
}

func TestEmptySS(t *testing.T) {
	ss := NewStringSliceValue(nil, nil)
	f := setUpSSFlagSet(ss)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if l := len(ss.Value); l != 0 {
		t.Fatalf("got ss %v with len=%d but expected length=0", ss.Value, l)
	}
}

func TestEmptySSValue(t *testing.T) {
	ss := NewStringSliceValue(nil, nil)
	f := setUpSSFlagSet(ss)
	err := f.Parse([]string{"--ss="})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if l := len(ss.Value); l != 0 {
		t.Fatalf("got ss %v with len=%d but expected length=0", ss.Value, l)
	}
}

func TestSS(t *testing.T) {
	ss := NewStringSliceValue(nil, nil)
	f := setUpSSFlagSet(ss)

	vals := []string{"one", "two", "4", "3"}
	arg := fmt.Sprintf("--ss=%s", strings.Join(vals, ","))

	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range ss.Value {
		if actual := v.(*StringValue).Value; vals[i] != actual {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], actual)
		}
	}
}

func TestSSDefault(t *testing.T) {
	vals := []string{"default", "value"}
	ss := NewStringSliceValue(vals, nil)
	f := setUpSSFlagSet(ss)

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ss.Value {
		if actual := v.(*StringValue).Value; vals[i] != actual {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], actual)
		}
	}
}

func TestSSWithDefault(t *testing.T) {
	ss := NewStringSliceValue([]string{"default", "value"}, nil)
	f := setUpSSFlagSet(ss)

	vals := []string{"one", "two", "4", "3"}
	arg := fmt.Sprintf("--ss=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range ss.Value {
		if actual := v.(*StringValue).Value; vals[i] != actual {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], actual)
		}
	}
}

func TestSSCalledTwice(t *testing.T) {
	ss := NewStringSliceValue(nil, nil)
	f := setUpSSFlagSet(ss)

	in := []string{"one, two", "three"}
	expected := []string{"one", "two", "three"}
	argfmt := "--ss=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(expected) != len(ss.Value) {
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss.Value))
	}

	for i, v := range ss.Value {
		if actual := v.(*StringValue).Value; expected[i] != actual {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], actual)
		}
	}
}

func TestSSWithComma(t *testing.T) {
	ss := NewStringSliceValue(nil, nil)
	f := setUpSSFlagSet(ss)

	in := []string{`"one,two"`, `"three"`, `"four,five",six`}
	expected := []string{"one,two", "three", "four,five", "six"}
	argfmt := "--ss=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	arg3 := fmt.Sprintf(argfmt, in[2])
	err := f.Parse([]string{arg1, arg2, arg3})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(expected) != len(ss.Value) {
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss.Value))
	}
	for i, v := range ss.Value {
		if actual := v.(*StringValue).Value; expected[i] != actual {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], actual)
		}
	}
}

func TestSSWithSquareBrackets(t *testing.T) {
	ss := NewStringSliceValue(nil, nil)
	f := setUpSSFlagSet(ss)

	in := []string{`"[a-z]"`, `"[a-z]+"`}
	expected := []string{"[a-z]", "[a-z]+"}
	argfmt := "--ss=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(expected) != len(ss.Value) {
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss.Value))
	}
	for i, v := range ss.Value {
		if actual := v.(*StringValue).Value; expected[i] != actual {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], actual)
		}
	}
}
