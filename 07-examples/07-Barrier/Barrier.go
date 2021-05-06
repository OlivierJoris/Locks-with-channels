package main

import ("fmt"
  "time"
  "sync"
  "math/rand"
)

func randomSleep () {
  time.Sleep(time.Duration(50 + rand.Intn(100)) * time.Millisecond)
}

type barrier struct {
  in chan bool; out chan bool
}

func wait (b * barrier) { b.in <- true; <- b.out }

func newBarrier (n int) * barrier {
  b := & barrier {
    in : make(chan bool),
    out: make(chan bool) }
  go b.loop(n)
  return b
}

func (b * barrier) loop(n int) {
  for {
    for i := 0; i < n; i++ { <-b.in }
    for i := 0; i < n; i++ { b.out <- true }
  }
}

func proc (wg * sync.WaitGroup, id int, b * barrier) {
  defer wg.Done()
  for i := 30; i > 0; i-- {
    fmt.Printf("%d outside barrier\n", id)
    randomSleep()
    wait(b)
    fmt.Printf("%d inside barrier\n", id)
    randomSleep()
    wait(b)    
  }
}

func main () {
  var n = 3
  var b = newBarrier(n)
  var w sync.WaitGroup; w.Add(n)
  rand.Seed(42)

  for i := 0; i < n; i++ { go proc(&w, i, b) }
  w.Wait()
}
