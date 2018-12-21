package main

import (
	"log"
)

func main() {
	err := newFullMethodsGenerator().generate()
	if err != nil {
		log.Println(err)
	}
}
