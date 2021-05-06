package main

import ("fmt"
  "strconv"
  "time"
)

func producer (q chan string) {
	for i := 0; ; i++ {
		q <- "Producer message " + strconv.Itoa(i)
	}
}

func consumer (q chan string) {
	for {
		msg := <- q
		fmt.Println("Consumer received: " + msg)
	}
}

func main () {
  var Size = 3
  var q = make(chan string, Size)

  go producer(q)
  go consumer(q)
  time.Sleep(2 * time.Second)
}
