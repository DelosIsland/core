package plock

import (
	"fmt"
	"sync"
	"testing"
	// "time"
)

func TestLockUnlock(t *testing.T) {
	want := 100
	a := 0
	pl := NewPriorityLock()
	wg := sync.WaitGroup{}
	for i := 0; i < want; i++ {
		go func(id int) {
			wg.Add(1)
			pl.Lock()
			a++
			// fmt.Println("I'm sleeping now: ", id)
			// time.Sleep(1 * time.Second)
			// fmt.Println("I've slept for 1s: ", id)
			pl.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println("a=", a)
	if a != want {
		t.Fail()
	}
}
