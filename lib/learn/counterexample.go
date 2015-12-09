package learn

import "../sul"

// suppose we have a getaccessseq function in obs
// e.g. inst.GetAccessSeq(sul.InputSeq) -> sul.InputSeq

// this function returns a suffix
func (inst *Obs) CEAnalyze(cin sul.InputSeq) sul.InputSeq {
	// generally this is a bisection algorithm
	// written by Yiwu

	oc := inst.orac.MQuery(cin)
	//	acts := inst.orac.GetInputs()
	sin, d, d2 := cin, cin, cin
	lower, upper := 2, len(cin)-1
	for {
		mid := (lower + upper) / 2
		sin = cin[:mid-1]
		sout, _ := inst.Run(sin)
		d = cin[mid-1:]
		d2 = cin[mid:]
		sout = append(sout, d...)
		omid := inst.orac.MQuery(sout)
		if oc.EqualTo(&omid) {
			lower = mid + 1
			if upper < lower {
				return d2
			}
		} else {
			upper = mid - 1
			if upper < lower {
				return d
			}
		}
	}
}
