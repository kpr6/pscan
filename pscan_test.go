package main

import "testing"

func TestDummy(t *testing.T) {
	want := "hello"
	got := "hello"

	if got != want {
		t.Errorf("got %q want %q \n", got, want)
	}
}
