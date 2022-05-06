package go_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	fmt.Println(context.Background())
	fmt.Println(context.TODO())
}

func TestContextParentChild(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")
	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	fmt.Println(contextA.Value("b"))
	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextB.Value("e"))
}

func CreateCounterWithLeak() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()

	return destination
}

func TestGoroutineLeak(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	destination := CreateCounterWithLeak()

	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Goroutine 3

	// ini akan selalu listening ke goroutine di CreateCounter
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}

		/*
			Disini akan terjadi Goroutine Leak karena di perulangan ini
			ketika sudah mencapai 10 , perulangan akan berhenti ,
			sementara perulangan di dalam goroutine terus berjalan dengan berusaha memasukkan data
			ke dalam channel, namun tidak pernah ada yg menerima

			Untuk mencegah ini bisa menggunakan context
		*/
	}

	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Goroutine 3
	// Disini goroutine masih tetap 3 , padahal sudah tidak digunakan goroutine nya
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulasi slow
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(context.Background())
	destination := CreateCounter(ctx)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	// Cancel context (akan menghentikan goroutine berjalan)
	cancel()
	time.Sleep(4 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel() // Cancel context (akan menghentikan goroutine berjalan)

	destination := CreateCounter(ctx)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(4 * time.Second)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel() // Cancel context (akan menghentikan goroutine berjalan)

	destination := CreateCounter(ctx)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(4 * time.Second)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
