package reo

import "time"
import "fmt"

/*
	Written by Li Yi
	@ 4th Nov 2015
	Ports have similar functions with nodes in Reo.
	In this library, we declare:
	- the type Port and Ports
	- a simple shake hand protocal of comunication between Port
*/
type Port struct {
	Main  chan string
	Slave chan string
}

// how many milliseconds we will wait until a
// TryRead() operations finishes
var Delay time.Duration = 5

func SetDelay(delay time.Duration) {
	Delay = delay
}

/*

	  Writer                        Reader
	    |                             |
	WaitWrite() ---- "write" ----> WaitRead()
	    |                             |
Confirmwrite() <-- "read" ---- ConfirmRead()
      |                             |
		Write() ------- datum -------> Read()

*/
func (self Port) WaitRead() {
	<-self.Slave
}

func (self Port) ConfirmRead() {
	self.Slave <- "read"
}

func (self Port) Read() string {
	return <-self.Main
}

func (self Port) WaitWrite() {
	self.Slave <- "write"
}

func (self Port) ConfirmWrite() {
	<-self.Slave
}

func (self Port) Write(c string) {
	self.Main <- c
}

func MakePort() Port {
	m := make(chan string)
	s := make(chan string)
	return Port{m, s}
}

/*
	From now on there are several encapsulated functions.
	They are not required by just act as some shortcuts
*/
func (self Port) SyncRead() string {
	self.WaitRead()
	self.ConfirmRead()
	return self.Read()
}

func (self Port) SyncWrite(c string) {
	self.WaitWrite()
	self.ConfirmWrite()
	self.Write(c)
}

func (p Port) TryRead(buf chan string) chan bool {
	stopflag := make(chan bool)
	go func() {
		c := make(chan string, 1)
		status := TimedStepExec(
			time.Millisecond*Delay,
			Operation{"read", p.Main, "buf"},
			Operation{"write", c, "buf"},
		)
		if !status {
			buf <- "<NONE>"
			logger.Println("<TRY READ>", "TIMEOUT")
		} else {
			s := <-c
			logger.Println("<TRY READ>", "TRANS", s)
			buf <- s
		}
		close(stopflag)
	}()
	return stopflag
}

func (p Port) LossyWrite(c string) chan bool {
	stopflag := make(chan bool)
	go func() {
		s := TimedStepExec(
			time.Millisecond*Delay,
			Operation{"write", p.Slave, "write"},
			Operation{"read", p.Slave, ""},
			Operation{"write", p.Main, c},
		)
		if !s {
			logger.Println("<TRY WRITE>", "TIMEOUT", c)
		}
		close(stopflag)
	}()
	return stopflag
}

func GenerateStopPort(n int) []Port {
	stopports := []Port{}
	stopflag := make(chan string)
	for i := 0; i < n; i++ {
		stopports = append(stopports, Port{stopflag, make(chan string)})
	}
	return stopports
}

func main() {
	// this line is just make sure that
	// there're some lines using fmt
	fmt.Println("Self Build")
}
