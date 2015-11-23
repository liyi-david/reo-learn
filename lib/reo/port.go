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
const Delay = 1

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

type Ports []Port

func (self Ports) WaitWrite() {
	for i := 0; i < len(self); i++ {
		self[i].Slave <- "write"
	}
}

func (self Ports) ConfirmWrite() {
	for i := 0; i < len(self); i++ {
		<-self[i].Slave
	}
}

func (self Ports) Write(c string) {
	for i := 0; i < len(self); i++ {
		self[i].Write(c)
	}
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
		select {
		case <-p.Slave:
			// FIXME there must be something strange here
			// Two "waiting finish" are emiited
			// thats why we have a deadlock
			p.ConfirmRead()
			// fmt.Println("[Try Read] Waiting Comfirmed")
			c := p.Read()
			buf <- c
		case <-time.After(time.Millisecond * Delay):
			buf <- "<NONE>"
		}
		close(stopflag)
	}()
	return stopflag
}

func (p Port) LossyWrite(c string) chan bool {
	stopflag := make(chan bool)
	go func() {
		select {
		case p.Slave <- "write":
			p.ConfirmWrite()
			p.Write(c)
		case <-time.After(time.Millisecond * Delay):
			// nothing done.
		}
		close(stopflag)
	}()
	return stopflag
}

func (p Port) UselessRead(stop chan bool) {
	select {
	case <-p.Slave:
		select {
		case <-stop:
		case p.Slave <- "read":
			select {
			case <-stop:
			case <-p.Main:
			}
		}
	case <-stop:
	}
}

func (p Port) UselessWrite(stop chan bool) {
	select {
	case p.Slave <- "write":
		select {
		case <-stop:
		case <-p.Slave:
			select {
			case <-stop:
			case p.Main <- "":
			}
		}
	case <-stop:
	}
}

func GenerateStopPort(n int) Ports {
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
