package reo

import "sync"
import "log"
import "io"
import "io/ioutil"
import "os"

// FIXME maybe we need to send a stop signal to any potential
// blocking operation? any SyncRead may leads to this kind of bugs
// NOTE when later some more complicated example would be trapped
// into deadlocks plz kindly check this
// FIXME maybe there're lots of WaitRead that need to replaced by
// select ...

var logger *log.Logger = log.New(os.Stderr, "REO - ", 2)

func SetLog(w io.Writer) {
	logger = log.New(w, "REO", 2)
}

func GetLogger() *log.Logger {
	return logger
}

func CloseLog() {
	SetLog(ioutil.Discard)
}

func SyncChannel(in, out, stop Port) {
	defer close(stop.Slave)
	for {
		c := make(chan string, 1)
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
			Operation{"write", c, "datum"},
		)
		if !status {
			return
		}
		logger.Println("[SYNC] TRANS", <-c)
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
		logger.Println("[SYNCDRAIN] TRIGGED")
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
		c := make(chan string, 1)
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
				Operation{"write", c, "datum"},
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
				Operation{"write", c, "datum"},
			)
			if !status {
				return
			}
		}
		logger.Println("[MERGER] TRANS", <-c)
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
		logger.Println("[REPLICATOR] TRANS")
	}
}

func OutputChannel(in, out, stop Port) {
	defer close(stop.Slave)
	for {
		c := make(chan string, 1)
		status := StepExec(
			stop.Main,
			Operation{"read", in.Slave, ""},
			Operation{"write", in.Slave, "read"},
			Operation{"read", in.Main, "datum"},
			Operation{"write", out.Main, "datum"},
			Operation{"write", c, "datum"},
		)
		if !status {
			return
		}
		logger.Println("[OUTPUT] TRANS", <-c)
	}
}

func BufferChannel(size int, in, out, stop Port) {
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
				t := <-c
				if len(buf) < size {
					buf = append(buf, t)
				}
				logger.Println("[BUFFER] PUSHED", t)
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
				logger.Println("[BUFFER] WRITTEN", buf[0])
				buf = buf[1:]
			}
		}
	}()
	wg.Add(2)
	wg.Wait()
}
