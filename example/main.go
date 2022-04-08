package main

import (
	"fmt"

	"github.com/mel2oo/go-vbox"
)

func main() {
	if err := vbox.Connect("root", "root", "10.20.53.139", 22); err != nil {
		return
	}

	machines, err := vbox.ListMachines()
	if err != nil {
		return
	}

	for _, m := range machines {
		fmt.Println(m.Name)
	}
}
