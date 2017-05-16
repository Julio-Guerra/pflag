// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpISFlagSet(isp *IntSliceValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.IntSliceVar(isp, "is", "Command separated list!")
	return f
}

func TestEmptyIS(t *testing.T) {
	is := NewIntSliceValue(nil, nil)
	f := setUpISFlagSet(is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(is.Value) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", is.Value, len(is.Value))
	}
}

func TestIS(t *testing.T) {
	is := NewIntSliceValue(nil, nil)
	f := setUpISFlagSet(is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is.Value {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*IntValue).Value; d != actual {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], actual)
		}
	}
}

func TestISDefault(t *testing.T) {
	is := NewIntSliceValue([]int{0, 1}, nil)
	f := setUpISFlagSet(is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is.Value {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*IntValue).Value; d != actual {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, actual)
		}
	}
}

func TestISWithDefault(t *testing.T) {
	is := NewIntSliceValue([]int{0, 1}, nil)
	f := setUpISFlagSet(is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range is.Value {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*IntValue).Value; d != actual {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, actual)
		}
	}
}

func TestISCalledTwice(t *testing.T) {
	is := NewIntSliceValue(nil, nil)
	f := setUpISFlagSet(is)

	in := []string{"1,2", "3"}
	expected := []int{1, 2, 3}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range is.Value {
		if actual := v.(*IntValue).Value; expected[i] != actual {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, expected[i], actual)
		}
	}
}
