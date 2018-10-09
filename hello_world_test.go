package main

import "testing"

func TestVerboseAdder(t *testing.T) {
	if verboseAdder(5,3) != 8 {
		t.Fail()
	}
}
