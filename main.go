package main

import (
	_interface "download/interface"
	"log"
)

func main() {
	cli := _interface.NewCLI()

	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
