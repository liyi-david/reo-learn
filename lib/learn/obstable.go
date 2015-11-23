package learn

import "../sul"
import "fmt"

/*
	Created By Li Yi @ Nov 17
	first of all, we have to define some data structures
*/

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
	self.SpLoc = len(self.SL) - 1
	newLp := []ObsLine{}
	acts := self.orac.GetInputs()
	for _, l := range self.SL {
		l.Dist = []int{}
		for j, d := range acts {
			newLp = append(newLp, NewLine(append(l.Index, d)))
			l.Dist = append(l.Dist, self.SpLoc+1+j)
		}
	}
	self.SL = append(self.SL, newLp...)
}

func (self *Obs) Run(in sul.InputSeq) sul.Output {
	// TODO
	return sul.Output{}
}

func (self *Obs) AddSuffix(suf sul.InputSeq) {
	// TODO
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
