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
  in [] chan bool; out [] chan bool
}

func wait (b * barrier, id int) {
  b.in[id] <- true; <- b.out[id]
}

func newBarrier (n int) * barrier {
  Size := 10
  tmp_in := make([] chan bool, n)
  tmp_out := make([] chan bool, n)
  for i := 0; i < n; i++ {
    tmp_in[i]  = make(chan bool, Size)
    tmp_out[i] = make(chan bool, Size)
  }
  b := & barrier {
    in : tmp_in, out: tmp_out }
  go b.loop(n)
  return b
}

func (b * barrier) loop(n int) {
  for {
    for id := 0; id < n; id++ { <- b.in[id] }
    for id := 0; id < n; id++ { b.out[id] <- true }
  }
}

func proc (wg * sync.WaitGroup, id int, b * barrier) {
  defer wg.Done()
  for i := 30; i > 0; i-- {
    fmt.Printf("%d outside barrier\n", id)
    randomSleep()
    wait(b, id)
    fmt.Printf("%d inside barrier\n", id)
    randomSleep()
    wait(b, id)
  }
}

func main () {
  var n = 3
  var b = newBarrier(n)
  var w sync.WaitGroup; w.Add(n)
  rand.Seed(42)

  for i := 0; i < n; i++ {
    go proc(&w, i, b) }
  w.Wait()
}
