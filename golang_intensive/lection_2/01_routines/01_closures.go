package main

import (
	"fmt"
	"math/rand"
	"time"
)

func func1() {
	const nRoutines = 10
	ch := make(chan int, nRoutines)

	for i := 0; i < nRoutines; i++ {
		go func() {
			sleepFor := rand.Int31n(1000)
			fmt.Printf("GoRoutine %d: sleeping for %d ms\n", i, sleepFor)

			time.Sleep(time.Duration(sleepFor) * time.Millisecond)
			fmt.Printf("GoRoutine %d: done\n", i)

			ch <- i
		}()
	}

	for i := 0; i < nRoutines; i++ {
		fmt.Println("main: waiting for goroutine...")
		fmt.Printf("main: goroutine %d done!\n", <-ch)
	}

	fmt.Println("main: done")
}

func main() {
	func1()
}
