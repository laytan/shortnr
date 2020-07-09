package main

import "testing"

func TestOnePlusOne(t *testing.T) {
	answer := 1 + 1;
	if(answer != 2) {
		t.Error("One plus one is 2");
	}
}
