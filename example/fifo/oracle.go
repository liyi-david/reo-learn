package sulfifo

import "../../lib/sul"
import "../../lib/reo"
import "time"

func GetOracle() *sul.Oracle {
	o := new(sul.Oracle)
	o.InPorts = []string{"A"}
	o.OutPorts = []string{"B"}
	o.TimeUnit = 100 * time.Millisecond
	o.GenerateInst = func() *sul.SutInst {
		r := new(sul.SutInst)
		r.InPorts = map[string]reo.Port{"A": reo.MakePort()}
		r.OutPorts = map[string]reo.Port{"B": reo.MakePort()}
		// if there're several channels, a better solution is that
		// we use one stop flag to close all of them
		// and multiple stop finish flag to make sure that all of them
		// are closed
		r.StopPorts = []reo.Port{}
		r.Start = func() {
			stopflag := make(chan string)
			stopport := reo.Port{stopflag, make(chan string)}
			go reo.FifoChannel(r.InPorts["A"], r.OutPorts["B"], stopport)
			r.StopPorts = append(r.StopPorts, stopport)
		}
		return r
	}
	return o
}
