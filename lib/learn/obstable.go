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
	Result  []*sul.Output
	Partion int // denote the current partion index of this line
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
	D        []sul.InputSeq
	SL       []ObsLine // including Sp and Lp
	SpLoc    int       // indicating [0 .. SpLen] of SL is Sp
	Partions [][]int
	// part 2. private variables
	// -- orac
	// -- lastadd: the last added suffix (from counter-example)
	orac    *sul.Oracle
	lastadd sul.InputSeq
}

// suppose all the elements in SL are now in Sp
// we need to expand all possible Sp to a new Lp
func (self *Obs) expandLp() {
	// refresh the amount of Sp
	self.SpLoc = len(self.SL) - 1
	newLp := []ObsLine{}
	for _, l := range self.SL {

		for _, d := range self.D {
			newLp = append(newLp, ObsLine{append(l.Index, d...), []*sul.Output{}, -1})
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
	inst.SL = []ObsLine{ObsLine{[]*sul.Input{}, []*sul.Output{}, -1}}
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
