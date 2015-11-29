package reo

// import "fmt"
import "sync"

// FIXME maybe we need to send a stop signal to any potential
// blocking operation? any SyncRead may leads to this kind of bugs
// NOTE when later some more complicated example would be trapped
// into deadlocks plz kindly check this
// FIXME maybe there're lots of WaitRead that need to replaced by
// select ...

func SyncChannel(in, out, stop Port) {
	defer close(stop.Slave)
	for {
		status := StepExec(
			stop.Main,
			//Operation{"debug", in.Slave, "SYNC LISTENING"},
			Operation{"read", in.Slave, ""},
			//Operation{"debug", in.Slave, "SYNC READ TRIGGED"},
			Operation{"write", out.Slave, "write"},
			Operation{"write", in.Slave, "read"},
			Operation{"read", out.Slave, ""},
			Operation{"read", in.Main, "datum"},
			Operation{"write", out.Main, "datum"},
		)
		if !status {
			return
		}
	}
}

func SyncdrainChannel(in1, in2, stop Port) {
	defer close(stop.Slave)
	for {
		status := StepExec(
			stop.Main,
			Operation{"read", in1.Slave, ""},
			Operation{"read", in2.Slave, ""},
			Operation{"write", in1.Slave, "read"},
			Operation{"write", in2.Slave, "read"},
			Operation{"read", in1.Main, ""},
			Operation{"read", in2.Main, ""},
		)
		if !status {
			return
		}
	}
}

func FifoChannel(in, out, stop Port) {
	defer close(stop.Slave)
	for {
		status := StepExec(
			stop.Main,
			Operation{"read", in.Slave, ""},
			Operation{"write", in.Slave, "read"},
			Operation{"read", in.Main, "datum"},
			// following are the syncwrite part
			Operation{"write", out.Slave, "write"},
			Operation{"read", out.Slave, ""},
			Operation{"write", out.Main, "datum"},
		)
		if !status {
			return
		}
	}
}

func LossysyncChannel(in, out, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		}
		// FIXME the SyncRead operation may blocks this channel
		// and hence it cannot be closed by stop Port
		c := in.SyncRead()
		select {
		// try WaitWrite
		case out.Slave <- "write":
			out.ConfirmWrite()
			out.Write(c)
		default:
			// do nothing
		}
	}
}

func MergerChannel(in1, in2, out, stop Port) {
	defer close(stop.Slave)
	for {
		// considering the syntax of select, here we use
		// <-in.slave instead of in.WaitRead()
		select {
		case <-stop.Main:
			return
		case <-in1.Slave:
			status := StepExec(
				stop.Main,
				Operation{"write", out.Slave, "write"},
				Operation{"write", in1.Slave, "read"},
				Operation{"read", out.Slave, ""},
				Operation{"read", in1.Main, "datum"},
				Operation{"write", out.Main, "datum"},
			)
			if !status {
				return
			}
		case <-in2.Slave:
			status := StepExec(
				stop.Main,
				Operation{"write", out.Slave, "write"},
				Operation{"write", in2.Slave, "read"},
				Operation{"read", out.Slave, ""},
				Operation{"read", in2.Main, "datum"},
				Operation{"write", out.Main, "datum"},
			)
			if !status {
				return
			}
		}
	}
}

func ReplicatorChannel(in, out1, out2 Port, stop Port) {
	defer close(stop.Slave)
	for {
		status := StepExec(
			stop.Main,
			Operation{"read", in.Slave, ""},
			//Operation{"debug", in.Slave, "REPLICATOR FIN READ"},
			Operation{"write", out1.Slave, "write"},
			//Operation{"debug", in.Slave, "REPLICATOR FIN FIRST SHAKEHAND"},
			Operation{"write", out2.Slave, "write"},
			Operation{"write", in.Slave, "read"},
			Operation{"read", out1.Slave, ""},
			Operation{"read", out2.Slave, ""},
			Operation{"read", in.Main, "datum"},
			Operation{"write", out1.Main, "datum"},
			Operation{"write", out2.Main, "datum"},
		)
		if !status {
			return
		}
	}
}

func BufferChannel(in, out, stop Port) {
	defer close(stop.Slave)
	buf := []string{}
	var wg sync.WaitGroup
	// listening input
	go func() {
		defer wg.Done()
		c := make(chan string, 1)
		for {
			status := StepExec(
				stop.Main,
				Operation{"read", in.Slave, ""},
				Operation{"write", in.Slave, "read"},
				Operation{"read", in.Main, "datum"},
				Operation{"write", c, "datum"},
			)
			if !status {
				return
			} else {
				buf = append(buf, <-c)
				// fmt.Println("PUSHED", buf)
			}
		}
	}()
	// listening output
	go func() {
		defer wg.Done()
		for {
			if len(buf) == 0 {
				select {
				case <-stop.Main:
					return
				default:
					continue
				}
			}
			status := StepExec(
				stop.Main,
				Operation{"write", out.Slave, "write"},
				Operation{"read", out.Slave, ""},
				Operation{"write", out.Main, buf[0]},
			)
			if !status {
				return
			} else {
				// fmt.Println("WRITTEN", buf[0])
				buf = buf[1:]
			}
		}
	}()
	wg.Add(2)
	wg.Wait()
}
