package rank

//"fmt"

const realRankCount int = 1000
const maxItemCount int = 1000

type rankItemBlock struct {
	items    []node
	nextFree int
}

func (rb *rankItemBlock) get() *node {
	if rb.nextFree >= cap(rb.items) {
		return nil
	} else {
		item := &rb.items[rb.nextFree]
		rb.nextFree++
		return item
	}
}
func (rb *rankItemBlock) reset() {
	for i, _ := range rb.items {
		item := &rb.items[i]
		item.sl = nil
		for j := 0; j < maxLevel; j++ {
			item.links[j].skip = 0
			item.links[j].pnext, item.links[j].pprev = nil, nil
		}
	}
	rb.nextFree = 0
}
func newRankItemBlock() *rankItemBlock {
	return &rankItemBlock{
		items: make([]node, 10000),
	}
}

type rankItemPool struct {
	blocks   []*rankItemBlock
	nextFree int
}

func newRankItemPool() *rankItemPool {
	return &rankItemPool{
		blocks: []*rankItemBlock{newRankItemBlock()},
	}
}
func (rp *rankItemPool) reset() {
	for _, v := range rp.blocks {
		v.reset()
	}
	rp.nextFree = 0
}
func (rp *rankItemPool) get() *node {
	item := rp.blocks[rp.nextFree].get()
	if nil == item {
		block := newRankItemBlock()
		rp.blocks = append(rp.blocks, block)
		rp.nextFree++
		item = block.get()
	}
	return item
}

type Rank struct {
	id2Item   map[uint64]*node
	spans     []*skiplists
	nextShink int
	cc        int
	itemPool  *rankItemPool
}

func NewRank() *Rank {
	return &Rank{
		id2Item:  map[uint64]*node{},
		spans:    make([]*skiplists, 0, 8192),
		itemPool: newRankItemPool(),
	}
}

func (r *Rank) Reset() {
	r.id2Item = map[uint64]*node{}
	r.spans = make([]*skiplists, 0, 8192)
	r.itemPool.reset()
}

func (r *Rank) getRankPersent(item *node, idxInSpan int) int {
	if len(r.spans) < 100 {
		return 100 - idxInSpan*100/maxItemCount*len(r.spans)
	} else {
		return 100 - maxItemCount*item.sl.idx/(len(r.spans)-1)
	}
}

func (r *Rank) getRankPersentByItem(item *node) int {
	if len(r.spans) < 100 {
		return 100 - r.getRank(item)*100/maxItemCount*len(r.spans)
	} else {
		return 100 - maxItemCount*item.sl.idx/(len(r.spans)-1)
	}
}

func (r *Rank) GetRankPersent(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else if len(r.spans) < 100 {
		return 100 - r.getRank(item)*100/maxItemCount*len(r.spans)
	} else {
		return 100 - maxItemCount*item.sl.idx/(len(r.spans)-1)
	}
}

func (r *Rank) getFrontSpanItemCount(item *node) int {
	c := 0
	i := 0
	for ; i < len(r.spans); i++ {
		v := r.spans[i]
		if item.sl == v {
			break
		} else {
			c += v.size
			if c >= realRankCount {
				break
			}
		}
	}

	if i < item.sl.idx {
		c += maxItemCount * (item.sl.idx - i)
	}

	return c
}

func (r *Rank) getRank(item *node) int {
	return r.getFrontSpanItemCount(item) + item.sl.GetNodeRank(item)
}

func (r *Rank) GetRank(id uint64) int {

	r.cc++
	defer func() {
		if r.cc%100 == 0 {
			r.shrink(nil)
		}
	}()

	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return r.getRank(item)
	}
}

func (r *Rank) Check() bool {
	if len(r.spans) == 0 {
		return true
	}
	vv := r.spans[0].max
	for i, v := range r.spans {
		vv = v.check(vv)
		if vv == -1 && i > 0 {
			return false
		}
	}

	return true
}

func (r *Rank) getRankItem(id uint64) *node {
	return r.id2Item[id]
}

