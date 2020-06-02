package rank

import (
	"fmt"
	"math"
	"strings"
)

const maxItemCount int = 100
const maxLevel int = 2

type link struct {
	pnext    *node
	pprev    *node
	nextSkip int //向前跳过多少个节点
	prevSkip int //向后跳过多少个节点
}

type node struct {
	id    uint64
	score int
	links [maxLevel]link
	sl    *skiplist
}

type skiplist struct {
	idx       int
	head      node
	tail      node
	nodeCount int
	max       int
	min       int
}

func newSkipList(idx int) *skiplist {
	sk := &skiplist{
		idx: idx,
	}

	for i := 0; i < maxLevel; i++ {
		sk.head.links[i].pnext = &sk.tail
		sk.tail.links[i].pprev = &sk.head
	}

	return sk
}

func (sl *skiplist) check(max int) int {
	cc := 0
	cur := sl.head.links[0].pnext
	for cur != &sl.tail {
		if !(max >= cur.score) {
			fmt.Println(sl.idx, "check failed 1")
			return -1
		} else {
			max = cur.score
		}
		cur = cur.links[0].pnext
		cc++
	}

	if cc != sl.nodeCount {
		fmt.Println(sl.idx, "check failed 2", cc, sl.nodeCount)
		return -1
	}

	return max
}

func (sl *skiplist) fixMinMax() {
	if sl.nodeCount > 0 {
		sl.max = sl.head.links[0].pnext.score
		sl.min = sl.tail.links[0].pprev.score
	} else {
		sl.max = 0
		sl.min = 0
	}
}

func (sl *skiplist) getToplinkCount() int {
	n := sl.head.links[maxLevel-1].pnext
	c := 0
	for n != &sl.tail {
		c++
		n = n.links[maxLevel-1].pnext
	}
	return c
}

func (sl *skiplist) show() {
	fmt.Println("--------------------idx", sl.idx, "--------------------------------")
	for i := 0; i < maxLevel; i++ {
		fmt.Println("level", i+1)
		n := sl.head.links[i].pnext
		s := []string{}
		c := 0
		for n != &sl.tail {
			c++
			//s = append(s, fmt.Sprintf("(score:%d,prevSkip:%d,nextSkip:%d)", n.score, n.links[i].prevSkip, n.links[i].nextSkip))
			s = append(s, fmt.Sprintf("(id:%d,score:%d)", n.id, n.score))
			n = n.links[i].pnext
		}
		fmt.Println(strings.Join(s, "->"), "link size", c)
	}
}

func (sl *skiplist) getRank(n *node) int {
	l := maxLevel - 1
	for {
		if n.links[l].pprev == nil {
			l--
		} else {
			break
		}
	}
	return sl.get_rank(l, n)
}

func (sl *skiplist) get_rank(l int, n *node) int {
	c := 0
	cur := n
	for {
		p := cur.links[l].pprev
		c += cur.links[l].prevSkip
		if p == &sl.head {
			return c + 1
		} else if l+1 < maxLevel && p.links[l+1].pprev != nil {
			return c + sl.get_rank(l+1, p)
		} else {
			cur = p
		}
	}
	return c
}

//弹出首节点
func (sl *skiplist) pop_front() *node {
	if sl.nodeCount == 0 {
		return nil
	} else {
		f := sl.head.links[0].pnext
		l0n := f.links[0].pnext
		for i := 0; i < maxLevel; i++ {
			n := f.links[i].pnext
			sl.head.links[i].pnext = n
			n.links[i].pprev = &sl.head
			if i == 0 {
				if n != &sl.tail {
					n.links[i].prevSkip--
				}
			} else {
				if l0n != &sl.tail && l0n != n {

					sl.head.links[i].pnext = l0n
					l0n.links[i].pprev = &sl.head

					l0n.links[i].pnext = n
					n.links[i].pprev = l0n
					n.links[i].prevSkip--
					l0n.links[i].nextSkip = n.links[i].prevSkip
				}
			}
			f.links[i].pnext, f.links[i].pprev = nil, nil
			f.links[i].nextSkip, f.links[i].prevSkip = 0, 0
		}
		sl.nodeCount--
		f.sl = nil
		sl.fixMinMax()
		return f
	}
}

//弹出尾节点
func (sl *skiplist) pop_back() *node {
	if sl.nodeCount == 0 {
		return nil
	} else {
		b := sl.tail.links[0].pprev
		l0p := b.links[0].pprev
		for i := 0; i < maxLevel; i++ {
			p := b.links[i].pprev
			sl.tail.links[i].pprev = p
			p.links[i].pnext = &sl.tail

			if i == 0 {
				if p != &sl.head {
					p.links[i].nextSkip--
				}
			} else {
				if l0p != &sl.head && l0p != p {

					sl.tail.links[i].pprev = l0p
					l0p.links[i].pnext = &sl.tail

					l0p.links[i].pprev = p
					p.links[i].pnext = l0p
					p.links[i].nextSkip--
					l0p.links[i].prevSkip = p.links[i].nextSkip
				}
			}

			b.links[i].pnext, b.links[i].pprev = nil, nil
			b.links[i].nextSkip, b.links[i].prevSkip = 0, 0
		}

		sl.nodeCount--
		b.sl = nil
		sl.fixMinMax()
		return b
	}
}

