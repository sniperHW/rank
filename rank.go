package rank

import (
//"fmt"
)

const realRankCount int = 10000
const maxItemCount int = 1000
const vacancyRate int = 10 //空缺率10%
const vacancy int = maxItemCount * vacancyRate / 100

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

	item.sl = nil
	return item
}

type Rank struct {
	id2Item   map[uint64]*node
	spans     []*skiplists
	itemPool  *rankItemPool
	nextShink int
	cc        int
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

func (r *Rank) GetPercentRank(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
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

func (r *Rank) getExactRank(item *node) int {
	return r.getFrontSpanItemCount(item) + item.sl.GetNodeRank(item)
}

func (r *Rank) GetExactRank(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return r.getExactRank(item)
	}
}

func (r *Rank) Check() bool {
	return true
}

func (r *Rank) Show() {
	for _, v := range r.spans {
		v.show()
	}
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

func (r *Rank) add(sl *skiplists, item *node) (*node, int) {
	rank := sl.InsertNode(item)
	if sl.size > maxItemCount {
		tail := sl.tail.links[0].pprev
		sl.DeleteNode(tail)
		sl.fixMinMax()
		return tail, rank
	} else {
		return nil, rank
	}
}

func (r *Rank) down(sl *skiplists, item *node) *node {
	sl.InsertFront(item)
	if sl.size > maxItemCount {
		tail := sl.tail.links[0].pprev
		sl.DeleteNode(tail)
		sl.fixMinMax()
		return tail
	} else {
		return nil
	}
}

func (r *Rank) UpdateScore(id uint64, score int) int {

	r.cc++

	defer func() {
		if r.cc%100 == 0 {
			r.shrink(vacancy, nil)
		}
	}()

	rank := 0
	var downItem *node

	item := r.getRankItem(id)
	if nil == item {
		item = r.itemPool.get()
		r.id2Item[id] = item
	} else {
		if item.value == score {
			return r.getExactRank(item)
		}
	}

	item.key = 0 - score
	item.value = score

	if item.sl != nil && item.sl.max > score && item.sl.min <= score {
		sl := item.sl
		sl.DeleteNode(item)
		rank = sl.InsertNode(item)
		sl.fixMinMax()
	} else {
		c := r.findSpan(score)

		oldC := item.sl

		if item.sl != nil {
			sl := item.sl
			sl.DeleteNode(item)
			sl.fixMinMax()
		}

		if downItem, rank = r.add(c, item); nil != downItem {
			downIdx := c.idx
			for nil != downItem {
				downIdx++
				if downIdx >= len(r.spans) {
					r.spans = append(r.spans, newSkipLists(downIdx))

				} else if r.spans[downIdx].size >= maxItemCount {

					if len(r.spans) < cap(r.spans) {

						//还有空间,扩张containers,将downIdx开始的元素往后挪一个位置，空出downIdx所在位置
						l := len(r.spans)
						r.spans = r.spans[:len(r.spans)+1]
						for i := l - 1; i >= downIdx; i-- {
							r.spans[i+1] = r.spans[i]
							r.spans[i+1].idx = i + 1
						}
						r.spans[downIdx] = newSkipLists(downIdx)

					} else {

						//下一个container满了，新建一个
						spans := make([]*skiplists, 0, len(r.spans)+1)
						for i := 0; i <= downIdx-1; i++ {
							spans = append(spans, r.spans[i])
						}

						spans = append(spans, newSkipLists(len(spans)))

						for i := downIdx; i < len(r.spans); i++ {
							c := r.spans[i]
							c.idx = len(spans)
							spans = append(spans, c)
						}

						r.spans = spans
					}
				}
				downItem = r.down(r.spans[downIdx], downItem)
			}
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
			} else if oldC.idx != item.sl.idx && maxItemCount-oldC.size > vacancy {
				r.shrink(vacancy, oldC)
			}
		}
	}

	return rank + r.getFrontSpanItemCount(item)
}

func (r *Rank) merge(to *skiplists, from *skiplists) {
	for to.size < maxItemCount && from.size > 0 {
		item := from.head.links[0].pnext
		from.DeleteNode(item)
		to.InsertNode(item)
	}

	to.fixMinMax()
	from.fixMinMax()
}

func (r *Rank) shrink(vacancy int, s *skiplists) {
	if nil == s {
		if r.nextShink >= len(r.spans)-1 {
			r.nextShink = 0
			return
		} else {
			s = r.spans[r.nextShink]
			r.nextShink++
		}
	}

	if maxItemCount-s.size > vacancy && s.idx+1 < len(r.spans) {
		n := r.spans[s.idx+1]
		r.merge(s, n)
		if n.size == 0 {
			//n已经空了，删除
			for i := n.idx + 1; i < len(r.spans); i++ {
				c := r.spans[i]
				r.spans[i-1] = c
				c.idx = i - 1
			}
			r.spans[len(r.spans)-1] = nil
			r.spans = r.spans[:len(r.spans)-1]
		}
	}
}