func (r *Rank) binarySearch(score int, left int, right int) *skiplists {

	if left >= right {
		return r.spans[left]
	}

	mIdx := (right-left)/2 + left
	m := r.spans[mIdx]

	if m.max > score {
		nIdx := mIdx + 1
		if nIdx >= len(r.spans) || r.spans[nIdx].max < score {
			return m
		}
		return r.binarySearch(score, mIdx+1, right)
	} else {
		pIdx := mIdx - 1
		if pIdx < 0 || r.spans[pIdx].min > score {
			return m
		}
		return r.binarySearch(score, left, mIdx-1)
	}

}

func (r *Rank) findSpan(score int) *skiplists {
	var c *skiplists
	if len(r.spans) == 0 {
		c = newSkipLists(len(r.spans))
		r.spans = append(r.spans, c)
	} else {
		c = r.binarySearch(score, 0, len(r.spans)-1)
	}

	return c
}

//返回排名和排名百分比
func (r *Rank) UpdateScore(id uint64, score int) (int, int) {

	r.cc++

	defer func() {
		if r.cc%100 == 0 {
			r.shrink(nil)
		}
	}()

	rank := 0

	item := r.getRankItem(id)
	if nil == item {
		item = r.itemPool.get()
		r.id2Item[id] = item
		item.value = id
	} else {
		if item.key == score {
			return r.getRank(item), r.getRankPersentByItem(item)
		}
	}

	item.key = score

	c := r.findSpan(score)

	oldC := item.sl

	if item.sl != nil {
		sl := item.sl
		sl.DeleteNode(item)
		sl.fixMinMax()
	}

	rank = c.InsertNode(item)

	if c.size > maxItemCount+maxItemCount/2 {

		if o := c.split(); nil != o {

			c.fixMinMax()
			o.fixMinMax()

			o.idx = c.idx + 1
			if o.idx >= len(r.spans) {
				r.spans = append(r.spans, o)
			} else {
				if len(r.spans) < cap(r.spans) {
					//还有空间,扩张containers,将downIdx开始的元素往后挪一个位置，空出downIdx所在位置
					l := len(r.spans)
					r.spans = r.spans[:len(r.spans)+1]
					for i := l - 1; i >= o.idx; i-- {
						r.spans[i+1] = r.spans[i]
						r.spans[i+1].idx = i + 1
					}
					r.spans[o.idx] = o

				} else {

					//下一个container满了，新建一个
					spans := make([]*skiplists, 0, len(r.spans)+1)
					for i := 0; i <= o.idx-1; i++ {
						spans = append(spans, r.spans[i])
					}

					spans = append(spans, o)

					for i := o.idx; i < len(r.spans); i++ {
						c := r.spans[i]
						c.idx = len(spans)
						spans = append(spans, c)
					}

					r.spans = spans
				}
			}
		} else {
			c.fixMinMax()
		}
	} else {
		c.fixMinMax()
	}

	if nil != oldC {
		if oldC.size == 0 {
			//oldC已经被清空，需要删除
			for i := oldC.idx + 1; i < len(r.spans); i++ {
				c := r.spans[i]
				r.spans[i-1] = c
				c.idx = i - 1
			}

			r.spans[len(r.spans)-1] = nil
			r.spans = r.spans[:len(r.spans)-1]
		} else if oldC.idx != item.sl.idx && oldC.idx+1 < len(r.spans) && oldC.size+r.spans[oldC.idx+1].size <= maxItemCount {
			r.shrink(oldC)
		}
	}

	return rank + r.getFrontSpanItemCount(item), r.getRankPersent(item, rank)
}

func (r *Rank) shrink(s *skiplists) {
	if nil == s {
		if r.nextShink >= len(r.spans)-1 {
			r.nextShink = 0
			return
		} else {
			s = r.spans[r.nextShink]
			r.nextShink++
		}
	}

	if s.idx+1 < len(r.spans) && s.size+r.spans[s.idx+1].size <= maxItemCount {
		n := r.spans[s.idx+1]
		s.merge(n)
		s.fixMinMax()
		for i := n.idx + 1; i < len(r.spans); i++ {
			c := r.spans[i]
			r.spans[i-1] = c
			c.idx = i - 1
		}
		r.spans[len(r.spans)-1] = nil
		r.spans = r.spans[:len(r.spans)-1]
	}
}
