package main

import "github.com/choria-io/choria-emulator/emulator"

func main() {
	err := emulator.Run()
	if err != nil {
		panic(err)
	}
}
