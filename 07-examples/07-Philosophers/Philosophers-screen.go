package main

import ("fmt"
  "time"
  "math/rand"
)

func randomSleep () {
  time.Sleep(time.Duration(50 + rand.Intn(100)) * time.Millisecond)
}

type Fork struct {
  id int
  Pick chan bool
  Drop chan bool
}

func (self * Fork) Init (id int) {
  self.id = id
  self.Pick = make(chan bool)
  self.Drop = make(chan bool)
}

func (self * Fork) Run () {
  for {
    <- self.Pick
    fmt.Printf("%d is picked\n", self.id)    
    <- self.Drop
    fmt.Printf("%d is dropped\n", self.id)    
  }
}

type Philosopher struct {
  id int
  first Fork
  second Fork
}

func (self * Philosopher) Init (id int, first Fork, second Fork) {
  self.id = id
  self.first = first
  self.second = second
}

func (self * Philosopher) Run () {
  for {
    fmt.Printf("%d is thinking\n", self.id)    
    self.first.Pick <- true
    self.second.Pick <- true
    fmt.Printf("%d is eating\n", self.id)
    self.first.Drop <- true
    self.second.Drop <- true
  }
}

func main () {
  const n = 3
  var fork [n] Fork;
  var philosopher [n] Philosopher;
  for i := 0; i < n; i++ { fork[i].Init(i); go fork[i].Run() }
  for i := 0; i < n; i++ {
    if i == 0 { philosopher[i].Init(i, fork[1], fork[0])
    } else { philosopher[i].Init(i, fork[i], fork[(i + 1) % n]) }
    go philosopher[i].Run() }
  time.Sleep(10 * time.Second)
}
