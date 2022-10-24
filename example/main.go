package main

import (
	"fmt"

	"github.com/mel2oo/go-vbox/vbox"
)

func main() {
	mgr, err := vbox.NewManage()
	if err != nil {
		return
	}

	machines, err := mgr.ListMachines()
	if err != nil {
		return
	}

	for _, m := range machines {
		fmt.Println(m.Name)
		if m.Name == "win7_64_1" {
			m.Start()
		}
	}
}
