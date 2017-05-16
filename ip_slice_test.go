package pflag

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func setUpIPSFlagSet(ipsp *IPSliceValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.IPSliceVar(ipsp, "ips", "Command separated list!")
	return f
}

func setUpIPSFlagSetWithDefault() (*FlagSet, *IPSliceValue) {
	f := NewFlagSet("test", ContinueOnError)
	ipsp := NewIPSliceValue([]net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("0:0:0:0:0:0:0:1"),
	}, nil)

	f.IPSliceVar(ipsp, "ips", "Command separated list!")
	return f, ipsp
}

func TestEmptyIP(t *testing.T) {
	ips := NewIPSliceValue(nil, nil)
	f := setUpIPSFlagSet(ips)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(ips.Value) != 0 {
		t.Fatalf("got ips %v with len=%d but expected length=0", ips.Value, len(ips.Value))
	}
}

func TestIPS(t *testing.T) {
	ips := NewIPSliceValue(nil, nil)
	f := setUpIPSFlagSet(ips)

	vals := []string{"192.168.1.1", "10.0.0.1", "0:0:0:0:0:0:0:2"}
	arg := fmt.Sprintf("--ips=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips.Value {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v.(*IPValue).Value) {
			t.Fatalf("expected ips[%d] to be %s but got: %s from GetIPSlice", i, vals[i], v)
		}
	}
}

func TestIPSDefault(t *testing.T) {
	f, ips := setUpIPSFlagSetWithDefault()

	vals := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips.Value {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v.(*IPValue).Value) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		}
	}
}

func TestIPSWithDefault(t *testing.T) {
	f, ips := setUpIPSFlagSetWithDefault()

	vals := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
	arg := fmt.Sprintf("--ips=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips.Value {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v.(*IPValue).Value) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		}
	}
}

func TestIPSCalledTwice(t *testing.T) {
	ips := NewIPSliceValue(nil, nil)
	f := setUpIPSFlagSet(ips)

	in := []string{"192.168.1.2,0:0:0:0:0:0:0:1", "10.0.0.1"}
	expected := []net.IP{net.ParseIP("192.168.1.2"), net.ParseIP("0:0:0:0:0:0:0:1"), net.ParseIP("10.0.0.1")}
	argfmt := "ips=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips.Value {
		if !expected[i].Equal(v.(*IPValue).Value) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, expected[i], v)
		}
	}
}

func TestIPSBadQuoting(t *testing.T) {
	tests := []struct {
		Want    []net.IP
		FlagArg []string
	}{
		{
			Want: []net.IP{
				net.ParseIP("a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568"),
				net.ParseIP("203.107.49.208"),
				net.ParseIP("14.57.204.90"),
			},
			FlagArg: []string{
				"a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568",
				"203.107.49.208",
				"14.57.204.90",
			},
		},
		{
			Want: []net.IP{
				net.ParseIP("204.228.73.195"),
				net.ParseIP("86.141.15.94"),
			},
			FlagArg: []string{
				"204.228.73.195",
				"86.141.15.94",
			},
		},
		{
			Want: []net.IP{
				net.ParseIP("c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f"),
				net.ParseIP("4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472"),
			},
			FlagArg: []string{
				"c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f",
				"4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472",
			},
		},
		{
			Want: []net.IP{
				net.ParseIP("5170:f971:cfac:7be3:512a:af37:952c:bc33"),
				net.ParseIP("93.21.145.140"),
				net.ParseIP("2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca"),
			},
			FlagArg: []string{
				" 5170:f971:cfac:7be3:512a:af37:952c:bc33  , 93.21.145.140     ",
				"2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca",
			},
		},
		{
			Want: []net.IP{
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
			},
			FlagArg: []string{
				`2e5e:66b2:6441:848:5b74:76ea:574c:3a7b,        2e5e:66b2:6441:848:5b74:76ea:574c:3a7b,2e5e:66b2:6441:848:5b74:76ea:574c:3a7b     `,
				" 2e5e:66b2:6441:848:5b74:76ea:574c:3a7b",
			},
		},
	}

	for i, test := range tests {
		ips := NewIPSliceValue(nil, nil)
		f := setUpIPSFlagSet(ips)

		if err := f.Parse([]string{fmt.Sprintf("--ips=%s", strings.Join(test.FlagArg, ","))}); err != nil {

			t.Fatalf("flag parsing failed with error: %s", err)
		}

		for j, b := range ips.Value {
			actual := b.(*IPValue).Value
			if !actual.Equal(test.Want[j]) {
				t.Fatalf("bad value parsed for test %d on net.IP %d:\nwant:\t%s\ngot:\t%s", i, j, test.Want[j], actual)
			}
		}
	}
}
