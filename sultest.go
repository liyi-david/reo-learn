package main

import "./example/fifo"
import "./lib/learn"
import "./lib/sul"
import "fmt"

func main() {
	s := sulfifo.GetOracle()
	obs := learn.LStar(s)
	fmt.Println(obs)
	fmt.Println(sul.Counter())
}
