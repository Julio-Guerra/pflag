package pflag

import (
	"fmt"
	"os"
	"testing"
)

func setUpIP(ip *IPValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.IPVar(ip, "address", "IP Address")
	return f
}

func TestIP(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		{"0.0.0.0", true, "0.0.0.0"},
		{" 0.0.0.0 ", true, "0.0.0.0"},
		{"1.2.3.4", true, "1.2.3.4"},
		{"127.0.0.1", true, "127.0.0.1"},
		{"255.255.255.255", true, "255.255.255.255"},
		{"", false, ""},
		{"0", false, ""},
		{"localhost", false, ""},
		{"0.0.0", false, ""},
		{"0.0.0.", false, ""},
		{"0.0.0.0.", false, ""},
		{"0.0.0.256", false, ""},
		{"0 . 0 . 0 . 0", false, ""},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		addr := NewIPValue(nil, nil)
		f := setUpIP(addr)

		tc := &testCases[i]

		arg := fmt.Sprintf("--address=%s", tc.input)
		err := f.Parse([]string{arg})
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure")
			continue
		} else if tc.success {
			if addr.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, addr.String())
			}
		}
	}
}
