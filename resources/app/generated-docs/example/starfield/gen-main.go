package main

import (
	"fmt"

	"github.com/telecoda/pico-go-electron/console"
)

/*
	This is the generated bootstrap code for running a pico-go cartridge.

	Don't change any of this code, you shouldn't even really be seeing it, I guess you're the curious type.
*/

func NewCart() console.Cartridge {
	return &cartridge{
		BaseCartridge: console.NewBaseCart(),
	}
}

func main() {
	cart := NewCart()
	if err := console.Init(consoleType); err != nil {
		fmt.Printf("Failed to init console: %s\n", err)
		return
	}
	if err := console.Run(cart); err != nil {
		fmt.Printf("Failed to run cartridge: %s\n", err)
		return
	}
}
