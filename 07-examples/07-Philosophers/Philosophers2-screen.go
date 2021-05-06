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
  for { <- self.Pick ; <- self.Drop }
}

type Waiter struct {
  n int
  sit_queue chan int
  leave_queue chan int
  ok [] chan bool
}

func (self * Waiter) Init (n int) {
  self.n = n
  self.sit_queue = make(chan int, n)
  self.leave_queue = make(chan int, n)
  self.ok = make([] chan bool, n)
  for i := 0; i < n; i++ { self.ok[i] = make(chan bool, 1) }
}

func (self * Waiter) Run () {
  var i int = 0
  var waiting bool = false
  var waiting_id int
  for {
    select {
    case v := <- self.sit_queue :
      if i < self.n - 1 { self.ok[v] <- true; i++
      } else { waiting = true; waiting_id = v }
    case <- self.leave_queue :
      if (waiting) { self.ok[waiting_id] <- true; waiting = false
      } else { i-- }
      default :
    }
  }
}

func (self * Waiter) sit (id int) { self.sit_queue <- id; <- self.ok[id] }

func (self * Waiter) leave (id int) { self.leave_queue <- id }

type Philosopher struct {
  id int
  first Fork
  second Fork
  waiter Waiter
}

func (self * Philosopher) Init (id int, first Fork, second Fork, waiter Waiter) {
  self.id = id
  self.first = first
  self.second = second
  self.waiter = waiter
}

func (self * Philosopher) Run () {
  for {
    fmt.Printf("%d is thinking\n", self.id)    
    randomSleep ()
    self.waiter.sit(self.id)
    self.first.Pick <- true
    self.second.Pick <- true
    fmt.Printf("%d is eating\n", self.id)
    randomSleep ()
    self.first.Drop <- true
    self.second.Drop <- true
    self.waiter.leave(self.id)
  }
}

func main () {
  const n = 5
  var fork [n] Fork;
  var philosopher [n] Philosopher;
  var waiter Waiter;
  waiter.Init(n)
  go waiter.Run()
  for i := 0; i < n; i++ { fork[i].Init(i); go fork[i].Run() }
  for i := 0; i < n; i++ {
    philosopher[i].Init(i, fork[i], fork[(i + 1) % n], waiter)
    go philosopher[i].Run() }
  time.Sleep(10 * time.Second)
}
