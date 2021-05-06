package main

import ("fmt"
  "time"
  "sync"
)

func wait (q chan bool) { <- q }

func signal (q chan bool) { q <- true }

func semaphore (q chan bool, value int) {
  for {
    if value > 0 {
      select {
      case <- q: value++
      case q <- true: value--
      }
    } else {
      <- q
      value++
    }
  }
}

func proc (wg * sync.WaitGroup, id int, q chan bool) {
  defer wg.Done()
  i := 30
  for ; i > 0; i-- {
    time.Sleep(100 * time.Millisecond)
    wait(q)
    fmt.Printf("Proc %d entering critical section\n", id)
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Proc %d exiting critical section\n", id)
    signal(q)
  }
}

func main () {
  var n = 3
  var Size = 0
  var q = make(chan bool, Size)
  var w sync.WaitGroup

  go semaphore(q, 1)
  
  w.Add(n)
  for i := 1; i <= n; i++ {
    go proc(&w, i, q) }
  w.Wait()
}