//向尾巴追加元素，如果sl非空,必须保证n->score<=尾部元素
func (sl *skiplist) push_back(n *node) {
	last := sl.tail.links[0].pprev
	if last == &sl.head {
		for i := 0; i < maxLevel; i++ {
			prev := &sl.head
			next := &sl.tail

			prev.links[i].pnext = n
			n.links[i].pprev = prev
			n.links[i].prevSkip = 0

			next.links[i].pprev = n
			n.links[i].pnext = next
			n.links[i].nextSkip = 0
		}
	} else {
		for i := 0; i < maxLevel; i++ {
			next := &sl.tail
			maxSkip := int(math.Pow10(i))

			if last.links[i].pprev != &sl.head && last.links[i].prevSkip < maxSkip {

				prev := last.links[i].pprev

				prev.links[i].nextSkip++
				n.links[i].prevSkip = prev.links[i].nextSkip
				prev.links[i].pnext = n
				n.links[i].pprev = prev

				last.links[i].pprev = nil
				last.links[i].pnext = nil

			} else {
				last.links[i].nextSkip = 1
				last.links[i].pnext = n
				n.links[i].pprev = last
				n.links[i].prevSkip = 1
			}

			n.links[i].nextSkip = 0
			n.links[i].pnext = next
			next.links[i].pprev = n
		}
	}
	sl.nodeCount++
	n.sl = sl
	sl.fixMinMax()
}

//向头部添加元素，如果sl非空,必须保证n->score>=头部元素
func (sl *skiplist) push_front(n *node) {
	first := sl.head.links[0].pnext
	if first == &sl.tail {
		for i := 0; i < maxLevel; i++ {
			prev := &sl.head
			next := &sl.tail

			prev.links[i].pnext = n
			n.links[i].pprev = prev
			n.links[i].prevSkip = 0

			next.links[i].pprev = n
			n.links[i].pnext = next
			n.links[i].nextSkip = 0
		}
	} else {
		for i := 0; i < maxLevel; i++ {
			prev := &sl.head
			maxSkip := int(math.Pow10(i))

			if first.links[i].pnext != &sl.tail && first.links[i].nextSkip < maxSkip {
				next := first.links[i].pnext

				next.links[i].prevSkip++
				n.links[i].nextSkip = next.links[i].prevSkip
				next.links[i].pprev = n
				n.links[i].pnext = next

				first.links[i].pprev = nil
				first.links[i].pnext = nil

			} else {
				first.links[i].prevSkip = 1
				first.links[i].pprev = n
				n.links[i].pnext = first
				n.links[i].nextSkip = 1
			}

			n.links[i].prevSkip = 0
			n.links[i].pprev = prev
			prev.links[i].pnext = n
		}
	}
	sl.nodeCount++
	n.sl = sl
	sl.fixMinMax()
}

//返回n在level的prev及next节点
func (sl *skiplist) find(head *node, n *node, level int) (pprev *node, pnext *node) {
	if level < 0 {
		return
	} else {
		next := head.links[level].pnext
		if nil == next {
			nn := head.links[0].pnext
			for nn != &sl.tail {
				if nn.links[level].pnext != nil {
					next = nn
					break
				} else {
					nn = nn.links[0].pnext
				}
			}

			if nil == next {
				return nil, nil
			}
		}

		for {
			if n.score >= next.score {
				pprev, pnext = next.links[level].pprev, next
				break
			} else {
				next = next.links[level].pnext
			}
		}
		return
	}
}

func (sl *skiplist) add(n *node) *node {
	sl.insert(n)
	if sl.nodeCount > maxItemCount {
		return sl.pop_back()
	} else {
		return nil
	}
}

func (sl *skiplist) down(n *node) *node {
	sl.push_front(n)
	if sl.nodeCount > maxItemCount {
		return sl.pop_back()
	} else {
		return nil
	}
}

func (sl *skiplist) merge(o *skiplist) {
	for sl.nodeCount < maxItemCount {
		n := o.pop_front()
		if nil == n {
			return
		} else {
			sl.push_back(n)
		}
	}
}

