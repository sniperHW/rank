package rank

import (
//"fmt"
)

var realRankIdx int = 100

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
	rb.nextFree = 0
}

func newRankItemBlock() *rankItemBlock {
	return &rankItemBlock{
		items: make([]node, 100000),
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
	spans     []*skiplist
	itemPool  *rankItemPool
	nextShink int
	cc        int
}

func NewRank() *Rank {
	return &Rank{
		id2Item:  map[uint64]*node{},
		spans:    make([]*skiplist, 0, 65536/2),
		itemPool: newRankItemPool(),
	}
}

func (r *Rank) Reset() {
	r.id2Item = map[uint64]*node{}
	r.spans = make([]*skiplist, 0, 65536/2)
	r.itemPool.reset()
}

func (r *Rank) GetPercentRank(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return 100 - 100*item.sl.idx/(len(r.spans)-1)
	}
}

func (r *Rank) getFrontSpanItemCount(item *node) int {
	c := 0
	if item.sl.idx < realRankIdx {
		for i := 0; i < item.sl.idx; i++ {
			c += r.spans[i].nodeCount
		}
	} else {
		c = 100 * item.sl.idx
	}
	return c
}

func (r *Rank) getExactRank(item *node) int {
	return r.getFrontSpanItemCount(item) + item.sl.getRank(item)
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
	if len(r.spans) > 0 {
		max := r.spans[0].max
		for _, v := range r.spans {
			max = v.check(max)
			if max == -1 {
				return false
			}
		}
	}
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

func (r *Rank) binarySearch(score int, left int, right int) *skiplist {

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

func (r *Rank) findSpan(score int) *skiplist {
	var c *skiplist
	if len(r.spans) == 0 {
		c = newSkipList(len(r.spans))
		r.spans = append(r.spans, c)
	} else {
		c = r.binarySearch(score, 0, len(r.spans)-1)
	}

	return c
}

func (r *Rank) UpdateScore(id uint64, score int) int {

	r.cc++

	//defer func() {
	//	if r.cc%100 == 0 {
	//		r.shrink(10)
	//	}
	//}()

	//var realRank int
	item := r.getRankItem(id)
	if nil == item {
		item = r.itemPool.get()
		item.id = id
		item.score = score

		r.id2Item[id] = item
	} else {
		if item.score == score {
			return r.getExactRank(item)
		}

		item.score = score
	}

	if item.sl != nil && item.sl.max > score && item.sl.min <= score {
		sl := item.sl
		sl.remove(item)
		sl.insert(item)
	} else {
		c := r.findSpan(score)

		oldC := item.sl

		if item.sl != nil {
			item.sl.remove(item)
		}

		if downItem := c.add(item); nil != downItem {
			downCount := 0
			downIdx := c.idx

			for nil != downItem {
				downIdx++
				downCount++
				if downIdx < realRankIdx || downCount <= 15 {
					if downIdx >= len(r.spans) {
						r.spans = append(r.spans, newSkipList(downIdx))
					}
				} else {
					//超过down次数，创建一个新的span接纳下降item终止下降过程
					if downIdx >= len(r.spans) {
						r.spans = append(r.spans, newSkipList(downIdx))

					} else if r.spans[downIdx].nodeCount >= maxItemCount {

						if len(r.spans) < cap(r.spans) {

							//还有空间,扩张containers,将downIdx开始的元素往后挪一个位置，空出downIdx所在位置
							l := len(r.spans)
							r.spans = r.spans[:len(r.spans)+1]
							for i := l - 1; i >= downIdx; i-- {
								r.spans[i+1] = r.spans[i]
								r.spans[i+1].idx = i + 1
							}
							r.spans[downIdx] = newSkipList(downIdx)

						} else {

							//下一个container满了，新建一个
							spans := make([]*skiplist, 0, len(r.spans)+1)
							for i := 0; i <= downIdx-1; i++ {
								spans = append(spans, r.spans[i])
							}

							spans = append(spans, newSkipList(len(spans)))

							for i := downIdx; i < len(r.spans); i++ {
								c := r.spans[i]
								c.idx = len(spans)
								spans = append(spans, c)
							}

							r.spans = spans
						}
					}
				}

				downItem = r.spans[downIdx].down(downItem)

			}
		}

		if nil != oldC && oldC.nodeCount == 0 {
			//oldC已经被清空，需要删除
			for i := oldC.idx + 1; i < len(r.spans); i++ {
				c := r.spans[i]
				r.spans[i-1] = c
				c.idx = i - 1
			}

			r.spans[len(r.spans)-1] = nil
			r.spans = r.spans[:len(r.spans)-1]
		}
	}

	return item.sl.getRank(item) + r.getFrontSpanItemCount(item)
}

func (r *Rank) shrink(emptyCount int) {
	if r.nextShink >= len(r.spans)-1 {
		r.nextShink = 0
	} else {
		s := r.spans[r.nextShink]
		if maxItemCount-s.nodeCount > emptyCount {
			//如果当前span有空间，将后续span的元素吸纳进当前span
			n := r.spans[r.nextShink+1]
			s.merge(n)
			if n.nodeCount == 0 {
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
		r.nextShink++
	}
}
