package sul

/*
	Written by Li Yi
	@ 6th Nov 2015
*/

import (
	"../reo"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

/********************* Basic Configurations **********************/

var logger *log.Logger = log.New(os.Stderr, "SUL - ", 2)
var ibound = 3
var ebound = 4

func SetLog(w io.Writer) {
	logger = log.New(w, "SUL", 2)
}

func CloseLog() {
	SetLog(ioutil.Discard)
}

func CloseReoLog() {
	reo.CloseLog()
}

func SetReoDelay(t time.Duration) {
	reo.SetDelay(t)
}

func SetBound(b int) {
	ibound = b
}

func SetEquivBound(b int) {
	ebound = b
}

/********************* Definitions of Input/Output **********************/

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

func (self *Input) String() string {
	rel := ""
	if self.IsTime {
		return "T"
	} else {
		for key, val := range self.Datum {
			if val {
				rel += key + ","
			}
		}
		if rel == "" {
			return "Ø"
		} else {
			return rel[:len(rel)-1]
		}
	}
}

func (self InputSeq) String() string {
	rel := ""
	for _, d := range self {
		rel += d.String() + "-"
	}
	if rel == "" {
		rel += "ϵ"
	} else {
		rel = rel[:len(rel)-1]
	}
	return rel
}

func (self Output) String() string {
	rel := ""
	for key, val := range self {
		if !val.IsEmpty {
			rel += fmt.Sprintf("%s:%s,", key, val.Datum)
		}
	}
	if rel == "" {
		rel += "ϵ"
	}
	return rel
}

func (self Input) EqualTo(i Input) bool {
	if i.IsTime != self.IsTime {
		return false
	}
	if i.IsTime == true {
		return true
	}
	// suppose all the keys in i and self are the same ones
	// otherwise this will be difficult to handle
	for key, val := range self.Datum {
		if val != i.Datum[key] {
			return false
		}
	}
	return true
}

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

func (self OutputSeq) EqualTo(o OutputSeq) bool {
	if len(self) != len(o) {
		return false
	} else {
		for i := 0; i < len(o); i++ {
			if !self[i].EqualTo(&o[i]) {
				return false
			}
		}
	}
	return true
}

/********************* Definition of SulInst ******************/

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
	Cache        *tnode
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
	// fmt.Println("STOP FLAG SET ON")
	// we need more iterations to stop all channels
	// cmflag is used to terminate the monitor goroutine
	// the monitor goroutine is used to deal with the
	// waiting SyncRead/SyncWrite operations
	cmflag := make(chan bool)
	// wait until all the channels terminate
	for _, p := range self.StopPorts {
		<-p.Slave
	}
	// fmt.Println("STOP WAIT FIN")
	close(cmflag)
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

func (self *Oracle) GetInputIndex(item Input) int {
	ins := self.GetInputs()
	for i := 0; i < len(ins); i++ {
		if ins[i].EqualTo(item) {
			return i
		}
	}
	panic("there's an undefined action " + item.String())
}

/********************************* PERFORMANCE ANALYSIS ******************************************/
// counter of membership query
var mqcounter int = 0
var rdcounter int = 0

func CounterReset() {
	mqcounter = 0
	rdcounter = 0
}

func Counter() (int, int) {
	return mqcounter, rdcounter
}

// time-cost analysis
var mquerytime float64 = 0

func MembershipTime() float64 {
	return mquerytime
}

// tree-optimization switch
var treeopt = true

func ToggleTreeOptimization() {
	treeopt = !treeopt
}

/************************************************************************************************/

// NOTE this function is used to process directly simultion on suls
func (self *Oracle) SeqSimulateIteration(ins InputSeq) OutputSeq {
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
	for index, in := range ins {
		// this log line is used to divide different behaviors in reolib
		reo.GetLogger().Println("[SEQ SIM] ITERATE", index, "======================================")
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
				fmt.Println("FATAL ERROR: empty data fetched.")
			} else {
				curr[name] = SingleOutput{data, false}
			}
		}
		out = append(out, curr)
	}
	return out
}

