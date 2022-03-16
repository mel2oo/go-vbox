package main

import (
	"fmt"

	"github.com/mel2oo/go-vbox"
)

func main() {
	machines, err := vbox.ListMachines()
	if err != nil {
		return
	}

	for _, m := range machines {
		fmt.Println(m.Name)
	}
}
