package main

import ("fmt"
  "time"
  "sync"
  "math/rand"
)

func randomSleep () {
  time.Sleep(time.Duration(50 + rand.Intn(100)) * time.Millisecond)
}

...

type RWlock struct {
  ...
}

func read_enter (lock * RWlock, id int) {
  ...
}

func read_exit (lock * RWlock, id int) {
  ...
}

func write_enter (lock * RWlock, id int) {
  ...
}

func write_exit (lock * RWlock, id int) {
  ...
}

func newRWlock (n int) * RWlock {
  ...
  go lock.loop(n)
  return lock
}

func (lock * RWlock) loop(n int) {
  ...
}

func proc (wg * sync.WaitGroup, id int, lock * RWlock) {
  defer wg.Done()
  for {
    switch rand.Intn(4) {
    case 0: randomSleep()
    case 1, 2: {
      read_enter(lock, id)
      fmt.Printf("%d enter read\n", id)
      randomSleep()
      fmt.Printf("%d exit read\n", id)
      read_exit(lock, id)
    }
    case 3: {
      write_enter(lock, id)
      fmt.Printf("%d enter write\n", id)
      randomSleep()
      fmt.Printf("%d exit write\n", id)
      write_exit(lock, id)
    }
    }
  }
}

func main () {
  var n = 10
  var lock = newRWlock(n)
  var w sync.WaitGroup; w.Add(n)
  rand.Seed(42)

  for i := 0; i < n; i++ { go proc(&w, i, lock) }
  w.Wait()
}