func (sl *skiplist) insert(n *node) {
	prev_nexts := [maxLevel][2]*node{}
	heads := [maxLevel]*node{}
	heads[maxLevel-1] = &sl.head
	for i := maxLevel - 1; i >= 0; i-- {
		pprev, pnext := sl.find(heads[i], n, i)
		prev_nexts[i][0] = pprev
		prev_nexts[i][1] = pnext
		if i > 0 {
			heads[i-1] = pprev
		}
	}

	if prev_nexts[0][0] == &sl.head {
		sl.push_front(n)
	} else if prev_nexts[0][1] == &sl.tail {
		sl.push_back(n)
	} else {
		//从level 0开始，将n加入
		for i := 0; i < maxLevel; i++ {
			prev := prev_nexts[i][0]
			next := prev_nexts[i][1]
			maxSkip := int(math.Pow10(i))

			if i == 0 {
				prev.links[i].pnext = n
				n.links[i].pprev = prev
				n.links[i].prevSkip = prev.links[i].nextSkip
				n.links[i].pnext = next
				next.links[i].pprev = n
				n.links[i].nextSkip, next.links[i].prevSkip = 1, 1
			} else {
				if prev.links[i].nextSkip+1 >= maxSkip {
					np := next.links[0].pprev
					nn := next.links[i].pnext

					prev.links[i].pnext = np
					np.links[i].pprev = prev

					np.links[i].prevSkip = prev.links[i].nextSkip
					if nil != nn && nn != &sl.tail && nn.links[i].prevSkip < maxSkip {

						np.links[i].pnext = nn
						nn.links[i].pprev = np

						nn.links[i].prevSkip++
						np.links[i].nextSkip = nn.links[i].prevSkip

						next.links[i].pnext, next.links[i].pprev = nil, nil

					} else {
						np.links[i].pnext = next
						next.links[i].pprev = np
						np.links[i].nextSkip, next.links[i].prevSkip = 1, 1
					}

				} else {
					prev.links[i].nextSkip++
					next.links[i].prevSkip++
					n.links[i].pprev, n.links[i].pnext = nil, nil
				}
			}
		}
		n.sl = sl
		sl.nodeCount++
		sl.fixMinMax()
	}
}

func (sl *skiplist) getNext(l int, n *node) *node {
	if n.links[l].pnext != nil || n.links[l].pprev != nil {
		return n
	} else {
		pnext := n.links[l-1].pnext
		for pnext != &sl.tail {
			if pnext.links[l].pprev != nil || pnext.links[l].pnext != nil {
				return pnext
			} else {
				pnext = pnext.links[l-1].pnext
			}
		}

		return nil
	}
}

func (sl *skiplist) getPre(l int, n *node) *node {

	if n.links[l].pnext != nil || n.links[l].pprev != nil {
		return n
	} else {
		pprev := n.links[l-1].pprev
		for pprev != &sl.head {
			if pprev.links[l].pprev != nil || pprev.links[l].pnext != nil {
				return pprev
			} else {
				pprev = pprev.links[l-1].pprev
			}
		}
		return nil
	}
}

func (sl *skiplist) remove(n *node) {
	if sl.nodeCount == 0 {
		panic("nodeCount == 0")
	} else if sl.nodeCount == 1 {
		for i := 0; i < maxLevel; i++ {
			sl.head.links[i].pnext = &sl.tail
			sl.tail.links[i].pprev = &sl.head
			n.links[i].pprev, n.links[i].pnext = nil, nil
			n.links[i].nextSkip, n.links[i].prevSkip = 0, 0
		}
		sl.nodeCount--
		sl.fixMinMax()
	} else {
		if n == sl.head.links[0].pnext {
			sl.pop_front()
		} else if n == sl.tail.links[0].pprev {
			sl.pop_back()
		} else {
			prev_nexts := [maxLevel][2]*node{}
			prev_nexts[0][0] = n.links[0].pprev
			prev_nexts[0][1] = n.links[0].pnext

			for i := 1; i < maxLevel; i++ {
				if n.links[i].pprev != nil || n.links[i].pnext != nil {
					prev_nexts[i][0] = n.links[i].pprev
					prev_nexts[i][1] = n.links[i].pnext
				} else {
					prev_nexts[i][0] = sl.getPre(i, prev_nexts[i-1][0])
					prev_nexts[i][1] = sl.getNext(i, prev_nexts[i-1][1])
				}
			}

			l0n := n.links[0].pnext

			for i := 0; i < maxLevel; i++ {
				pprev, pnext := prev_nexts[i][0], prev_nexts[i][1]

				if i == 0 {
					pprev.links[i].pnext = pnext
					pnext.links[i].pprev = pprev
				} else {
					if n.links[i].pprev == nil {
						//在pprev和pnext之间的节点，但不属于i级链
						pprev.links[i].nextSkip--
						pnext.links[i].prevSkip = pprev.links[i].nextSkip
					} else {

						//i级链上的节点
						pprev.links[i].pnext = l0n
						l0n.links[i].pprev = pprev
						l0n.links[i].prevSkip = pprev.links[i].nextSkip

						pnext.links[i].prevSkip--
						if l0n != pnext {
							l0n.links[i].pnext = pnext
							pnext.links[i].pprev = l0n
							l0n.links[i].nextSkip = pnext.links[i].prevSkip

						}
					}
				}

				n.links[i].pprev, n.links[i].pnext = nil, nil
				n.links[i].nextSkip, n.links[i].prevSkip = 0, 0
			}

			sl.nodeCount--
			n.sl = nil
			sl.fixMinMax()
		}
	}
}