func (self *Oracle) SeqSimulate(ins InputSeq) OutputSeq {
	// we use cache technique to improve the effiency of MQuery,
	// otherwise this would make it really slow
	if self.Cache == nil {
		self.Cache = makenode()
	} else {
		r, ok := self.Cache.search(ins)
		if ok {
			rdcounter++
			logger.Println("[SEQSIMULATE]", ins.String(), "REDUCE: ", rdcounter)
			return r
		}
	}

	logger.Println("[SEQSIMULATE]", ins.String(), "COUNTER: ", mqcounter)
	mqcounter++

	var ct = 0
	var count = 0
	var seq OutputSeq
	var rec OutputSeq = OutputSeq{}
	starttime := time.Now()

	for ct <= ibound {
		count++
		reo.GetLogger().Println("[SEQSIMULATE] ITERATE", count, "---------------------------------------")
		seq = self.SeqSimulateIteration(ins)
		if len(seq) == 0 {
			// a panic happens
			logger.Println("[SEQSIMULATE] PANIC CAUGHT")
			continue
		} else {
			if ct > 0 && seq.EqualTo(rec) {
				ct++
			} else {
				// ct = 0 : the first iteration || not equal : means there's an error
				ct = 1
				rec = seq
			}
		}
	}

	// save simulation result to cache
	if treeopt {
		self.Cache.insert(ins, seq)
	}

	mquerytime += time.Now().Sub(starttime).Seconds()
	return seq
}

func (self *Oracle) MQuery(in InputSeq) Output {
	seq := self.SeqSimulate(in)
	return seq[len(seq)-1]
}

type Executable interface {
	Run(InputSeq) (InputSeq, Output)
	SeqRun(InputSeq) ([]InputSeq, []Output)
}

func (self Oracle) EQuery(hypo Executable) (InputSeq, bool) {
	// FIXME bound is set for iteration depth, but who will tell us its value ?
	var bound = ebound
	/*
		  1. generate a series of all permutation
		  2. TODO check if these sequences are in the tree
			   - if true, just ignore them
				 - if false, use SeqRun to generate it's corresponding hypo-execution
				   to compare and see if we can find a counter-example
	*/
	seqs := []InputSeq{InputSeq{}}
	acts := self.GetInputs()
	// FIXME coud we encapsulate this part into a single function? this may lead to better
	// performance
	for i := 0; i < bound; i++ {
		nseqs := []InputSeq{}
		for j := 0; j < len(seqs); j++ {
			for k := 0; k < len(acts); k++ {
				// check if the sequence is existing in the tree
				curr := InputSeq{}
				// NOTE still not sure what cause the problem
				curr = append(curr, seqs[j]...)
				curr = append(curr, acts[k])
				nseqs = append(nseqs, curr)
			}
		}
		seqs = nseqs
	}
	// existing-check and generation of it's corresponding execution under hypothesis
	for i := 0; i < len(seqs); i++ {
		_, ok := self.Cache.search(seqs[i])
		if !ok {
			_, hypout := hypo.SeqRun(seqs[i])
			sysout := self.SeqSimulate(seqs[i])
			logger.Println("[EQUERY] HypoCheck", seqs[i], hypout, sysout)
			for j := 0; j < len(hypout); j++ {
				if hypout[j].String() != sysout[j].String() {
					// counter-example found
					// FIXME further analysis maybe required
					logger.Println("[EQUERY] counter-example found", seqs[i][:j+1], "with")
					logger.Println("         hypout:", hypout)
					logger.Println("         sysout:", sysout)
					return seqs[i][:j+1], true
				}
			}
		} else {
			// optimization counter
			rdcounter++
		}
	}
	return InputSeq{}, false
}

func main() {
	fmt.Println("Compiled as Main")
}
