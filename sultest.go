package main

import "./example/fifo"
import "./lib/learn"
import "fmt"

func main() {
	s := sulfifo.GetOracle()
	obs := learn.LStar(s)
	fmt.Println(obs)
}
