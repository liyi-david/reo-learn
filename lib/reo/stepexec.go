package reo

import "fmt"
import "time"

type Operation struct {
	opr  string
	chn  chan string
	name string
}

type Flagproc struct {
	// -- flag: when it is closed, all the operations should be halted
	flag chan string
	vars map[string]string
	oprs []Operation
}

func (self Flagproc) Var(name string) string {
	v, ok := self.vars[name]
	if !ok {
		return name
	} else {
		return v
	}
}

func (self Flagproc) Execute() bool {
	for _, o := range self.oprs {
		if o.opr == "read" {
			select {
			case <-self.flag:
				return false
			case t := <-o.chn:
				// fmt.Println("DATA READ: " + t + " TO " + o.name)
				self.vars[o.name] = t
			}
		} else if o.opr == "debug" {
			fmt.Println("DEBUG", o.name)
		} else {
			// operation should be write
			select {
			case <-self.flag:
				return false
			case o.chn <- self.Var(o.name):
			}
		}
	}
	return true
}

func StepExec(flag chan string, oprs ...Operation) bool {
	proc := Flagproc{flag, map[string]string{}, oprs}
	return proc.Execute()
}

func TimedStepExec(timeout time.Duration, oprs ...Operation) bool {
	flag := make(chan string)
	proc := Flagproc{flag, map[string]string{}, oprs}
	go func() {
		<-time.After(timeout)
		close(flag)
	}()
	status := proc.Execute()
	return status
}
