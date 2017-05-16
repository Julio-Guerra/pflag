// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	testBool, _                  = Bool("test_bool", "bool value")
	testInt                      = Int("test_int", "int value")
	testInt64                    = Int64("test_int64", "int64 value")
	testUint                     = UInt("test_uint", "uint value")
	testUint64                   = UInt64("test_uint64", "uint64 value")
	testString                   = String("test_string", "string value")
	testFloat                    = Float64("test_float64", "float64 value")
	testDuration                 = Duration("test_duration", "time.Duration value")
	testOptionalInt              = Int("test_optional_int", "optional int value")
	normalizeFlagNameInvocations = 0
)

func boolString(s string) string {
	if s == "0" || s == "<nil>" {
		return "false"
	}
	return "true"
}

func TestEverything(t *testing.T) {
	m := make(map[string]*Flag)
	visitor := func(expected string) func(f *Flag) {
		return func(f *Flag) {
			if len(f.Name) > 5 && f.Name[0:5] == "test_" {
				m[f.Name] = f
				ok := false
				switch {
				case f.Value.String() == expected:
					ok = true
				case f.Name == "test_string" && expected == "0" && f.Value.String() == "":
					ok = true
				case f.Name == "test_bool" && f.Value.String() == boolString(expected):
					ok = true
				case f.Name == "test_duration" && f.Value.String() == expected+"s":
					ok = true
				}
				if !ok {
					t.Errorf("Visit: bad value `%s` for %s", f.Value.String(), f.Name)
				}
			}
		}
	}
	VisitAll(visitor("0"))
	if len(m) != 9 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	m = make(map[string]*Flag)
	Visit(visitor("0"))
	if len(m) != 0 {
		t.Errorf("Visit sees unset flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now set all flags
	Set("test_bool", "true")
	Set("test_int", "1")
	Set("test_int64", "1")
	Set("test_uint", "1")
	Set("test_uint64", "1")
	Set("test_string", "1")
	Set("test_float64", "1")
	Set("test_duration", "1s")
	Set("test_optional_int", "1")
	Visit(visitor("1"))
	if len(m) != 9 {
		t.Error("Visit fails after set")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now test they're visited in sort order.
	var flagNames []string
	Visit(func(f *Flag) { flagNames = append(flagNames, f.Name) })
	if !sort.StringsAreSorted(flagNames) {
		t.Errorf("flag names not sorted: %v", flagNames)
	}
}

func TestUsage(t *testing.T) {
	called := false
	ResetForTesting(func() { called = true })
	if GetCommandLine().Parse([]string{"--x"}) == nil {
		t.Error("parse did not fail for unknown flag")
	}
	if !called {
		t.Error("did not call Usage for unknown flag")
	}
}

func TestAddFlagSet(t *testing.T) {
	oldSet := NewFlagSet("old", ContinueOnError)
	newSet := NewFlagSet("new", ContinueOnError)

	oldSet.String("flag1", "flag1")
	oldSet.String("flag2", "flag2")

	newSet.String("flag2", "flag2")
	newSet.String("flag3", "flag3")

	oldSet.AddFlagSet(newSet)

	if len(oldSet.formal) != 3 {
		t.Errorf("Unexpected result adding a FlagSet to a FlagSet %v", oldSet)
	}
}

func TestAnnotation(t *testing.T) {
	f := NewFlagSet("shorthand", ContinueOnError)

	if err := f.SetAnnotation("missing-flag", "key", nil); err == nil {
		t.Errorf("Expected error setting annotation on non-existent flag")
	}

	f.StringP("stringa", "a", "string value")
	if err := f.SetAnnotation("stringa", "key", nil); err != nil {
		t.Errorf("Unexpected error setting new nil annotation: %v", err)
	}
	if annotation := f.Lookup("stringa").Annotations["key"]; annotation != nil {
		t.Errorf("Unexpected annotation: %v", annotation)
	}

	f.StringP("stringb", "b", "string2 value")
	if err := f.SetAnnotation("stringb", "key", []string{"value1"}); err != nil {
		t.Errorf("Unexpected error setting new annotation: %v", err)
	}
	if annotation := f.Lookup("stringb").Annotations["key"]; !reflect.DeepEqual(annotation, []string{"value1"}) {
		t.Errorf("Unexpected annotation: %v", annotation)
	}

	if err := f.SetAnnotation("stringb", "key", []string{"value2"}); err != nil {
		t.Errorf("Unexpected error updating annotation: %v", err)
	}
	if annotation := f.Lookup("stringb").Annotations["key"]; !reflect.DeepEqual(annotation, []string{"value2"}) {
		t.Errorf("Unexpected annotation: %v", annotation)
	}
}

func testParse(f *FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolValue, _ := f.Bool("bool", "bool value")
	bool2Value, _ := f.Bool("bool2", "bool2 value")
	bool3Value, _ := f.Bool("bool3", "bool3 value")
	intValue := f.Int("int", "int value")
	int8Value := f.Int8("int8", "int value")
	int32Value := f.Int32("int32", "int value")
	int64Value := f.Int64("int64", "int64 value")
	uintValue := f.UInt("uint", "uint value")
	uint8Value := f.UInt8("uint8", "uint value")
	uint16Value := f.UInt16("uint16", "uint value")
	uint32Value := f.UInt32("uint32", "uint value")
	uint64Value := f.UInt64("uint64", "uint64 value")
	stringValue := f.String("string", "string value")
	float32Value := f.Float32("float32", "float32 value")
	float64Value := f.Float64("float64", "float64 value")
	ipValue := f.IP("ip", "ip value")
	maskValue := f.IPMask("mask", "mask value")
	durationValue := f.Duration("duration", "time.Duration value")
	optionalIntNoValue := NewIntValue(nil, 9)
	f.IntVar(optionalIntNoValue, "optional-int-no-value", "int value")
	optionalIntWithValue := NewIntValue(nil, 9)
	f.IntVar(optionalIntWithValue, "optional-int-with-value", "int value")
	extra := "one-extra-argument"
	args := []string{
		"--bool",
		"--bool2=true",
		"--bool3=false",
		"--int=22",
		"--int8=-8",
		"--int32=-32",
		"--int64=0x23",
		"--uint", "24",
		"--uint8=8",
		"--uint16=16",
		"--uint32=32",
		"--uint64=25",
		"--string=hello",
		"--float32=-172e12",
		"--float64=2718e28",
		"--ip=10.11.12.13",
		"--mask=255.255.255.0",
		"--duration=2m",
		"--optional-int-no-value",
		"--optional-int-with-value=42",
		extra,
	}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if boolValue.Value != true {
		t.Error("bool flag should be true, is ", boolValue.Value)
	}
	if bool2Value.Value != true {
		t.Error("bool2 flag should be true, is ", bool2Value.Value)
	}
	if bool3Value.Value != false {
		t.Error("bool3 flag should be false, is ", bool2Value.Value)
	}
	if intValue.Value != 22 {
		t.Error("int flag should be 22, is ", intValue.Value)
	}
	if int8Value.Value != -8 {
		t.Error("int8 flag should be 0x23, is ", int8Value.Value)
	}
	if int32Value.Value != -32 {
		t.Error("int32 flag should be 0x23, is ", int32Value.Value)
	}
	if int64Value.Value != 0x23 {
		t.Error("int64 flag should be 0x23, is ", int64Value.Value)
	}
	if uintValue.Value != 24 {
		t.Error("uint flag should be 24, is ", uintValue.Value)
	}
	if uint8Value.Value != 8 {
		t.Error("uint8 flag should be 8, is ", uint8Value.Value)
	}
	if uint16Value.Value != 16 {
		t.Error("uint16 flag should be 16, is ", uint16Value.Value)
	}
	if uint32Value.Value != 32 {
		t.Error("uint32 flag should be 32, is ", uint32Value.Value)
	}
	if uint64Value.Value != 25 {
		t.Error("uint64 flag should be 25, is ", uint64Value.Value)
	}
	if stringValue.Value != "hello" {
		t.Error("string flag should be `hello`, is ", stringValue.Value)
	}
	if float32Value.Value != -172e12 {
		t.Error("float32 value should be -172e12, is ", float32Value.Value)
	}
	if float64Value.Value != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", float64Value.Value)
	}
	if !ipValue.Value.Equal(net.ParseIP("10.11.12.13")) {
		t.Error("ip flag should be 10.11.12.13, is ", ipValue.Value)
	}
	if maskValue.String() != "255.255.255.0" {
		t.Error("mask flag should be 255.255.255.0, is ", maskValue.String())
	}
	if durationValue.Value != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", durationValue.Value)
	}
	if optionalIntNoValue.Value != 9 {
		t.Error("optional int flag should be the default value, is ", optionalIntNoValue.Value)
	}
	if optionalIntWithValue.Value != 42 {
		t.Error("optional int flag should be 42, is ", optionalIntWithValue.Value)
	}
	if len(f.Args()) != 1 {
		t.Error("expected one argument, got", len(f.Args()), f.Args())
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	}
}

func testParseAll(f *FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.StdBoolP("boola", "a", "bool value")
	f.BoolP("boolb", "b", "bool2 value")
	f.StdBoolP("boolc", "c", "bool3 value")
	f.BoolP("boold", "d", "bool4 value")
	f.StringP("strings", "s", "string value")
	f.StringP("stringz", "z", "string value")
	stringx := NewStringValue(nil, "1")
	f.StringVarP(stringx, "stringx", "x", "string value")
	f.StringP("stringy", "y", "string value")

	args := []string{
		"-ab",
		"-csxx",
		"--stringz=something",
		"-d=true",
		"-x",
		"-y",
		"ee",
	}
	want := []string{
		"boola", "true",
		"boolb", "true",
		"boolc", "true",
		"strings", "xx",
		"stringz", "something",
		"boold", "=true",
		"stringx", "1",
		"stringy", "ee",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.ParseAll() fail to restore the args")
		t.Errorf("Got: %v", got)
		t.Errorf("Want: %v", want)
	}
}

func TestShorthand(t *testing.T) {
	f := NewFlagSet("shorthand", ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolaFlag, _ := f.StdBoolP("boola", "a", "bool value")
	boolbFlag, _ := f.BoolP("boolb", "b", "bool2 value")
	boolcFlag, _ := f.StdBoolP("boolc", "c", "bool3 value")
	booldFlag, _ := f.BoolP("boold", "d", "bool4 value")
	stringaFlag := f.StringP("stringa", "s", "string value")
	stringzFlag := f.StringP("stringz", "z", "string value")
	extra := "interspersed-argument"
	notaflag := "--i-look-like-a-flag"
	args := []string{
		"-a",
		"-b", "true",
		extra,
		"-cs", "hello",
		"-z", "something",
		"-dtrue",
		"--",
		notaflag,
	}
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil {
		t.Error("expected no error, got ", err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if boolaFlag.Value != true {
		t.Error("boola flag should be true, is ", boolaFlag.Value)
	}
	if boolbFlag.Value != true {
		t.Error("boolb flag should be true, is ", boolbFlag.Value)
	}
	if boolcFlag.Value != true {
		t.Error("boolc flag should be true, is ", boolcFlag.Value)
	}
	if booldFlag.Value != true {
		t.Error("boold flag should be true, is ", booldFlag.Value)
	}
	if stringaFlag.Value != "hello" {
		t.Error("stringa flag should be `hello`, is ", stringaFlag.Value)
	}
	if stringzFlag.Value != "something" {
		t.Error("stringz flag should be `something`, is ", stringzFlag.Value)
	}
	if len(f.Args()) != 2 {
		t.Error("expected one argument, got", len(f.Args()))
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	} else if f.Args()[1] != notaflag {
		t.Errorf("expected argument %q got %q", notaflag, f.Args()[1])
	}
	if f.ArgsLenAtDash() != 1 {
		t.Errorf("expected argsLenAtDash %d got %d", f.ArgsLenAtDash(), 1)
	}
}

func TestShorthandLookup(t *testing.T) {
	f := NewFlagSet("shorthand", ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.StdBoolP("boola", "a", "bool value")
	f.BoolP("boolb", "b", "bool2 value")
	args := []string{
		"-ab",
	}
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil {
		t.Error("expected no error, got ", err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	flag := f.ShorthandLookup("a")
	if flag == nil {
		t.Errorf("f.ShorthandLookup(\"a\") returned nil")
	}
	if flag.Name != "boola" {
		t.Errorf("f.ShorthandLookup(\"a\") found %q instead of \"boola\"", flag.Name)
	}
	flag = f.ShorthandLookup("")
	if flag != nil {
		t.Errorf("f.ShorthandLookup(\"\") did not return nil")
	}
	defer func() {
		recover()
	}()
	flag = f.ShorthandLookup("ab")
	// should NEVER get here. lookup should panic. defer'd func should recover it.
	t.Errorf("f.ShorthandLookup(\"ab\") did not panic")
}

func TestParse(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParse(GetCommandLine(), t)
}

func TestParseAll(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParseAll(GetCommandLine(), t)
}

func TestFlagSetParse(t *testing.T) {
	testParse(NewFlagSet("test", ContinueOnError), t)
}

func TestPresentHelper(t *testing.T) {
	f := NewFlagSet("changedtest", ContinueOnError)
	f.Bool("changed", "changed bool")
	f.Bool("settrue", "true to true")
	f.Bool("setfalse", "false to false")
	f.Bool("unchanged", "unchanged bool")

	args := []string{"--changed", "--settrue", "--setfalse=false"}
	if err := f.Parse(args); err != nil {
		t.Error("f.Parse() error", err)
	}
	if !f.Present("changed") {
		t.Errorf("--changed wasn't changed!")
	}
	if !f.Present("settrue") {
		t.Errorf("--settrue wasn't changed!")
	}
	if !f.Present("setfalse") {
		t.Errorf("--setfalse wasn't changed!")
	}
	if f.Present("unchanged") {
		t.Errorf("--unchanged was changed!")
	}
	if f.ArgsLenAtDash() != -1 {
		t.Errorf("Expected argsLenAtDash: %d but got %d", -1, f.ArgsLenAtDash())
	}
}

func replaceSeparators(name string, from []string, to string) string {
	result := name
	for _, sep := range from {
		result = strings.Replace(result, sep, to, -1)
	}
	// Type convert to indicate normalization has been done.
	return result
}

func wordSepNormalizeFunc(f *FlagSet, name string) NormalizedName {
	seps := []string{"-", "_"}
	name = replaceSeparators(name, seps, ".")
	normalizeFlagNameInvocations++

	return NormalizedName(name)
}

func testWordSepNormalizedNames(args []string, t *testing.T) {
	f := NewFlagSet("normalized", ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	withDashFlag, _ := f.Bool("with-dash-flag", "bool value")
	// Set this after some flags have been added and before others.
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	withUnderFlag, _ := f.Bool("with_under_flag", "bool value")
	withBothFlag, _ := f.Bool("with-both_flag", "bool value")
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if withDashFlag.Value != true {
		t.Error("withDashFlag flag should be true, is ", withDashFlag.Value)
	}
	if withUnderFlag.Value != true {
		t.Error("withUnderFlag flag should be true, is ", withUnderFlag.Value)
	}
	if withBothFlag.Value != true {
		t.Error("withBothFlag flag should be true, is ", withBothFlag.Value)
	}
}

func TestWordSepNormalizedNames(t *testing.T) {
	args := []string{
		"--with-dash-flag",
		"--with-under-flag",
		"--with-both-flag",
	}
	testWordSepNormalizedNames(args, t)

	args = []string{
		"--with_dash_flag",
		"--with_under_flag",
		"--with_both_flag",
	}
	testWordSepNormalizedNames(args, t)

	args = []string{
		"--with-dash_flag",
		"--with-under_flag",
		"--with-both_flag",
	}
	testWordSepNormalizedNames(args, t)
}

func aliasAndWordSepFlagNames(f *FlagSet, name string) NormalizedName {
	seps := []string{"-", "_"}

	oldName := replaceSeparators("old-valid_flag", seps, ".")
	newName := replaceSeparators("valid-flag", seps, ".")

	name = replaceSeparators(name, seps, ".")
	switch name {
	case oldName:
		name = newName
		break
	}

	return NormalizedName(name)
}

func TestCustomNormalizedNames(t *testing.T) {
	f := NewFlagSet("normalized", ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}

	validFlag, _ := f.Bool("valid-flag", "bool value")
	f.SetNormalizeFunc(aliasAndWordSepFlagNames)
	someOtherFlag, _ := f.Bool("some-other-flag", "bool value")

	args := []string{"--old_valid_flag", "--some-other_flag"}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}

	if validFlag.Value != true {
		t.Errorf("validFlag is %v even though we set the alias --old_valid_falg", validFlag.Value)
	}
	if someOtherFlag.Value != true {
		t.Error("someOtherFlag should be true, is ", someOtherFlag.Value)
	}
}

// Every flag we add, the name (displayed also in usage) should normalized
func TestNormalizationFuncShouldChangeFlagName(t *testing.T) {
	// Test normalization after addition
	f := NewFlagSet("normalized", ContinueOnError)

	f.Bool("valid_flag", "bool value")
	if f.Lookup("valid_flag").Name != "valid_flag" {
		t.Error("The new flag should have the name 'valid_flag' instead of ", f.Lookup("valid_flag").Name)
	}

	f.SetNormalizeFunc(wordSepNormalizeFunc)
	if f.Lookup("valid_flag").Name != "valid.flag" {
		t.Error("The new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name)
	}

	// Test normalization before addition
	f = NewFlagSet("normalized", ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)

	f.Bool("valid_flag", "bool value")
	if f.Lookup("valid_flag").Name != "valid.flag" {
		t.Error("The new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name)
	}
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (_ *flagVar) Type() string {
	return "flagVar"
}

func (_ *flagVar) DefaultValue() string {
	return ""
}

func (_ *flagVar) DefaultArg() string {
	return ""
}

func TestUserDefined(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var v flagVar
	flags.VarP(&v, "v", "v", true, "usage")
	if err := flags.Parse([]string{"--v=1", "-v2", "-v", "3"}); err != nil {
		t.Error(err)
	}
	if len(v) != 3 {
		t.Fatal("expected 3 args; got ", len(v))
	}
	expect := "[1 2 3]"
	if v.String() != expect {
		t.Errorf("expected value %q got %q", expect, v.String())
	}
}

func TestSetOutput(t *testing.T) {
	var flags FlagSet
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.Init("test", ContinueOnError)
	flags.Parse([]string{"--unknown"})
	if out := buf.String(); !strings.Contains(out, "--unknown") {
		t.Logf("expected output mentioning unknown; got %q", out)
	}
}

// This tests that one can reset the flags. This still works but not well, and is
// superseded by FlagSet.
func TestChangingArgs(t *testing.T) {
	ResetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--before", "subcmd"}
	before, _ := StdBool("before", "")
	if err := CommandLine.Parse(os.Args[1:]); err != nil {
		t.Fatal(err)
	}
	cmd := Arg(0)
	os.Args = []string{"subcmd", "--after", "args"}
	after, _ := StdBool("after", "")
	if err := CommandLine.Parse(os.Args[1:]); err != nil {
		t.Fatal(err)
	}
	args := Args()

	if !before.Value || cmd != "subcmd" || !after.Value || len(args) != 1 || args[0] != "args" {
		t.Fatalf("expected true subcmd true [args] got %v %v %v %v", *before, cmd, *after, args)
	}
}

// Test that -help invokes the usage message and returns ErrHelp.
func TestHelp(t *testing.T) {
	var helpCalled = false
	fs := NewFlagSet("help test", ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	flag, _ := fs.Bool("flag", "regular flag")
	// Regular flag invocation should work
	err := fs.Parse([]string{"--flag=true"})
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if !flag.Value {
		t.Error("flag was not set by --flag")
	}
	if helpCalled {
		t.Error("help called for regular flag")
		helpCalled = false // reset for next test
	}
	// Help flag should work as expected.
	err = fs.Parse([]string{"--help"})
	if err == nil {
		t.Fatal("error expected")
	}
	if err != ErrHelp {
		t.Fatal("expected ErrHelp; got ", err)
	}
	if !helpCalled {
		t.Fatal("help was not called")
	}
	// If we define a help flag, that should override.
	fs.Bool("help", "help flag")
	helpCalled = false
	err = fs.Parse([]string{"--help"})
	if err != nil {
		t.Fatal("expected no error for defined --help; got ", err)
	}
	if helpCalled {
		t.Fatal("help was called; should not have been for defined help flag")
	}
}

func TestNoInterspersed(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetInterspersed(false)
	f.Bool("true", "always true")
	f.Bool("false", "always false")
	err := f.Parse([]string{"--true", "--", "break", "--false"})
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	args := f.Args()
	if len(args) != 2 || args[0] != "break" || args[1] != "--false" {
		t.Fatal("expected interspersed options/non-options to fail")
	}
}

func TestTermination(t *testing.T) {
	f := NewFlagSet("termination", ContinueOnError)
	boolFlag, _ := f.BoolP("bool", "l", "bool value")
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	arg1 := "ls"
	arg2 := "-l"
	args := []string{
		"--",
		arg1,
		arg2,
	}
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if boolFlag.Value {
		t.Error("expected boolFlag=false, got true")
	}
	if len(f.Args()) != 2 {
		t.Errorf("expected 2 arguments, got %d: %v", len(f.Args()), f.Args())
	}
	if f.Args()[0] != arg1 {
		t.Errorf("expected argument %q got %q", arg1, f.Args()[0])
	}
	if f.Args()[1] != arg2 {
		t.Errorf("expected argument %q got %q", arg2, f.Args()[1])
	}
	if f.ArgsLenAtDash() != 0 {
		t.Errorf("expected argsLenAtDash %d got %d", 0, f.ArgsLenAtDash())
	}
}

func TestDeprecatedFlagInDocs(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("badflag", "always true")
	f.MarkDeprecated("badflag", "use --good-flag instead")

	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "badflag") {
		t.Errorf("found deprecated flag in usage!")
	}
}

func TestDeprecatedFlagShorthandInDocs(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	name := "noshorthandflag"
	f.BoolP(name, "n", "always true")
	f.MarkShorthandDeprecated("noshorthandflag", fmt.Sprintf("use --%s instead", name))

	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "-n,") {
		t.Errorf("found deprecated flag shorthand in usage!")
	}
}

func parseReturnStderr(t *testing.T, f *FlagSet, args []string) (string, error) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := f.Parse(args)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stderr = oldStderr
	out := <-outC

	return out, err
}

func TestDeprecatedFlagUsage(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("badflag", "always true")
	usageMsg := "use --good-flag instead"
	f.MarkDeprecated("badflag", usageMsg)

	args := []string{"--badflag"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if !strings.Contains(out, usageMsg) {
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	}
}

func TestDeprecatedFlagShorthandUsage(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	name := "noshorthandflag"
	f.BoolP(name, "n", "always true")
	usageMsg := fmt.Sprintf("use --%s instead", name)
	f.MarkShorthandDeprecated(name, usageMsg)

	args := []string{"-n"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if !strings.Contains(out, usageMsg) {
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	}
}

func TestDeprecatedFlagUsageNormalized(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("bad-double_flag", "always true")
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	usageMsg := "use --good-flag instead"
	f.MarkDeprecated("bad_double-flag", usageMsg)

	args := []string{"--bad_double_flag"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if !strings.Contains(out, usageMsg) {
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	}
}

// Name normalization function should be called only once on flag addition
func TestMultipleNormalizeFlagNameInvocations(t *testing.T) {
	normalizeFlagNameInvocations = 0

	f := NewFlagSet("normalized", ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	f.Bool("with_under_flag", "bool value")

	if normalizeFlagNameInvocations != 1 {
		t.Fatal("Expected normalizeFlagNameInvocations to be 1; got ", normalizeFlagNameInvocations)
	}
}

//
func TestHiddenFlagInUsage(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("secretFlag", "shhh")
	f.MarkHidden("secretFlag")

	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "secretFlag") {
		t.Errorf("found hidden flag in usage!")
	}
}

//
func TestHiddenFlagUsage(t *testing.T) {
	f := NewFlagSet("bob", ContinueOnError)
	f.Bool("secretFlag", "shhh")
	f.MarkHidden("secretFlag")

	args := []string{"--secretFlag"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if strings.Contains(out, "shhh") {
		t.Errorf("usage message printed when using a hidden flag!")
	}
}

const defaultOutput = `      --A                         for bootstrapping, allow 'any' type (defaults to "false")
      --Alongflagname             disable bounds checking (defaults to "false")
  -C, --CCC                       a boolean defaulting to true (defaults to "true")
      --D=path                    set relative path for local imports
  -E, --EEE[=1234]                a num with NoOptDefVal (defaults to "4321")
      --F=number                  a non-zero number (defaults to "2.7")
      --G=float                   a float that defaults to zero
      --IP=ip                     IP address with no default
      --IPMask=ipmask             Netmask address with no default
      --IPNet=ipNetwork           IP network with no default
      --Ints=intSlice             int slice with zero default
      --N=int                     a non-zero int (defaults to "27")
      --ND1[="bar"]               a string with default values (defaults to "foo")
      --ND2[=4321]                a num with default values (defaults to "1234")
      --StringSlice=stringSlice   string slice with zero default
      --Z=int                     an int that defaults to zero
      --custom=custom             custom Value implementation
      --customP=custom            a VarP with default (defaults to "10")
      --maxT=timeout              set timeout for dial
`

// Custom value that satisfies the Value interface.
type customValue int

func (cv *customValue) String() string { return fmt.Sprintf("%v", *cv) }

func (cv *customValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*cv = customValue(v)
	return err
}

func (_ *customValue) Type() string { return "custom" }

func (_ *customValue) DefaultArg() string { return "" }

func (c *customValue) DefaultValue() string {
	if *c == 0 {
		return ""
	}
	return strconv.FormatInt(int64(*c), 10)
}

func TestPrintDefaults(t *testing.T) {
	fs := NewFlagSet("print defaults test", ContinueOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Bool("A", "for bootstrapping, allow 'any' type")
	fs.Bool("Alongflagname", "disable bounds checking")
	ccc := NewBoolValue(true, true)
	fs.BoolVarP(ccc, "CCC", "C", "a boolean defaulting to true")
	fs.String("D", "set relative `path` for local imports")
	f := NewFloat64Value(2.7, nil)
	fs.Float64Var(f, "F", "a non-zero `number`")
	fs.Float64("G", "a float that defaults to zero")
	n := NewIntValue(27, nil)
	fs.IntVar(n, "N", "a non-zero int")
	fs.IntSlice("Ints", "int slice with zero default")
	fs.IP("IP", "IP address with no default")
	fs.IPMask("IPMask", "Netmask address with no default")
	fs.IPNet("IPNet", "IP network with no default")
	fs.Int("Z", "an int that defaults to zero")
	fs.Duration("maxT", "set `timeout` for dial")
	fs.StringVar(NewStringValue("foo", "bar"), "ND1", "a string with default values")
	fs.IntVar(NewIntValue(1234, 4321), "ND2", "a `num` with default values")
	fs.IntVarP(NewIntValue(4321, 1234), "EEE", "E", "a `num` with NoOptDefVal")
	fs.StringSlice("StringSlice", "string slice with zero default")

	var cv customValue
	fs.Var(&cv, "custom", true, "custom Value implementation")

	cv2 := customValue(10)
	fs.VarP(&cv2, "customP", "", true, "a VarP with default")

	fs.PrintDefaults()
	got := buf.String()
	if got != defaultOutput {
		fmt.Println("\n" + got)
		fmt.Println("\n" + defaultOutput)
		t.Errorf("got %q want %q\n", got, defaultOutput)
	}
}

func TestVisitAllFlagOrder(t *testing.T) {
	fs := NewFlagSet("TestVisitAllFlagOrder", ContinueOnError)
	fs.SortFlags = false
	// https://github.com/spf13/pflag/issues/120
	fs.SetNormalizeFunc(func(f *FlagSet, name string) NormalizedName {
		return NormalizedName(name)
	})

	names := []string{"C", "B", "A", "D"}
	for _, name := range names {
		fs.Bool(name, "")
	}

	i := 0
	fs.VisitAll(func(f *Flag) {
		if names[i] != f.Name {
			t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
		}
		i++
	})
}

func TestVisitFlagOrder(t *testing.T) {
	fs := NewFlagSet("TestVisitFlagOrder", ContinueOnError)
	fs.SortFlags = false
	names := []string{"C", "B", "A", "D"}
	for _, name := range names {
		fs.Bool(name, "")
		fs.Set(name, "true")
	}

	i := 0
	fs.Visit(func(f *Flag) {
		if names[i] != f.Name {
			t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
		}
		i++
	})
}
