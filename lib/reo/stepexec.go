package reo

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

func (self *Flagproc) Execute() (result, ok bool) {
	defer func() {
		// FIXME recover from unhandled panic
		recover()
	}()
	ok = true
	result = false
	for _, o := range self.oprs {
		if o.opr == "read" {
			select {
			case <-self.flag:
				// logger.Println("{STEPEXEC}", "READ INTERRUPTED", o.name)
				return false, true
			case t := <-o.chn:
				if t != "read" && t != "write" {
					logger.Println("{STEPEXEC}", "READ", t, "TO", o.name)
				}
				self.vars[o.name] = t
			}
		} else if o.opr == "debug" {
			logger.Println("{STEPEXEC}", o.name)
		} else {
			// operation should be write
			select {
			case <-self.flag:
				return false, true
			case o.chn <- self.Var(o.name):
			}
		}
	}
	return true, true
}

func StepExec(flag chan string, oprs ...Operation) bool {
	proc := Flagproc{flag, map[string]string{}, oprs}
	for {
		r, ok := proc.Execute()
		if ok {
			return r
		}
	}
}

func TimedStepExec(timeout time.Duration, oprs ...Operation) bool {
	for {
		flag := make(chan string)
		proc := Flagproc{flag, map[string]string{}, oprs}
		go func() {
			<-time.After(timeout)
			close(flag)
		}()
		r, ok := proc.Execute()
		if ok {
			return r
		} else {
			logger.Println(r)
		}
	}
}
