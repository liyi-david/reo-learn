package learn

import "../sul"
import "fmt"
import "strconv"
import "log"
import "os"
import "io"
import "io/ioutil"
import "time"

/*
	Created By Li Yi @ Nov 17
	first of all, we have to define some data structures
*/

var logger *log.Logger = log.New(os.Stderr, "LEARN - ", 2)

func SetLog(w io.Writer) {
	logger = log.New(w, "SUL", 2)
}

func CloseLog() {
	SetLog(ioutil.Discard)
}

// we use ObsLine to store a single line in Obs table

type ObsLine struct {
	Index sul.InputSeq
	// using link array here leads to some inconvinence
	// however, we need to make an element in Result able
	// to be assigned nil, indicating that the cell hasn't
	// been calculated yet
	Result     []*sul.Output
	AccessLine int
	// lines can be equivalent if they share the same output signature
	// AccessLine indicates the first line which has the same output
	// with self
	Dist []int
	// Dist maps an action (by it's index) to its successor line
	// denoted by its index too
	// Every time we expandLp, the Dist of each line would be refreshed
}

func NewLine(index sul.InputSeq) ObsLine {
	newl := ObsLine{
		index, []*sul.Output{}, -1, []int{},
	}
	return newl
}

func (self ObsLine) LimitEqualTo(l ObsLine, limit int) bool {
	for i := 0; i < limit; i++ {
		if !self.Result[i].EqualTo(l.Result[i]) {
			return false
		}
	}
	return true
}

func (self ObsLine) EqualTo(l ObsLine) bool {
	// first check if the two lines has the same length
	// actually we suppose they have, otherwise this could be
	// a terrible bug
	if len(self.Result) != len(l.Result) {
		return false
	} else {
		for i := 0; i < len(l.Result); i++ {
			if !self.Result[i].EqualTo(l.Result[i]) {
				return false
			}
		}
	}
	return true
}

type Obs struct {
	// part 1. basic variables
	// describing same arguments as in the paper
	D     []sul.InputSeq
	SL    []ObsLine // including Sp and Lp
	SpLoc int       // indicating [0 .. SpLen] of SL is Sp
	// part 2. private variables
	// -- orac
	// -- lastadd: the last added suffix (from counter-example)
	orac    *sul.Oracle
	lastadd sul.InputSeq
	states  []int // a list of valid states (there indexes in SL)
}

// suppose all the elements in SL are now in Sp
// we need to expand all possible Sp to a new Lp
// also we call call it one-step expansion
func (self *Obs) expandLp() {
	// refresh the amount of Sp
	// remove redundant states in SL
	newSL := []ObsLine{}
	for i, _ := range self.SL {
		flag := true
		for j := 0; j < len(newSL); j++ {
			if newSL[j].EqualTo(self.SL[i]) {
				flag = false
				break
			}
		}
		if flag {
			newSL = append(newSL, self.SL[i])
		}
	}
	self.SL = newSL
	self.SpLoc = len(self.SL) - 1
	newLp := []ObsLine{}
	acts := self.orac.GetInputs()
	for i, _ := range self.SL {
		self.SL[i].Dist = []int{}
		for _, d := range acts {
			newLp = append(newLp, NewLine(append(self.SL[i].Index, d)))
			self.SL[i].Dist = append(self.SL[i].Dist, len(self.SL)+len(newLp)-1)
		}
	}
	self.SL = append(self.SL, newLp...)
}

/********************************* PERFORMANCE ANALYSIS ******************************************/
// time-cost analysis
var hyporuntime float64 = 0

func RunTime() float64 {
	return hyporuntime
}

/*************************************************************************************************/

func (self *Obs) SeqRun(in sul.InputSeq) ([]sul.InputSeq, []sul.Output) {
	var loc = 0
	var rel *sul.Output
	starttime := time.Now()

	indarr, outarr := []sul.InputSeq{}, []sul.Output{}
	for i := 0; i < len(in); i++ {
		inputindex := self.orac.GetInputIndex(*in[i])
		rel = self.SL[loc].Result[inputindex]
		loc = self.SL[loc].Dist[inputindex]
		indarr = append(indarr, self.SL[loc].Index)
		outarr = append(outarr, *rel)
	}

	hyporuntime += time.Now().Sub(starttime).Seconds()
	return indarr, outarr
}

func (self *Obs) Run(in sul.InputSeq) (sul.InputSeq, sul.Output) {
	seqs, outs := self.SeqRun(in)
	tail := len(seqs) - 1
	if tail < 0 {
		panic("hypothesis cannot be executed with no input")
	}
	return seqs[tail], outs[tail]
}

func (self *Obs) AddSuffix(suf sul.InputSeq) {
	self.D = append(self.D, suf)
	logger.Println("[SUFFIX ADD]", suf)
	if len(suf) == 0 {
		panic("fatal error: empty suffixes added")
	}
}

func (self *Obs) GetHypoStr() string {
	rel := "Hypothesis Acquired: \n"
	acts := self.orac.GetInputs()
	for i := 0; i <= self.SpLoc; i++ {
		if self.SL[i].AccessLine == i {
			// then this is a state
			rel += "> state " + strconv.Itoa(i) + ": "
			rel += self.SL[i].Index.String()
			rel += "  with edges: \n"
			for j := 0; j < len(acts); j++ {
				rel += fmt.Sprintf("[%s]\t -> state %2d with output %s \n", acts[j].String(), self.SL[i].Dist[j], self.SL[i].Result[j].String())
			}
		}
	}
	rel += "Hypothesis's description is finished."
	return rel
}

func (self *Obs) String() string {
	rel := "Observation Table: \n"
	rel += "\t|"
	for i := 0; i < len(self.D); i++ {
		rel += fmt.Sprintf("%s\t", self.D[i].String())
	}
	rel += "\n"
	rel += "---------------------------------------\n"
	// every line
	for i := 0; i < len(self.SL); i++ {
		rel += fmt.Sprintf("|%s\t|", self.SL[i].Index.String())
		for j := 0; j < len(self.D); j++ {
			// fmt.Println(self.SL[i].Result)
			var nxt = ""
			// NOTE in one line, there could be more columns than actions
			// especially when some suffixed are added. so we use j < len(self.SL[i].Dist) to make sure
			// panics won't happen
			if i <= self.SpLoc && j < len(self.SL[i].Dist) {
				nxt = fmt.Sprintf("#%d", self.SL[i].Dist[j])
			}
			rel += fmt.Sprintf("%s%s\t", self.SL[i].Result[j].String(), nxt)
		}
		// check if the current line has its corresponding state
		if i > self.SpLoc && self.SL[i].AccessLine == -1 {
			rel += "\t [UNCLOSED]"
		}
		rel += "\n"
		if i == self.SpLoc {
			rel += "---------------------------------------\n"
		}
	}
	return rel
}

func ObsInit(orac *sul.Oracle) *Obs {
	inst := new(Obs)
	inst.orac = orac
	// first element of SL should be \varpesilon
	inst.SL = []ObsLine{NewLine([]*sul.Input{})}
	// here I've used a lambda function to simulate
	// map operation in functional languages
	inst.D = func(items []*sul.Input) []sul.InputSeq {
		// map the single actions to a trivial sequence
		rel := []sul.InputSeq{}
		for _, item := range items {
			rel = append(rel, sul.InputSeq{item})
		}
		return rel
	}(inst.orac.GetInputs())
	inst.SpLoc = 0
	inst.expandLp()
	return inst
}

func main() {
	fmt.Println("Juse make sure fmt would not lead to error ...")
}
