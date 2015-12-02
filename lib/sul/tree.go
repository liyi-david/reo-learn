package sul

type tnode struct {
	result Output
	child  map[string]*tnode
}

func makenode() *tnode {
	p := new(tnode)
	p.child = map[string]*tnode{}
	return p
}

func (self *tnode) insert(iseq InputSeq, oseq OutputSeq) {
	p := self
	for i := 0; i < len(iseq); i++ {
		index := iseq[i].String()
		_, ok := p.child[index]
		if !ok {
			// need to create the new node
			p.child[index] = makenode()
		}
		p = p.child[index]
		p.result = oseq[i]
	}
}

func (self *tnode) search(iseq InputSeq) *Output {
	p := self
	for i := 0; i < len(iseq); i++ {
		index := iseq[i].String()
		next, ok := p.child[index]
		if !ok {
			return nil
		} else {
			p = next
		}
	}
	return &p.result
}
