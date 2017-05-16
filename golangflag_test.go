// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"testing"
)

func TestGoflags(t *testing.T) {
	stringValue := goflag.String("stringFlag", "stringFlag", "stringFlag")
	boolValue := goflag.Bool("boolFlag", false, "boolFlag")

	f := NewFlagSet("test", ContinueOnError)

	f.AddGoFlagSet(goflag.CommandLine)
	err := f.Parse([]string{"--stringFlag=bob", "--boolFlag"})
	if err != nil {
		t.Fatal("expected no error; get", err)
	}

	if *stringValue != "bob" {
		t.Fatalf("expected getString=bob but got getString=%s", *stringValue)
	}

	if *boolValue != true {
		t.Fatalf("expected getBool=true but got getBool=%v", *boolValue)
	}
}
