package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	mem, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("bad arg: " + err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Could not allocate array: %v\n", r)
			os.Exit(1)
		}
	}()

	s := make([]byte, mem, mem)
	fmt.Printf("Allocated %d\n", len(s))
	os.Exit(0)
}
