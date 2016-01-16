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
import timer "./example/timer"
import expfifo "./example/expfifo"

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
	sul.SetReoDelay(10)
	sul.SetBound(2)
	sul.SetEquivBound(4)
	// logs on/off
	// sul.CloseLog()
	sul.CloseReoLog()
	// learn.CloseLog()
	// sul.ToggleTreeOptimization()
	// -------------- CONFIGURATION END --------------------
	var sulname = "expfifo"
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
	case "time":
		s = timer.GetOracle()
	case "expfifo":
		s = expfifo.GetOracle()
	}
	// -------------- ACTIVE LEARNING START --------------------
	// following are test code for MQuery
	var debug = false
	if debug {
		counter := 0
		for {
			var tin sul.InputSeq = sul.InputSeq{
				&sul.Input{map[string]bool{"A": true}, false},
				&sul.Input{map[string]bool{"A": false}, true},
				&sul.Input{map[string]bool{"A": false}, true},
				&sul.Input{map[string]bool{"A": false}, true},
			}
			r := s.SeqSimulate(tin)
			logger.Println("RESULT:", r, counter)
			counter++

			break
		}
		return
	}

	obs := learn.LStar(s)
	logger.Println(obs.GetHypoStr())
	logger.Println(sul.Counter())
	logger.Println("Time Cost in [SeqSimulate]:", sul.MembershipTime())
	logger.Println("Time Cost in [SeqRun]:", learn.RunTime())
}
