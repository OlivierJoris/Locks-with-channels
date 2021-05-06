package main

import ("fmt"
  "strconv"
  "time"
  "sync"
)

func producer (wg * sync.WaitGroup, q chan string, id int) {
  defer wg.Done()
	for i := 0; i < 10; i++ {
		q <- "Message " + strconv.Itoa(i) + " from " + strconv.Itoa(id)
    time.Sleep(100 * time.Millisecond)
	}
}

func consumer (wg * sync.WaitGroup, q chan string, id int) {
  defer wg.Done()
	for {
		msg, more := <- q
    if !more { break }
		fmt.Println(strconv.Itoa(id) + " received: " + msg)
	}
}

func main () {
  var nP, nC = 2, 4
  var Size = 3
  var q = make(chan string, Size)
  var wP, wC sync.WaitGroup
  
  wP.Add(nP)
  wC.Add(nC)
  for i := 1; i <= nP; i++ { go producer(&wP, q, i) }
  for i := 1; i <= nC; i++ { go consumer(&wC, q, i) }
  wP.Wait()
  close(q)
  wC.Wait()
}
