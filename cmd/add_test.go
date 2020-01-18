package cmd

import (
	"testing"
)

func TestIsDot(t *testing.T) {
	result := isDot(".")
	if result == false {
		t.Fatal("failed test")
	}
}
