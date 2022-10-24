package vbox

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	ver, err := manage().Version()
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(ver)
}
