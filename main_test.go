package main

import "testing"

func TestMain(t *testing.T) {
	oldMain := mainFunc
	defer func() {
		mainFunc = oldMain
	}()

	called := false
	mainFunc = func() {
		called = true
	}

	main()

	if !called {
		t.Fail()
	}
}
