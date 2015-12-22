package main

// the fifo lib is imported with a prefix '.' indicating
// that we can use any functions in fifo lib with Func()
// instead of sulfifo.Func()
// since I'm going to write all the oracles with function
// GetOracle, this way of importing is more easy because
// we can change the system-under-learn without changing
// main code but the library we're going to import

import buf2 "./example/2-buffer"
import altn "./example/alternator"
import fifo "./example/fifo"

import "./lib/learn"
import "./lib/sul"
import "log"
import "os"
import "runtime"

var logger = log.New(os.Stderr, "TOP - ", 0)

func main() {
	// ---------------- CONFIGURATION ----------------------
	// set up multicores
	runtime.GOMAXPROCS(4)
	// configurations in simulation
	sul.SetReoDelay(5)
	sul.SetBound(1)
	// logs on/off
	sul.CloseLog()
	sul.CloseReoLog()
	// learn.CloseLog()
	// -------------- CONFIGURATION END --------------------
	var sulname = "buf2"
	var s *sul.Oracle
	switch sulname {
	case "buf2":
		s = buf2.GetOracle()
		break
	case "altn":
		s = altn.GetOracle()
		break
	case "fifo":
		s = fifo.GetOracle()
		break
	}
	// -------------- ACTIVE LEARNING START --------------------
	// following are test code for MQuery
	var debug = false
	if debug {
		counter := 0
		for {
			var tin sul.InputSeq = sul.InputSeq{
				&sul.Input{map[string]bool{"A": true, "B": false}, false},
			}
			r := s.MQuery(tin)
			logger.Println("RESULT:", r, counter)
			counter++
			if r.String() != "C:B," {
				break
			}
		}
		return
	}
	obs := learn.LStar(s)
	logger.Println(obs.GetHypoStr())
	logger.Println(sul.Counter())
}
