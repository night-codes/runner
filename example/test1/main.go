package main

import (
	"fmt"
	"time"
)

func main() {
	for now := range time.Tick(1 * time.Second) {
		fmt.Printf("%v\n", now)
	}
}
