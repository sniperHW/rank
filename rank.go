package rank

import (
	"fmt"
)

var maxItemCount int = 100
var searchSetp int

type rankItemBlock struct {
	items    []rankItem
	nextFree int
}

func (rb *rankItemBlock) get() *rankItem {
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
		items: make([]rankItem, 100000),
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

func (rp *rankItemPool) get() *rankItem {
	item := rp.blocks[rp.nextFree].get()
	if nil == item {
		block := newRankItemBlock()
		rp.blocks = append(rp.blocks, block)
		rp.nextFree++
		item = block.get()
	}

	item.c = nil
	return item
}

type rankItem struct {
	id    uint64
	score int
	pprev *rankItem
	pnext *rankItem
	c     *span
}

type span struct {
	idx   int
	max   int
	min   int
	count int
	head  rankItem
	tail  rankItem
}

func newSpan(idx int) *span {
	c := &span{
		idx: idx,
	}

	c.head.pnext = &c.tail
	c.tail.pprev = &c.head

	return c
}

func (c *span) show() {
	cur := c.head.pnext
	for cur != &c.tail {
		fmt.Println(cur.id, cur.score)
		cur = cur.pnext
	}
}

func (c *span) remove(item *rankItem) {
	c.count--
	p := item.pprev
	n := item.pnext

	p.pnext = n
	n.pprev = p

	//item.pprev = nil
	//item.pnext = nil
	//item.c = nil

	c.fixMinMax()
}

func (c *span) update(item *rankItem, change int) {

	p := item.pprev
	n := item.pnext

	var cc *rankItem

	if change > 0 {
		//积分增加往前移动
		if p == &c.head {
			return
		}

		p.pnext = n
		n.pprev = p

		for p != &c.head {
			if p.score >= item.score {
				break
			} else {
				searchSetp++
				p = p.pprev
			}
		}

		cc = p

		//插入到cc后面

		n = cc.pnext
		n.pprev = item
		cc.pnext = item
		item.pnext = n
		item.pprev = cc

	} else {
		//积分减少往后移
		if n == &c.tail {
			return
		}

		p.pnext = n
		n.pprev = p

		for n != &c.tail {
			if n.score <= item.score {
				break
			} else {
				searchSetp++
				n = n.pnext
			}
		}

		cc = n

		//插入到cc前
		p = cc.pprev
		p.pnext = item
		cc.pprev = item
		item.pprev = p
		item.pnext = cc

	}

	c.fixMinMax()
}

func (c *span) down(item *rankItem) *rankItem {
	c.count++
	item.c = c

	n := c.head.pnext

	c.head.pnext = item
	item.pprev = &c.head
	n.pprev = item
	item.pnext = n

	var r *rankItem

	if c.count > maxItemCount {
		c.count--
		//最后一个元素
		r = c.tail.pprev
		r.pprev.pnext = &c.tail
		c.tail.pprev = r.pprev
	}

	c.fixMinMax()
	return r

}

func (c *span) add(item *rankItem) *rankItem {

	c.count++
	item.c = c
	//寻找合适的插入位置
	var cc *rankItem

	front := c.head.pnext
	back := c.tail.pprev
	for {
		if front.score <= item.score {
			cc = front
			break
		} else if back.score >= item.score {
			cc = back.pnext
			break
		} else if back.pprev.score >= item.score && back.score <= item.score {
			cc = back
			break
		}
		front = front.pnext
		back = back.pprev
	}

	//插入到cc前
	p := cc.pprev
	p.pnext = item
	cc.pprev = item
	item.pprev = p
	item.pnext = cc

	var r *rankItem

	if c.count > maxItemCount {
		c.count--
		//最后一个元素
		r = c.tail.pprev
		r.pprev.pnext = &c.tail
		c.tail.pprev = r.pprev
		//r.pprev = nil
		//r.pnext = nil
	}

	c.fixMinMax()

	return r
}

func (c *span) fixMinMax() {
	if c.count > 0 {
		c.max = c.head.pnext.score
		c.min = c.tail.pprev.score
	}
}

func (c *span) check(max int) int {
	cur := c.head.pnext
	for cur != &c.tail {
		if !(max >= cur.score) {
			return -1
		} else {
			max = cur.score
		}
		cur = cur.pnext
	}
	return max
}

type Rank struct {
	id2Item  map[uint64]*rankItem
	spans    []*span
	itemPool *rankItemPool
}

func NewRank() *Rank {
	return &Rank{
		id2Item:  map[uint64]*rankItem{},
		spans:    make([]*span, 0, 65536),
		itemPool: newRankItemPool(),
	}
}

func (r *Rank) Reset() {
	r.id2Item = map[uint64]*rankItem{}
	r.spans = make([]*span, 0, 65536)
	r.itemPool.reset()
}

func (r *Rank) GetPercentRank(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return 100 - 100*item.c.idx/(len(r.spans)-1)
	}
}

