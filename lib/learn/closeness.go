package learn

func (self *Obs) Partition() {
}

func (self *Obs) Canonical() bool {
	return false
}

func (self *Obs) fillTable() {
	for _, line := range self.SL {
		for i, d := range self.D {
			str := append(line.Index, d...)
			rel := self.orac.MQuery(str)
			if len(line.Result) <= i {
				line.Result = append(line.Result, &rel)
			} else {
				// simply assignment
				line.Result[i] = &rel
			}
		}
	}
}

func (self *Obs) TableClose() {
	var i, j int
	for self.fillTable(); ; self.fillTable() {
		// check if the table is closed
		flag := true
		for i = self.SpLoc + 1; i < len(self.SL); i++ {
			for j = 0; j <= self.SpLoc; j++ {
				if self.SL[i].EqualTo(self.SL[j]) {
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
			self.expandLp()
		}
	}
}

func (self *Obs) SuffixClose() {
	// TODO
}
