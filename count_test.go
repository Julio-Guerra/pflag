package pflag

import (
	"os"
	"testing"
)

func setUpCount(c *CountValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.CountVarP(c, "verbose", "v", "a counter")
	return f
}

func TestCount(t *testing.T) {
	testCases := []struct {
		input    []string
		success  bool
		expected int
	}{
		{[]string{"-v", "-v", "-v"}, true, 3},
		{[]string{"-v", "--verbose", "-v"}, true, 3},
		{[]string{"-v3", "-v"}, true, 4},
		{[]string{"-va"}, false, 0},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		count := NewStandardCountValue()
		f := setUpCount(count)

		tc := &testCases[i]

		err := f.Parse(tc.input)
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure, got success")
			continue
		} else if tc.success {
			if count.Value != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, count.Value)
			}
		}
	}
}
