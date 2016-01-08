package alternator

import "../../lib/sul"
import "../../lib/reo"
import "time"

/*
	[ARCHITECTURE]
*/

func GetOracle() *sul.Oracle {
	o := new(sul.Oracle)
	o.InPorts = []string{"A", "B"}
	o.MidPorts = []string{"M0", "M1", "M2", "M3", "M4", "M5"}
	o.OutPorts = []string{"C"}
	o.TimeUnit = 40 * time.Millisecond
	o.GenerateInst = func() *sul.SulInst {
		r := new(sul.SulInst)
		// if there're several channels, a better solution is that
		// we use one stop flag to close all of them
		// and multiple stop finish flag to make sure that all of them
		// are closed
		r.GeneratePort(o)
		// generating stop ports
		r.StopPorts = reo.GenerateStopPort(7)
		r.Start = func() {
			go reo.BufferChannel(1, r.InPorts["A"], r.MidPorts["M0"], r.StopPorts[0])
			go reo.BufferChannel(1, r.InPorts["B"], r.MidPorts["M1"], r.StopPorts[1])
			go reo.OutputChannel(r.MidPorts["M5"], r.OutPorts["C"], r.StopPorts[2])
			go reo.FifoChannel(r.MidPorts["M0"], r.MidPorts["M2"], r.StopPorts[3])
			go reo.FifoChannel(r.MidPorts["M2"], r.MidPorts["M3"], r.StopPorts[4])
			go reo.SyncdrainChannel(r.MidPorts["M1"], r.MidPorts["M4"], r.StopPorts[5])
			go reo.ReplicatorChannel(r.MidPorts["M3"], r.MidPorts["M4"], r.MidPorts["M5"], r.StopPorts[6])
		}
		return r
	}
	return o
}
