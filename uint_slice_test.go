package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpUISFlagSet(uisp *UIntSliceValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.UIntSliceVar(uisp, "uis", "Command separated list!")
	return f
}

func TestEmptyUIS(t *testing.T) {
	uis := NewUIntSliceValue(nil, nil)
	f := setUpUISFlagSet(uis)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(uis.Value) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", uis.Value, len(uis.Value))
	}
}

func TestUIS(t *testing.T) {
	uis := NewUIntSliceValue(nil, nil)
	f := setUpUISFlagSet(uis)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--uis=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range uis.Value {
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*UIntValue).Value; uint(u) != actual {
			t.Fatalf("expected uis[%d] to be %s but got %d", i, vals[i], actual)
		}
	}
}

func TestUISDefault(t *testing.T) {
	uis := NewUIntSliceValue([]uint{33, 44}, nil)
	f := setUpUISFlagSet(uis)

	vals := []string{"33", "44"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range uis.Value {
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*UIntValue).Value; uint(u) != actual {
			t.Fatalf("expect uis[%d] to be %d but got: %d", i, u, actual)
		}
	}
}

func TestUISWithDefault(t *testing.T) {
	uis := NewUIntSliceValue([]uint{33, 55}, nil)
	f := setUpUISFlagSet(uis)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--uis=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range uis.Value {
		u, err := strconv.ParseUint(vals[i], 10, 0)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*UIntValue).Value; uint(u) != actual {
			t.Fatalf("expected uis[%d] to be %d but got: %d", i, u, actual)
		}
	}
}

func TestUISCalledTwice(t *testing.T) {
	uis := NewUIntSliceValue(nil, nil)
	f := setUpUISFlagSet(uis)

	in := []string{"1,2", "3"}
	expected := []int{1, 2, 3}
	argfmt := "--uis=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range uis.Value {
		if actual := v.(*UIntValue).Value; uint(expected[i]) != actual {
			t.Fatalf("expected uis[%d] to be %d but got: %d", i, expected[i], actual)
		}
	}
}
