package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

func genData(writer io.Writer, timeout time.Duration) {
	promise := make(chan bool)
	ch := make(chan int)

	const nRoutines = 10
	const dataChunksCount = 20

	for i := 0; i < nRoutines; i++ {
		go func(routineNo int) {
			data := ""
			for i := 0; i < dataChunksCount; i++ {
				data += fmt.Sprintf("black_overlord_%d_magic_number_%d:", routineNo, i)
			}

			fmt.Printf("Routine %d: data prepare done\n%s\n", routineNo, data)

			if timeout > 0 {
				fmt.Printf("Routing %d: sleeping %s\n", routineNo, timeout)
				time.Sleep(timeout)
			}

			ch <- routineNo
			<-promise

			fmt.Fprintf(writer, "%s", data)

			ch <- routineNo
		}(i)
	}

	waitForRoutines(nRoutines, ch, "data preparing")
	close(promise)
	waitForRoutines(nRoutines, ch, "data priting")
}

func waitForRoutines(nRoutines int, ch chan int, ctx string) {
	for i := 0; i < nRoutines; i++ {
		routineNo := <-ch
		fmt.Printf("Waiter %q: routine %d done; waiting %d/%d routines...\n",
			ctx, routineNo, nRoutines-i-1, nRoutines)
	}
}

type localGenAndPrint interface {
	io.Writer
	io.WriterTo
}

func genAndPrint(writer localGenAndPrint) {
	genData(writer, 0)

	fmt.Printf("\n\n\nTOTAL DATA\n\n\n")
	writer.WriteTo(os.Stdout)
}

func func1() {
	genAndPrint(&bytes.Buffer{})
}

type localReadWriter struct {
	buf bytes.Buffer
	mu  sync.RWMutex
}

func (rw *localReadWriter) Read(dest []byte) (int, error) {
	rw.mu.RLock()
	defer rw.mu.RUnlock()

	return rw.buf.Read(dest)
}

func (rw *localReadWriter) Write(p []byte) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	return rw.buf.Write(p)
}

func (rw *localReadWriter) WriteTo(dest io.Writer) (int64, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	return rw.buf.WriteTo(dest)
}

func func2() {
	genAndPrint(&localReadWriter{})
}

func randDuration() time.Duration {
	return time.Duration(rand.Int31n(1000)) * time.Millisecond
}

func func3() {
	promise := make(chan struct{})

	go func() {
		genData(&localReadWriter{}, randDuration())

		promise <- struct{}{}
	}()

	select {
	case <-promise:
		fmt.Println("==============================================================")
		fmt.Println("GEN DATA DONE!")
		fmt.Println("==============================================================")
	case <-time.Tick(400 * time.Millisecond):
		fmt.Println("==============================================================")
		fmt.Println("TIMEOUT!")
		fmt.Println("==============================================================")
	}
}

func func4() {
	var wg sync.WaitGroup

	const nRoutines = 10
	for i := 0; i < nRoutines; i++ {
		wg.Add(1)

		go func(routineNo int) {
			time.Sleep(randDuration())
			fmt.Println("Routine ", routineNo, " done")
			wg.Done()
		}(i)
	}

	fmt.Println("Waiting for routines")
	wg.Wait()
	fmt.Println("Waiting done")
}

func main() {
	rand.Seed(time.Now().Unix())

	func1()
	//func2()
	//func3()
	//func4()
}
