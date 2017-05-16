package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpBSFlagSet(bsp *BoolSliceValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(bsp, "bs", "Command separated list!")
	return f
}

func TestEmptyBS(t *testing.T) {
	bs := NewBoolSliceValue(nil, nil)
	f := setUpBSFlagSet(bs)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(bs.Value) != 0 {
		t.Fatalf("got bs %v with len=%d but expected length=0", bs.Value, len(bs.Value))
	}
}

func TestBS(t *testing.T) {
	bs := NewBoolSliceValue(nil, nil)
	f := setUpBSFlagSet(bs)

	vals := []string{"1", "F", "TRUE", "0"}
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range bs.Value {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}

		if actual := v.(*BoolValue).Value; b != actual {
			t.Fatalf("expected is[%d] to be %s but got: %t", i, vals[i], actual)
		}
	}
}

func TestBSDefault(t *testing.T) {
	bs := NewBoolSliceValue([]bool{false, true}, nil)
	f := setUpBSFlagSet(bs)

	vals := []string{"false", "T"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range bs.Value {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*BoolValue).Value; b != actual {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, actual)
		}
	}
}

func TestBSWithDefault(t *testing.T) {
	bs := NewBoolSliceValue([]bool{false, true}, nil)
	f := setUpBSFlagSet(bs)

	vals := []string{"FALSE", "1"}
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range bs.Value {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if actual := v.(*BoolValue).Value; b != actual {
			t.Fatalf("expected bs[%d] to be %t but got: %t", i, b, actual)
		}
	}
}

func TestBSCalledTwice(t *testing.T) {
	bs := NewBoolSliceValue(nil, nil)
	f := setUpBSFlagSet(bs)

	in := []string{"T,F", "T"}
	expected := []bool{true, false, true}
	argfmt := "--bs=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	for i, v := range bs.Value {
		if actual := v.(*BoolValue).Value; expected[i] != actual {
			t.Fatalf("expected bs[%d] to be %t but got %t", i, expected[i], actual)
		}
	}
}

func TestBSBadQuoting(t *testing.T) {
	tests := []struct {
		Want    []bool
		FlagArg []string
	}{
		{
			Want:    []bool{true, false, true},
			FlagArg: []string{"1", "0", "true"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"True", "F"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"T", "0"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"1", "0"},
		},
		{
			Want:    []bool{true, false, false},
			FlagArg: []string{"true,false", "false"},
		},
		{
			Want:    []bool{true, false, false, true, false, true, false},
			FlagArg: []string{"true,false,false,1,0,     T", " false "},
		},
		{
			Want:    []bool{false, false, true, false, true, false, true},
			FlagArg: []string{"0, False,  T,false  , true,F", "true"},
		},
	}

	for i, test := range tests {
		bs := NewBoolSliceValue(nil, nil)
		f := setUpBSFlagSet(bs)

		if err := f.Parse([]string{fmt.Sprintf("--bs=%s", strings.Join(test.FlagArg, ","))}); err != nil {
			t.Fatalf("flag parsing failed with error: %s\nparsing:\t%#v", err, test.FlagArg)
		}

		for j, b := range bs.Value {
			if actual := b.(*BoolValue).Value; actual != test.Want[j] {
				t.Fatalf("bad value parsed for test %d on bool %d:\nwant:\t%t\ngot:\t%t", i, j, test.Want[j], actual)
			}
		}
	}
}
