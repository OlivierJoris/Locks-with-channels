/* Olivier JORIS -- 20182113

QUESTION

Why do I believe my code is correct ?

I think that my code is correct because when a writer enters the write_enter
function, a message is sent through the writer_in channel to the synchronizer
process (loop function) then, this process wait that all current readers
finish their execution through the reader_out channel and the reading map.
Then, using the ok_writer, the synchronizer informs the reader it can
start its execution (there is no current reader). The synchronizer is now
waiting for the end of the execution of the writer through the writer_end
channel. Only one writer can thus write at the same time without any reader.

If the synchronizer receives a read request (through the reader_in channel),
it means that there is no current writer. The reader_ok channel is used to
ensure that the map has been updated according to this read request.
Several readers can thus read at the same time without any writer.

Moreover, all previous mentionned channels, using message passing
synchronization, are ensuring that this code is thread safe and data race free.

Finally, because, in go, a select on several non-empty channel is performed
randomly, this code is starvation free : a new read request is not blocking
a current waiting write request and vice versa.
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

type RWlock struct {
	reader_in  chan int
	reader_out chan int
	writer_in  chan int
	writer_out chan int
	ok_reader  []chan bool
	ok_writer  []chan bool
}

func read_enter(lock *RWlock, id int) {
	lock.reader_in <- id
	<-lock.ok_reader[id]
}

func read_exit(lock *RWlock, id int) {
	lock.reader_out <- id
}

func write_enter(lock *RWlock, id int) {
	lock.writer_in <- id
	<-lock.ok_writer[id]
}

func write_exit(lock *RWlock, id int) {
	lock.writer_out <- id
}

func newRWlock(n int) *RWlock {
	size := 10

	tmp_ok_writer := make([]chan bool, n)
	tmp_ok_reader := make([]chan bool, n)

	for i := 0; i < n; i++ {
		tmp_ok_writer[i] = make(chan bool, size)
		tmp_ok_reader[i] = make(chan bool, size)
	}

	lock := &RWlock{
		reader_in:  make(chan int, n),
		reader_out: make(chan int, n),
		writer_in:  make(chan int, n),
		writer_out: make(chan int, n),
		ok_writer:  tmp_ok_writer,
		ok_reader:  tmp_ok_reader}

	go lock.loop(n)

	return lock
}

func (lock *RWlock) loop(n int) {
	// Map maintaining current readers
	reading := make(map[int]bool)

	for {
		select {
		// Write request
		case v1 := <-lock.writer_in:
			// Wait for current readers to finish their execution
			for id := range reading {
				<-lock.reader_out
				delete(reading, id)
			}
			// Informs that writer can start (no reader)
			lock.ok_writer[v1] <- true
			// Wait that the writer finish his execution
			<-lock.writer_out
		// Read request
		case v2 := <-lock.reader_in:
			reading[v2] = true
			// Informs that reader can start (map updated)
			lock.ok_reader[v2] <- true
		// Reader end
		case v3 := <-lock.reader_out:
			delete(reading, v3)
		default:
		}
	}
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
