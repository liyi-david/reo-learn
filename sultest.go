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
import "log"
import "os"
import "runtime"

var logger = log.New(os.Stderr, "TOP - ", 0)

func main() {
	// set up multicores
	runtime.GOMAXPROCS(4)
	s := GetOracle()
	// ---------------- CONFIGURATION ----------------------
	// sul.CloseLog()
	sul.CloseReoLog()
	sul.SetReoDelay(1)
	sul.SetBound(1)
	// -------------- CONFIGURATION END --------------------

	logger.Println("MAIN PROC START")
	// following are test code for MQuery
	/*
		for {
			var tin sul.InputSeq = sul.InputSeq{
				&sul.Input{map[string]bool{"A": true, "B": false}, false},
				&sul.Input{map[string]bool{"A": true, "B": true}, false},
				&sul.Input{map[string]bool{"A": false, "B": false}, false},
			}
			r := s.MQuery(tin)
			logger.Println("RESULT:", r)
			if r.String() == "ϵ" {
				break
			}
		}
	*/
	obs := learn.LStar(s)
	logger.Println(obs.GetHypo())
	logger.Println(sul.Counter())
	logger.Println(obs)
}
