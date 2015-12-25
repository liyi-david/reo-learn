package reo

import "sync"
import "log"
import "io"
import "io/ioutil"
import "os"
import "time"

// NOTE when later some more complicated example would be trapped
// into deadlocks plz kindly check this
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
			Operation{"debug", in1.Slave, "SYNCDRAIN first hand-shake finished"},
			Operation{"write", in1.Slave, "read"},
			Operation{"write", in2.Slave, "read"},
			Operation{"debug", in1.Slave, "SYNCDRAIN second hand-shake finished"},
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
	c := make(chan string, 1)
	for {
		status := StepExec(
			stop.Main,
			Operation{"read", in.Slave, ""},
			Operation{"write", in.Slave, "read"},
			Operation{"read", in.Main, "datum"},
			Operation{"write", c, "datum"},
			Operation{"debug", c, "[FIFO] BUFFERED"},
			// following are the syncwrite part
			Operation{"write", out.Slave, "write"},
			Operation{"read", out.Slave, ""},
			Operation{"write", out.Main, "datum"},
		)
		if !status {
			return
		} else {
			logger.Println("[FIFO] RELEASED", <-c)
		}
	}
}

func timer(t time.Duration) chan string {
	r := make(chan string)
	go func() {
		// use time.After
		<-time.After(t)
		close(r)
	}()
	return r
}

func TimerChannel(in, out Port, t time.Duration, stop Port) {
	defer close(stop.Slave)
	c := make(chan string, 1)
	for {
		status := StepExec(
			stop.Main,
			Operation{"read", in.Slave, ""},
			Operation{"write", in.Slave, "read"},
			Operation{"read", in.Main, "datum"},
			Operation{"write", c, "datum"},
			Operation{"debug", c, "[FIFO] BUFFERED"},
			Operation{"read", timer(t), ""},
			// following are the syncwrite part
			Operation{"write", out.Slave, "write"},
			Operation{"read", out.Slave, ""},
			Operation{"write", out.Main, "datum"},
		)
		if !status {
			return
		} else {
			logger.Println("[TIMER] RELEASED", <-c)
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

/* semantics of merger channel:
a) first hand-shake can be applied by both inputs and they are not
   mutual exclusive
*/

func MergerChannel(in1, in2, out, stop Port) {
	defer close(stop.Slave)
	var wg sync.WaitGroup
	var lock sync.Mutex
	var step int = 0
	// we need some assistant function to deal with step
	pushData := func(d string) {
		logger.Println("[MERGER] start pushing data")
		select {
		case <-stop.Main:
			logger.Println("[MERGER] pushdata timeout")
			return
		case out.Main <- d:
			step = 0
			logger.Println("[MERGER] data pushed.", d)
		}
	}

	listener := func(in Port, label string) {
		defer wg.Done()
		// start working
		for {
			// first hand-shake with input port
			select {
			case <-stop.Main:
				return
			case <-in.Slave:
				// first hand-shake with output port
				lock.Lock()
				if step == 0 {
					select {
					case <-stop.Main:
						lock.Unlock()
						return
					case out.Slave <- "write":
						step++
					}
				}
				lock.Unlock()
				// otherwise nothing has to be done
			}
			logger.Println("[MERGER] first handshake", label, step)
			// second hand-shake with output port
			select {
			case <-stop.Main:
				return
			case in.Slave <- "read":
				// second hand-shake with input port
				lock.Lock()
				if step == 0 {
					select {
					case <-stop.Main:
						lock.Unlock()
						return
					case out.Slave <- "write":
						step++
					}
				}
				lock.Unlock()
				lock.Lock()
				if step == 1 {
					select {
					case <-stop.Main:
						lock.Unlock()
						return
					case <-out.Slave:
						step++
					}
				}
				lock.Unlock()
				// otherwise nothing has to be done
			}
			logger.Println("[MERGER] second handshake", label, step)
			// final hand-shake with input port
			select {
			case <-stop.Main:
				return
			case d := <-in.Main:
				// final hand-shake with output
				lock.Lock()
				if step == 0 {
					select {
					case <-stop.Main:
						lock.Unlock()
						return
					case out.Slave <- "write":
						step++
					}
				}
				lock.Unlock()
				lock.Lock()
				if step == 1 {
					select {
					case <-stop.Main:
						lock.Unlock()
						return
					case <-out.Slave:
						step++
					}
				}
				lock.Unlock()
				// now we have step == 2
				lock.Lock()
				if step == 2 {
					pushData(d)
				}
				lock.Unlock()
			}
			logger.Println("[MERGER] FINAL handshake", label, step)
		}
	}

	wg.Add(2)
	go listener(in1, "in1")
	go listener(in2, "in2")
	wg.Wait()
}

//func NMergerChannel(in1, in2, out, stop Port) {
//defer close(stop.Slave)
//for {
//// considering the syntax of select, here we use
//// <-in.slave instead of in.WaitRead()
//c := make(chan string, 1)
//select {
//case <-stop.Main:
//return
//case <-in1.Slave:
//status := StepExec(
//stop.Main,
//Operation{"write", out.Slave, "write"},
//Operation{"write", in1.Slave, "read"},
//Operation{"read", out.Slave, ""},
//Operation{"read", in1.Main, "datum"},
//Operation{"write", out.Main, "datum"},
//Operation{"write", c, "datum"},
//)
//if !status {
//return
//}
//case <-in2.Slave:
//status := StepExec(
//stop.Main,
//Operation{"write", out.Slave, "write"},
//Operation{"write", in2.Slave, "read"},
//Operation{"read", out.Slave, ""},
//Operation{"read", in2.Main, "datum"},
//Operation{"write", out.Main, "datum"},
//Operation{"write", c, "datum"},
//)
//if !status {
//return
//}
//}
//logger.Println("[MERGER] TRANS", <-c)
//}
//}

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
			Operation{"read", out1.Slave, ""},
			Operation{"read", out2.Slave, ""},
			Operation{"write", in.Slave, "read"},
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
				Operation{"read", in.Main, "datum"},
				Operation{"write", c, "datum"},
			)
			if !status {
				return
			} else {
				t := <-c
				if len(buf) < size {
					buf = append(buf, t)
					logger.Println("[BUFFER] PUSHED", t)
				} else {
					logger.Println("[BUFFER] FULL", t)
				}

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
