package sul

/*
	Written by Li Yi
	@ 6th Nov 2015
*/

import (
	"../reo"
	"fmt"
	"time"
)

type Input struct {
	Datum  map[string]bool
	IsTime bool
}

type SingleOutput struct {
	Datum   string
	IsEmpty bool
}

type Output map[string]SingleOutput

type InputSeq []*Input
type OutputSeq []Output

func (self SingleOutput) EqualTo(so SingleOutput) bool {
	if self.IsEmpty {
		return so.IsEmpty
	} else {
		return self.Datum == so.Datum
	}
}

func (self *Output) EqualTo(o *Output) bool {
	// we assume that the two output share the same ports
	for key, _ := range *self {
		if !(*self)[key].EqualTo((*o)[key]) {
			return false
		}
	}
	return true
}

// Instance of System Under Test
type SulInst struct {
	// public fields
	InPorts, OutPorts, MidPorts map[string]reo.Port
	OutBufs                     map[string]chan string
	Start                       func()
	StopPorts                   []reo.Port
}

type Oracle struct {
	InPorts      []string
	MidPorts     []string
	OutPorts     []string
	Inputs       []*Input
	TimeUnit     time.Duration
	GenerateInst func() *SulInst
}

func (self *SulInst) GeneratePort(ref *Oracle) {
	self.InPorts = map[string]reo.Port{}
	self.OutPorts = map[string]reo.Port{}
	self.MidPorts = map[string]reo.Port{}
	for _, name := range ref.InPorts {
		self.InPorts[name] = reo.MakePort()
	}
	for _, name := range ref.OutPorts {
		self.OutPorts[name] = reo.MakePort()
	}
	for _, name := range ref.MidPorts {
		self.MidPorts[name] = reo.MakePort()
	}
}

func (self Input) deepcopy() *Input {
	r := new(Input)
	r.Datum = map[string]bool{}
	if self.IsTime {
		r.IsTime = true
	} else {
		for key, val := range self.Datum {
			r.Datum[key] = val
		}
	}
	return r
}

func (self *SulInst) Stop() {
	// NOTE theoretically the array StopPorts includes
	// at least one element since a connector usually
	// contains at least one channel
	close(self.StopPorts[0].Main)
	// we need more iterations to stop all channels
	// cmflag is used to terminate the monitor goroutine
	// the monitor goroutine is used to deal with the
	// waiting SyncRead/SyncWrite operations
	cmflag := make(chan bool)
	go self.monitor(cmflag)
	// wait until all the channels terminate
	for _, p := range self.StopPorts {
		<-p.Slave
	}
	close(cmflag)
}

func (self *SulInst) monitor(stop chan bool) {
	for {
		select {
		case <-stop:
			return
		default:
			// keep writing things and reading things
			for _, p := range self.InPorts {
				p.UselessWrite(stop)
			}
			for _, p := range self.OutPorts {
				p.UselessRead(stop)
			}
		}
	}
}

func (self *Oracle) GetInputs() []*Input {
	if len(self.Inputs) != 0 {
		return self.Inputs
	}
	rel := []*Input{new(Input)}
	// need to initialize the head element manually
	rel[0].Datum = map[string]bool{}
	temp := []*Input{}
	for _, port := range self.InPorts {
		for _, inp := range rel {
			inp.Datum[port] = false
			titm := inp.deepcopy()
			titm.Datum[port] = true
			temp = append(temp, titm)
		}
		rel = append(rel, temp...)
		temp = []*Input{}
	}
	tick := new(Input)
	tick.IsTime = true
	rel = append(rel, tick)
	self.Inputs = rel
	return rel
}

func (self *Oracle) SeqSimulate(ins InputSeq) OutputSeq {
	inst := self.GenerateInst()
	inst.OutBufs = map[string]chan string{}
	// initialization of buffers
	for name, _ := range inst.OutPorts {
		inst.OutBufs[name] = make(chan string, len(ins)+1)
	}
	inst.Start()
	var stopgroup []chan bool
	// use waitgroup to make sure all processes finished
	// before we continue dealing with data
	for _, in := range ins {
		stopgroup = []chan bool{}
		if in.IsTime {
			time.Sleep(self.TimeUnit)
		} else {
			for pname, exist := range in.Datum {
				// push data
				if exist {
					stopgroup = append(stopgroup, inst.InPorts[pname].LossyWrite(pname))
				}
			}
		}
		for name, port := range inst.OutPorts {
			stopgroup = append(stopgroup, port.TryRead(inst.OutBufs[name]))
		}
		// since all the TryRead/Lossy Operation won't take longer than
		// reo.Delay milliseconds
		// wait until all the channels in stopgroup are closed
		for _, c := range stopgroup {
			<-c
		}
	}
	// make sure all the i/o operations are finished
	// then we try to stop the execution of connector
	// fmt.Println("Going to STOP.")
	inst.Stop()
	// fmt.Println("STOP Finished.")
	// generate output
	var out OutputSeq
	for _, _ = range ins {
		curr := Output{}
		for name, _ := range inst.OutPorts {
			data := <-inst.OutBufs[name]
			if data == "<NONE>" {
				curr[name] = SingleOutput{"", true}
			} else if data == "" {
			} else {
				curr[name] = SingleOutput{data, false}
			}
		}
		out = append(out, curr)
	}
	return out
}

var mqcounter int = 0

func CounterReset() {
	mqcounter = 0
}

func Counter() int {
	return mqcounter
}

func (self *Oracle) MQuery(in InputSeq) Output {
	// TODO we should use cache technique to improve
	// the effiency of MQuery, otherwise this would make
	// it really slow
	// fmt.Println("Membership Query.")
	mqcounter++
	outputs := self.SeqSimulate(in)
	if len(outputs) == 0 {
		panic("Fatal Error: SeqSimulate returns an empty array.")
	} else {
		return outputs[len(outputs)-1]
	}
}

// TODO: EQuery should accept an argument
// descirbing the hypothesis
type Executable interface {
	Run(InputSeq) Output
}

func (self Oracle) EQuery(hypo Executable) (InputSeq, bool) {
	// TODO
	return InputSeq{}, false
}

func main() {
	fmt.Println("Compiled as Main")
}
