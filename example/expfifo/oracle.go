package expfifo

import "../../lib/sul"
import "../../lib/reo"
import "time"
import "strconv"

func GetOracle() *sul.Oracle {
	o := new(sul.Oracle)
	o.InPorts = []string{"A"}
	// dynamically generate the MidPorts
	o.MidPorts = []string{}
	for i := 0; i < 14; i++ {
		o.MidPorts = append(o.MidPorts, "M"+strconv.Itoa(i))
	}

	o.OutPorts = []string{"B"}
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
			// Input & Output
			go reo.BufferChannel(1, r.InPorts["A"], r.MidPorts["M0"], r.StopPorts[0])
			go reo.OutputChannel(r.MidPorts["M12"], r.OutPorts["B"], r.StopPorts[1])
			// Function Channels
			go reo.SyncChannel(r.MidPorts["M0"], r.MidPorts["M1"], r.StopPorts[2])
			go reo.ReplicatorChannel(r.MidPorts["M1"], r.MidPorts["M2"], r.MidPorts["M3"], r.StopPorts[3])
			go reo.FifoChannel(r.MidPorts["M2"], r.MidPorts["M4"], r.StopPorts[4])
			go reo.TimerChannel(r.MidPorts["M3"], r.MidPorts["M5"], 40*time.Millisecond, r.StopPorts[5])
			go reo.LossysyncChannel(r.MidPorts["M5"], r.MidPorts["M6"], r.StopPorts[6])
			go reo.ReplicatorChannel(r.MidPorts["M4"], r.MidPorts["M7"], r.MidPorts["M8"], r.StopPorts[7])
			go reo.ReplicatorChannel(r.MidPorts["M6"], r.MidPorts["M9"], r.MidPorts["M10"], r.StopPorts[8])
			go reo.SyncdrainChannel(r.MidPorts["M8"], r.MidPorts["M9"], r.StopPorts[9])
			go reo.LossysyncChannel(r.MidPorts["M7"], r.MidPorts["M11"], r.StopPorts[10])
			go reo.SyncChannel(r.MidPorts["M13"], r.MidPorts["M10"], r.StopPorts[11])
			go reo.ReplicatorChannel(r.MidPorts["M11"], r.MidPorts["M12"], r.MidPorts["M13"], r.StopPorts[12])
		}
		return r
	}
	return o
}
