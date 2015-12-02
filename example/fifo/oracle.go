package sulfifo

import "../../lib/sul"
import "../../lib/reo"
import "time"

func GetOracle() *sul.Oracle {
	o := new(sul.Oracle)
	o.InPorts = []string{"A"}
	o.MidPorts = []string{"M0"}
	o.OutPorts = []string{"B"}
	o.TimeUnit = 100 * time.Millisecond
	o.GenerateInst = func() *sul.SulInst {
		r := new(sul.SulInst)
		r.GeneratePort(o)
		// if there're several channels, a better solution is that
		// we use one stop flag to close all of them
		// and multiple stop finish flag to make sure that all of them
		// are closed
		r.StopPorts = reo.GenerateStopPort(2)
		r.Start = func() {
			go reo.FifoChannel(r.InPorts["A"], r.MidPorts["M0"], r.StopPorts[0])
			go reo.OutputChannel(r.MidPorts["M0"], r.OutPorts["B"], r.StopPorts[1])
		}
		return r
	}
	return o
}
