package main

import (
	"fmt"

	"github.com/mel2oo/go-vbox"
)

func main() {
	if _, err := vbox.NewSSHCmd("root", "Dbapp@2121", "10.20.152.15", 22); err != nil {
		return
	}

	machines, err := vbox.ListMachines()
	if err != nil {
		return
	}

	for _, m := range machines {
		fmt.Println(m.Name)
		if m.Name == "win7_64_1" {
			m.Start("11")
		}
	}
}
