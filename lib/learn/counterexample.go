package learn

import "../sul"

// suppose we have a getaccessseq function in obs
// e.g. inst.GetAccessSeq(sul.InputSeq) -> sul.InputSeq

// this function returns a suffix
func (inst *Obs) CEAnalyze(cin sul.InputSeq) sul.InputSeq {
	// generally this is a bisection algorithm
	// written by Yiwu
	logger.Println(cin)
	oc := inst.orac.MQuery(cin)
	//	acts := inst.orac.GetInputs()
	sin, d, d2 := cin, cin, cin
	lower, upper := 1, len(cin)-1
	for {
		mid := (lower + upper) / 2
		sin = cin[:mid]
		sout, _ := inst.Run(sin)
		logger.Println("hypothesis execution accessseq:", sin, "->", sout)
		d = cin[mid:]
		d2 = cin[mid+1:]
		sout = append(sout, d...)
		omid := inst.orac.MQuery(sout)
		logger.Println("mquery: ", sout, "->", omid)
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
