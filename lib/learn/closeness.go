package learn

import "fmt"

func (self *Obs) direct_hypothesis() {
	defer fmt.Println("Hypotheis Constructed")
	self.states = []int{}
	// start iteration through all lines in Sp
	// removing the dumplicated lines
	for i := 0; i <= self.SpLoc; i++ {
		self.SL[i].AccessLine = i
		for j := 0; j < i; j++ {
			if self.SL[i].EqualTo(self.SL[j]) {
				self.SL[i].AccessLine = j
				break
			}
		}
	}
	// updating the edges, make them target to correct states' index
	for i := self.SpLoc + 1; i < len(self.SL); i++ {
		if self.SL[i].AccessLine == -1 {
			panic("Fatal Error: Constructing hypothesis with an unclosed table.")
		}
		self.SL[i].AccessLine = self.SL[self.SL[i].AccessLine].AccessLine
	}
	// updating the distributations
	for i := 0; i <= self.SpLoc; i++ {
		for j, _ := range self.SL[i].Dist {
			// pre-cond : SL[i].Dist[j] indicates the index of corresponding edge
			// the edge is activated by action j
			self.SL[i].Dist[j] = self.SL[self.SL[i].Dist[j]].AccessLine
			// post-cond : SL[i].Dist[j] indicates the target node of corresponding edge
			// it means that after
		}
	}
	// NOTE after executing this function, there's no need of Lp lines
}

func (self *Obs) fillTable() {
	for i, _ := range self.SL {
		// if a fillTable operation is executed, the hypothesis
		// need to be reconstructed
		self.SL[i].AccessLine = -1
		for j, d := range self.D {
			str := append(self.SL[i].Index, d...)
			rel := self.orac.MQuery(str)
			if len(self.SL[i].Result) <= j {
				self.SL[i].Result = append(self.SL[i].Result, &rel)
			} else {
				// simply assignment
				self.SL[i].Result[j] = &rel
			}
		}
	}
}

func (self *Obs) TableClose() {
	var i, j int
	// after close the table, we need to construct a hypothesis
	defer self.direct_hypothesis()
	for self.fillTable(); ; self.fillTable() {
		// check if the table is closed now
		flag := true
		for i = self.SpLoc + 1; i < len(self.SL); i++ {
			for j = 0; j <= self.SpLoc; j++ {
				if self.SL[i].EqualTo(self.SL[j]) {
					// mark the corresponding line
					self.SL[i].AccessLine = j
					break
				}
			}
			if j > self.SpLoc {
				flag = false
				break
			}
		}
		if flag {
			// the obstable has been enclosed
			return
		} else {
			// try extra steps, see if we can make it close
			self.expandLp()
		}
	}
}

func (self *Obs) SuffixClose() {
	// TODO
}
