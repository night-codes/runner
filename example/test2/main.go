package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Test 2 started")
	go func() {
		for range time.Tick(5 * time.Second) {
			fmt.Printf("\n")
		}
	}()
	go func() {
		for range time.Tick(15 * time.Second) {
			os.Exit(0)
		}
	}()

	for range time.Tick(1 * time.Second) {
		fmt.Printf("+_+_+_+_+_+_+_ ")
	}
}
