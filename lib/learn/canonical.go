package learn

func arraycompare(a, b []int) bool {
	if len(a) != len(b) {
		return false
	} else {
		for i, _ := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}

func psearch(parts [][]int, state int) int {
	for i, p := range parts {
		for _, s := range p {
			if s == state {
				return i
			}
		}
	}
	return -1
}

func psucc(parts [][]int, dist []int) []int {
	ndist := []int{}
	for _, t := range dist {
		ndist = append(ndist, psearch(parts, t))
	}
	return ndist
}

// this function is used to divide a partition to several
// part w.r.t to their successor
// the function returns several sub-partitions
func (self *Obs) partitionDivide(p [][]int, index int) [][]int {
	logger.Println("Partition Division")
	logger.Println(p)
	// map maps an action to its following state (corresponding index)
	newpts := [][]int{}
	newdists := [][]int{}
	for _, i := range p[index] {
		flag := false
		currdist := psucc(p, self.SL[i].Dist)
		// search the existing sub partitions
		for j, _ := range newpts {
			if arraycompare(newdists[j], currdist) {
				newpts[j] = append(newpts[j], i)
				flag = true
				break
			}
		}
		if !flag {
			// we need to add a new sub-partition
			newpts = append(newpts, []int{i})
			newdists = append(newdists, currdist)
		}
	}
	logger.Println(newpts)
	return newpts
}

func (self *Obs) Canonical() bool {
	pt := [][]int{}
	// D[0 .. lim - 1] are the suffixes with length 1
	lim := len(self.orac.GetInputs())
	// step 1. initialize the partitions
	for i := 0; i <= self.SpLoc; i++ {
		flag := false
		// we see if this line belongs to some existing partition
		for j := 0; j < len(pt); j++ {
			// we assume that there're at least one element in each partition
			if self.SL[i].LimitEqualTo(self.SL[pt[j][0]], lim) {
				pt[j] = append(pt[j], i)
				flag = true
				break
			}
		}
		// if we need to create a new partition
		if !flag {
			pt = append(pt, []int{i})
		}
	}
	// step 2. iteration
	// if there's no change in an iteration, the loop will exit
	for changed := true; changed; {
		changed = false
		newpt := [][]int{}
		for i := 0; i < len(pt); i++ {
			// recheck partition i
			tgroup := self.partitionDivide(pt, i)
			if len(tgroup) > 1 {
				newpt = append(newpt, tgroup...)
				changed = true
			}
		}
		pt = newpt
	}
	// step 3. check if it is canonical
	for _, p := range pt {
		if len(p) > 1 {
			return false
		}
	}
	return true
}
