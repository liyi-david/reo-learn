package expfifo

import "../../lib/sul"
import "../../lib/reo"
import "time"
import "strconv"
import "fmt"

func GetOracle() *sul.Oracle {
	o := new(sul.Oracle)
	o.InPorts = []string{"A"}
	// dynamically generate the MidPorts
	o.MidPorts = []string{}
	for i := 0; i < 11; i++ {
		o.MidPorts = append(o.MidPorts, "M"+strconv.Itoa(i))
	}
	fmt.Println(o.MidPorts)

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
		r.StopPorts = reo.GenerateStopPort(10)
		r.Start = func() {
			// Input & Output
			go reo.BufferChannel(1, r.InPorts["A"], r.MidPorts["M0"], r.StopPorts[0])
			go reo.OutputChannel(r.MidPorts["M10"], r.OutPorts["B"], r.StopPorts[1])
			// Function Channels
			go reo.ReplicatorChannel(r.MidPorts["M0"], r.MidPorts["M1"], r.MidPorts["M2"], r.StopPorts[2])
			go reo.FifoChannel(r.MidPorts["M1"], r.MidPorts["M3"], r.StopPorts[3])
			go reo.TimerChannel(r.MidPorts["M2"], r.MidPorts["M4"], 80*time.Millisecond, r.StopPorts[4])
			go reo.LossysyncChannel(r.MidPorts["M4"], r.MidPorts["M5"], r.StopPorts[5])
			go reo.ReplicatorChannel(r.MidPorts["M3"], r.MidPorts["M6"], r.MidPorts["M7"], r.StopPorts[6])
			go reo.ReplicatorChannel(r.MidPorts["M5"], r.MidPorts["M8"], r.MidPorts["M9"], r.StopPorts[7])
			go reo.SyncdrainChannel(r.MidPorts["M7"], r.MidPorts["M8"], r.StopPorts[8])
			go reo.ReplicatorChannel(r.MidPorts["M6"], r.MidPorts["M10"], r.MidPorts["M9"], r.StopPorts[9])
		}
		return r
	}
	return o
}
