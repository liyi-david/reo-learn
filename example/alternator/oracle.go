package alternator

import "../../lib/sul"
import "../../lib/reo"
import "time"

/*
	[ARCHITECTURE]
[M6]A ---> M0 -- [Sync] --> M4 ---> C
     |                        |
     M1                 -->  M5
     |                 /
[SyncDrain]           /
     |          [FIFO]
     M2          /
     |          /
[M7]B --------> M3
*/

func GetOracle() *sul.Oracle {
	o := new(sul.Oracle)
	o.InPorts = []string{"A", "B"}
	o.MidPorts = []string{"M0", "M1", "M2", "M3", "M4", "M5", "M6", "M7", "M8"}
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
		r.StopPorts = reo.GenerateStopPort(9)
		r.Start = func() {
			go reo.ReplicatorChannel(r.MidPorts["M6"], r.MidPorts["M0"], r.MidPorts["M1"], r.StopPorts[0])
			go reo.ReplicatorChannel(r.MidPorts["M7"], r.MidPorts["M2"], r.MidPorts["M3"], r.StopPorts[1])
			go reo.MergerChannel(r.MidPorts["M4"], r.MidPorts["M5"], r.MidPorts["M8"], r.StopPorts[2])
			go reo.SyncdrainChannel(r.MidPorts["M1"], r.MidPorts["M2"], r.StopPorts[3])
			go reo.SyncChannel(r.MidPorts["M0"], r.MidPorts["M4"], r.StopPorts[4])
			go reo.FifoChannel(r.MidPorts["M3"], r.MidPorts["M5"], r.StopPorts[5])
			go reo.BufferChannel(1, r.InPorts["A"], r.MidPorts["M6"], r.StopPorts[6])
			go reo.BufferChannel(1, r.InPorts["B"], r.MidPorts["M7"], r.StopPorts[7])
			go reo.OutputChannel(r.MidPorts["M8"], r.OutPorts["C"], r.StopPorts[8])
		}
		return r
	}
	return o
}
