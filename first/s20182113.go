/*
Olivier JORIS -- s20182113

Why do I believe that my code is correct ?

I think that my code is correct because
the write channel ensures that either
a certain number of readers are reading
at the same time or only 1 writer.

Moreover, the read channel ensures that
only one reader can execute the read_enter
or read_exit function at the same time.

Finally, to ensure starvation a fair channel
is used in order that new readers do not block
a waiting writer.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randomSleep() {
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
}

func wait(q chan bool) { <-q }

func signal(q chan bool) { q <- true }

func semaphore(q chan bool, value int) {
	for {
		if value > 0 {
			select {
			case <-q:
				value++
			case q <- true:
				value--
			}
		} else {
			<-q
			value++
		}
	}
}

type RWlock struct {
	fair       chan bool
	read       chan bool
	write      chan bool
	read_count uint
}

func read_enter(lock *RWlock, id int) {
	wait(lock.fair)
	wait(lock.read)

	lock.read_count++
	if lock.read_count == 1 {
		wait(lock.write)
	}

	signal(lock.fair)
	signal(lock.read)
}

func read_exit(lock *RWlock, id int) {
	wait(lock.read)

	lock.read_count--
	if lock.read_count == 0 {
		signal(lock.write)
	}

	signal(lock.read)
}

func write_enter(lock *RWlock, id int) {
	wait(lock.fair)
	wait(lock.write)
	signal(lock.fair)
}

func write_exit(lock *RWlock, id int) {
	signal(lock.write)
}

func newRWlock(n int) *RWlock {
	size := 10
	lock := &RWlock{
		write:      make(chan bool, size),
		read:       make(chan bool, size),
		fair:       make(chan bool, size),
		read_count: 0}

	go lock.loop(n)

	return lock
}

func (lock *RWlock) loop(n int) {
	go semaphore(lock.read, 1)
	go semaphore(lock.fair, 1)
	go semaphore(lock.write, 1)
}

func proc(wg *sync.WaitGroup, id int, lock *RWlock) {
	defer wg.Done()
	for {
		switch rand.Intn(4) {
		case 0:
			randomSleep()
		case 1, 2:
			{
				read_enter(lock, id)
				fmt.Printf("%d enter read\n", id)
				randomSleep()
				fmt.Printf("%d exit read\n", id)
				read_exit(lock, id)
			}
		case 3:
			{
				write_enter(lock, id)
				fmt.Printf("%d enter write\n", id)
				randomSleep()
				fmt.Printf("%d exit write\n", id)
				write_exit(lock, id)
			}
		}
	}
}

func main() {
	var n = 10
	var lock = newRWlock(n)
	var w sync.WaitGroup
	w.Add(n)
	rand.Seed(42)

	for i := 0; i < n; i++ {
		go proc(&w, i, lock)
	}
	w.Wait()
}
