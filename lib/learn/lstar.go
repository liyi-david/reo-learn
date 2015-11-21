package learn

import "../sul"

/*
	Created by Li Yi @ Nov 20
	This file includes main algorithm of active learning
	basically implements the L* algorithm
*/

/********************************************************************
BASIC IDEA
> loop
>   loop: not suffixclose (canonical) ?
>     add new suffix
>     table close
>   end
>   counter example ?
>     yes - add new suffix
>     no - return
>   end
> end
********************************************************************/

func LStar(orac *sul.Oracle) *Obs {
	// step 1. initialize
	obs := ObsInit(orac)
	// if we have a existing counter-example
	c := sul.InputSeq{}
	// cexist indicates if there's a counter-example found in last round
	// it's initialized as true otherwise the loop will never be executed
	for cexist := true; cexist; {
		// enclose the table
		obs.TableClose()
		for sc := obs.Canonical(); sc; sc = obs.Canonical() {
			obs.SuffixClose()
			obs.TableClose()
		}
		// looking up for counter-examples
		c, cexist = orac.EQuery(obs)
		if cexist {
			// analyze the counter-example
			obs.AddSuffix(CEAnalyze(c))
		}
	}
	// TODO debug after finished
	return obs
}