func (r *Rank) GetExactRank(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		c := 0
		for i := 0; i < item.c.idx; i++ {
			c += r.spans[i].count
		}

		cc := item.c
		cur := item
		for cur != &cc.head {
			c++
			cur = cur.pprev
		}
		return c
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
		fmt.Println("----------", v.idx, "----------------")
		v.show()
	}
}

func (r *Rank) getRankItem(id uint64) *rankItem {
	return r.id2Item[id]
}

func (r *Rank) binarySearch(score int, left int, right int) *span {
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

func (r *Rank) findSpan(score int) *span {
	var c *span
	if len(r.spans) == 0 {
		c = newSpan(len(r.spans))
		r.spans = append(r.spans, c)
	} else {
		c = r.binarySearch(score, 0, len(r.spans)-1)
	}

	return c
}

func (r *Rank) UpdateScore(id uint64, score int) {
	var change int
	item := r.getRankItem(id)
	if nil == item {
		item = r.itemPool.get()
		item.id = id
		item.score = score

		r.id2Item[id] = item
	} else {
		if item.score == score {
			return
		}

		if item.score > score {
			change = -1
		} else {
			change = 1
		}

		item.score = score
	}

	c := r.findSpan(score)

	if c == item.c {
		c.update(item, change)
	} else {

		oldC := item.c

		if item.c != nil {
			item.c.remove(item)
		}

		if downItem := c.add(item); nil != downItem {

			downCount := 0
			downIdx := c.idx

			for nil != downItem {
				downIdx++
				downCount++
				if downCount >= 15 {
					if downIdx >= len(r.spans) {
						r.spans = append(r.spans, newSpan(downIdx))
					} else if r.spans[downIdx].count >= maxItemCount {

						if len(r.spans) < cap(r.spans) {
							//还有空间,扩张containers,将downIdx开始的元素往后挪一个位置，空出downIdx所在位置
							l := len(r.spans)
							r.spans = r.spans[:len(r.spans)+1]
							for i := l - 1; i >= downIdx; i-- {
								r.spans[i+1] = r.spans[i]
								r.spans[i+1].idx = i + 1
							}
							r.spans[downIdx] = newSpan(downIdx)

						} else {
							//下一个container满了，新建一个
							spans := make([]*span, 0, len(r.spans)+1)
							for i := 0; i <= c.idx; i++ {
								spans = append(spans, r.spans[i])
							}

							spans = append(spans, newSpan(len(spans)))

							for i := downIdx; i < len(r.spans); i++ {
								c = r.spans[i]
								c.idx = len(spans)
								spans = append(spans, c)
							}

							r.spans = spans
						}
					}
				} else {
					if downIdx >= len(r.spans) {
						r.spans = append(r.spans, newSpan(downIdx))
					}
				}

				downItem = r.spans[downIdx].down(downItem)
			}
		}

		if nil != oldC && oldC.count == 0 {
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
}
