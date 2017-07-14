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

	s := make([]byte, mem, mem)
	fmt.Printf("Allocated %d\n", len(s))
	os.Exit(0)
}
