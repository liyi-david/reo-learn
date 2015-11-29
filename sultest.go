package main

// the fifo lib is imported with a prefix '.' indicating
// that we can use any functions in fifo lib with Func()
// instead of sulfifo.Func()
// since I'm going to write all the oracles with function
// GetOracle, this way of importing is more easy because
// we can change the system-under-learn without changing
// main code but the library we're going to import
import . "./example/alternator"

// import . "./example/fifo"
import "./lib/learn"
import "./lib/sul"
import "fmt"

func main() {
	s := GetOracle()
	/*
		var tin sul.InputSeq = sul.InputSeq{
			&sul.Input{map[string]bool{"A": true, "B": true}, false},
			&sul.Input{map[string]bool{"A": true, "B": true}, false},
		}
		fmt.Println(s.MQuery(tin))
	*/
	obs := learn.LStar(s)
	fmt.Println(obs.GetHypo())
	fmt.Println(sul.Counter())
	fmt.Println(obs)
}
