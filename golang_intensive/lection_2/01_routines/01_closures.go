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

func sleepInRoutine() {
	for i := 0; i < 10; i++ {
		go func(routineNo int) {
			fmt.Println("Trying to sleep in goroutine ", routineNo)
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Sleep done in goroutine ", routineNo)
		}(i)
	}
}

func func2() {
	fmt.Println("Testing goroutines scope")
	sleepInRoutine()
	time.Sleep(2 * time.Second)
	fmt.Println("Main function done")
}

type T struct {
	A int
}

func (t T) fun() {
	fmt.Println("Into T: A == ", t.A)
}

func (t *T) fun2() {
	fmt.Println("Into T.func2: A == ", t.A)
	t.A++
}

func func3() {
	for _, t := range []T{T{1}, T{2}, T{3}} {
		//go t.fun2()
		go t.fun()
		/*
			go func() {
				time.Sleep(10 * time.Millisecond)
				t.fun()
			}()
		*/
	}
	time.Sleep(1 * time.Second)
}

func main() {
	func1()
	//func2()
	//func3()
}
