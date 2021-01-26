package locke_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/blukai/locke"
)

func Example() {
	locke := locke.New()
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()

		locker := locke.NewTxn("xd")

		fmt.Println("sypha:  request lock")
		locker.Lock()
		fmt.Println("sypha:  locked")

		fmt.Println("sypha:  sleep")
		time.Sleep(time.Millisecond)

		fmt.Println("sypha:  request unlock")
		locker.Unlock()
		fmt.Println("sypha:  unlocked")
	}()

	// hacky way to make sure that sypha runs first?
	time.Sleep(time.Microsecond * 10)

	wg.Add(1)
	go func() {
		defer wg.Done()

		locker := locke.NewTxn("xd")

		fmt.Println("trevor: request lock")
		locker.Lock()
		fmt.Println("trevor: locked")

		fmt.Println("trevor: request unlock")
		locker.Unlock()
		fmt.Println("trevor: unlocked")
	}()

	wg.Wait()

	// Output:
	// sypha:  request lock
	// sypha:  locked
	// sypha:  sleep
	// trevor: request lock
	// sypha:  request unlock
	// sypha:  unlocked
	// trevor: locked
	// trevor: request unlock
	// trevor: unlocked
}

func ExampleDiningPhilosophers() {
	locke := locke.New()
	wg := new(sync.WaitGroup)

	for n := 0; n < 5; n++ {
		// each of the five philosophers lives in their own goroutine
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			cl := n
			cr := (n + 5 - 1) % 5

			locker := locke.NewTxn(cl, cr)

			fmt.Printf("%d is thinking\n", n)
			time.Sleep(20 * time.Millisecond)
			locker.Lock()

			fmt.Printf("%d is eating with chopsticks %d and %d\n", n, cl, cr)
			time.Sleep(10 * time.Millisecond)
			locker.Unlock()
		}(n)
	}

	wg.Wait()

	// Unordered utput:
	// 4 is thinking
	// 0 is thinking
	// 3 is thinking
	// 2 is thinking
	// 1 is thinking
	// 2 is eating with chopsticks 2 and 1
	// 4 is eating with chopsticks 4 and 3
	// 1 is eating with chopsticks 1 and 0
	// 3 is eating with chopsticks 3 and 2
	// 0 is eating with chopsticks 0 and 4
}

func BenchmarkDiningPhilosophersX5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		locke := locke.New()
		wg := new(sync.WaitGroup)

		for n := 0; n < 5; n++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				left := n
				right := (n + 5 - 1) % 5
				locker := locke.NewTxn(left, right)
				for i := 0; i < 5; i++ {
					time.Sleep(20 * time.Millisecond)
					locker.Lock()

					time.Sleep(10 * time.Millisecond)
					locker.Unlock()
				}
			}(n)
		}
		wg.Wait()
	}
}
