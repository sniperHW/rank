package rank

import (
	"fmt"
	"math/rand"
	"strings"
)

const maxLevel int = 15

type link struct {
	pnext *node
	pprev *node
	skip  int //next跳过了多少个节点
}

type node struct {
	key   int
	value uint64
	links [maxLevel]link
	sl    *skiplists
}

type skiplists struct {
	level int
	size  int
	head  node
	tail  node
	step  int
	max   int
	min   int
	idx   int
}

func newSkipLists(idx int) *skiplists {
	sl := &skiplists{
		idx: idx,
	}

	for i := 0; i < maxLevel; i++ {
		sl.head.links[i].pnext = &sl.tail
		sl.tail.links[i].pprev = &sl.head
	}

	return sl
}

func (sl *skiplists) fixMinMax() {
	if sl.size > 0 {
		sl.max = sl.head.links[0].pnext.key
		sl.min = sl.tail.links[0].pprev.key
	} else {
		sl.max = 0
		sl.min = 0
	}
}

func (sl *skiplists) show() {

	fmt.Println("size", sl.size, "max", sl.max, "min", sl.min)

	for i := 0; i <= sl.level; i++ {
		cur := sl.head.links[i].pnext
		s := []string{}
		s = append(s, fmt.Sprintf("head skip:%d", sl.head.links[i].skip))
		for &sl.tail != cur {
			if nil == cur.links[i].pnext {
				fmt.Println(cur.key, cur.value)
			}

			s = append(s, fmt.Sprintf("(key:%d,value:%d,skip:%d)", cur.key, cur.value, cur.links[i].skip))
			cur = cur.links[i].pnext
		}
		fmt.Println("level", i+1, strings.Join(s, ","))
	}
	fmt.Println("--------------------------------------------------------------------------")
}

func (sl *skiplists) randomLevel() int {
	lvl := 0
	for rand.Float32() < 0.55 && lvl < maxLevel-1 {
		lvl++
	}
	return lvl
}

func (sl *skiplists) checkLink() bool {

	for i := 0; i < maxLevel; i++ {
		if i > sl.level && sl.head.links[i].pnext != &sl.tail {
			panic("check head failed")
		}
	}

	a := make([]*node, 0, sl.size)
	cur := sl.head.links[0].pnext
	for &sl.tail != cur {
		a = append(a, cur)
		cur = cur.links[0].pnext
	}

	for i := 0; i <= sl.level; i++ {

		tailPre := sl.tail.links[i].pprev
		if tailPre.links[i].pnext != &sl.tail {
			fmt.Println("check tail false level", i+1)
			return false
		}

		cur := &sl.head
		idx := 0
		for cur.links[i].pnext != &sl.tail {
			n := cur.links[i].pnext
			idx += cur.links[i].skip
			if idx-1 >= sl.size {
				fmt.Println(i, cur.key)
				panic("xxxxxx")
			}

			if n != a[idx-1] {
				fmt.Println("level", i, "idx", idx, n.key, a[idx-1].key)
				sl.show()
				panic("xxxxxxxxxxx1")
				return false
			}

			cur = cur.links[i].pnext
		}

		if cur != tailPre {
			sl.show()
			fmt.Println("cur", cur.key, tailPre.key)
			panic("cur != tailPre")
		}
	}

	return true
}

func (sl *skiplists) check(v int) int {

	if !sl.checkLink() {
		return -1
	}

	c := 0
	cur := sl.head.links[0].pnext
	for &sl.tail != cur {
		if cur.key > v {
			return -1
		}
		v = cur.key
		c++
		cur = cur.links[0].pnext
	}

	if c != sl.size {
		return -1
	} else {
		return v
	}
}

/*

func (sl *skiplists) Search(key int) *node {
	x := &sl.head
	for i := sl.level; i >= 0; i-- {
		for nil != x.links[i].pnext && x.links[i].pnext.key < key {
			x = x.links[i].pnext
		}
	}

	x = x.links[0].pnext

	if nil != x && x.key == key {
		return x
	} else {
		return nil
	}
}*/

/////////////////////////使用node的接口

func (sl *skiplists) InsertNode(n *node) int {

	key := n.key

	update := [maxLevel]*node{}     //n插入到update后面，n的offset=updateOffset[0] + 1
	updateOffset := [maxLevel]int{} //update节点的offset
	updateOffset[sl.level] = 1

	x := &sl.head
	for i := sl.level; i >= 0; i-- {
		if i != sl.level {
			updateOffset[i] = updateOffset[i+1]
		}

		for &sl.tail != x.links[i].pnext && x.links[i].pnext.key > key {
			updateOffset[i] += x.links[i].skip
			x = x.links[i].pnext
			sl.step++
		}

		update[i] = x
	}

	offset0 := updateOffset[0] + 1 //新节点在level1的位置

	lvl := sl.randomLevel()
	if lvl > sl.level {
		for i := sl.level + 1; i <= lvl; i++ {
			update[i] = &sl.head
			updateOffset[i] = 1
		}
		sl.level = lvl
	}

	x = n

	for i := 0; i <= sl.level; i++ {

		if i <= lvl {

			x.links[i].pprev = update[i]

			x.links[i].pnext = update[i].links[i].pnext

			update[i].links[i].pnext = x

			x.links[i].pnext.links[i].pprev = x

			oldSkip := update[i].links[i].skip

			update[i].links[i].skip = offset0 - updateOffset[i]

			if x.links[i].pnext != &sl.tail {
				x.links[i].skip = (updateOffset[i] + oldSkip + 1) - offset0
				if x.links[i].skip < 0 {
					fmt.Println("ssss", sl.idx)
					panic("3")
				}
			} else {
				x.links[i].skip = 0
			}

		} else if update[i].links[i].skip > 0 {
			update[i].links[i].skip++
		}
	}

	sl.size++
	n.sl = sl

	/*if !sl.checkLink() {
		panic("InsertNode checkLink failed")
	}*/

	return offset0 - 1
}

func (sl *skiplists) isInLink(lvl int, head *node, n *node) bool {
	cur := head.links[lvl].pnext
	for cur != &sl.tail {
		if cur == n {
			return true
		} else {
			cur = cur.links[lvl].pnext
		}
	}
	return false
}

func (sl *skiplists) DeleteNode(n *node) {

	update := [maxLevel]*node{}
	x := n
	lvl := 0

	for {

		if nil == update[lvl] {
			if x == n {
				update[lvl] = x.links[lvl].pprev
			} else {
				update[lvl] = x
			}
		}

		if lvl+1 > sl.level {
			break
		} else if nil != x.links[lvl+1].pprev || nil != x.links[lvl+1].pnext {
			lvl++
		} else {
			x = x.links[lvl].pprev
		}

	}

	/*for x != nil {
		pprev := x.links[lvl].pprev

		if (x.links[lvl].pprev != nil || x.links[lvl].pnext != nil) && nil == update[lvl] {
			//当前节点是lvl链上的节点
			if x == n {
				update[lvl] = pprev
			} else {
				update[lvl] = x
			}

			if update[lvl].links[lvl].pnext != n && update[lvl].links[lvl].skip == 0 {
				break
			}

		}

		if lvl < sl.level && (nil != x.links[lvl+1].pprev || nil != x.links[lvl+1].pnext) {
			//当前节点也在lvl+1的链上
			lvl++
		} else {
			x = pprev
		}
	}*/

	for i := 0; i <= sl.level; i++ {
		if nil == update[i] {
			break
		} else {

			if update[i].links[i].pnext != n {
				if update[i].links[i].skip > 0 {
					update[i].links[i].skip--
					if update[i].links[i].skip < 1 {
						panic("1")
					}
				} else {
					break
				}
			} else {
				update[i].links[i].pnext = n.links[i].pnext
				n.links[i].pnext.links[i].pprev = update[i]
				if update[i].links[i].pnext != &sl.tail {
					update[i].links[i].skip = update[i].links[i].skip + n.links[i].skip - 1
					if update[i].links[i].skip < 0 {
						panic("2")
					}
				} else {
					update[i].links[i].skip = 0
				}
			}
		}
		n.links[i].pnext = nil
		n.links[i].pprev = nil
		n.links[i].skip = 0
	}

	for sl.level > 0 && sl.head.links[sl.level].pnext == &sl.tail {
		sl.level--
	}

	sl.size--
	n.sl = nil
}

func (sl *skiplists) GetNodeRank(n *node) int {
	rank := 0
	x := n
	var pprev *node
	var lvl int

	for pprev != &sl.head {
		for i := sl.level; i >= 0; i-- {
			pprev = x.links[i].pprev
			if nil != pprev {
				lvl = i
				break
			}
		}
		rank += pprev.links[lvl].skip
		x = pprev

	}

	return rank
}

func (sl *skiplists) GetNodeRankCheck(n *node) int {

	rank := 1

	cur := sl.head.links[0].pnext
	for cur != &sl.tail {
		if cur == n {
			break
		} else {
			rank++
			cur = cur.links[0].pnext
		}
	}

	if cur == &sl.tail {
		return -1
	}

	/*x := n
	var pprev *node
	var lvl int

	for pprev != &sl.head {
		for i := sl.level; i >= 0; i-- {
			pprev = x.links[i].pprev
			if nil != pprev {
				lvl = i
				break
			}
		}
		rank += pprev.links[lvl].skip
		x = pprev

	}*/
	return rank
}

//将o合并到sl
func (sl *skiplists) merge(o *skiplists) {

	if o.size > 0 {

		//将o的第一个节点提升为最高等级节点
		oF := o.head.links[0].pnext
		for i := 0; i <= o.level; i++ {
			if o.head.links[i].pnext != oF {
				oF.links[i].pnext = o.head.links[i].pnext
				o.head.links[i].pnext.links[i].pprev = oF
				o.head.links[i].pnext = oF
				oF.links[i].pprev = &o.head
				oF.links[i].skip = o.head.links[i].skip - 1
				o.head.links[i].skip = 1
			}
		}

		maxL := sl.level
		minL := sl.level
		if maxL < o.level {
			maxL = o.level
		}
		if minL > o.level {
			minL = o.level
		}

		for i := maxL; i >= 0; i-- {

			if i <= minL {
				last := sl.tail.links[i].pprev
				skip := 1

				if i > 0 && last != sl.tail.links[0].pprev {
					lv := i - 1
					cur := last
					for lv >= 0 {
						skip += cur.links[lv].skip
						if &sl.tail != cur.links[lv].pnext {
							cur = cur.links[lv].pnext
						} else {
							lv--
						}
					}
				}

				last.links[i].pnext = oF
				oF.links[i].pprev = last
				last.links[i].skip = skip
			} else {
				if i > sl.level {
					sl.head.links[i].pnext = oF
					oF.links[i].pprev = &sl.head
					sl.head.links[i].skip = sl.size + 1
				}
			}

			//处理tail
			if i <= o.level {
				last := o.tail.links[i].pprev
				last.links[i].pnext = &sl.tail
				sl.tail.links[i].pprev = last
			}

			o.head.links[i].pnext = &o.tail
			o.tail.links[i].pprev = &o.head
			o.head.links[i].skip = 0
		}

		sl.size += o.size
		sl.level = maxL
		o.level = 0
		o.size = 0

		cur := oF
		for cur != &sl.tail {
			cur.sl = sl
			cur = cur.links[0].pnext
		}
	}
}

/*
 *  尽量将sl分成均匀的两段
 *
 *  如果分裂点不是最高level的节点，将其提升为最高level节点
 */
func (sl *skiplists) split() *skiplists {
	half := sl.size / 2
	if half >= 1 {
		update := [maxLevel]*node{}
		updateOffset := [maxLevel]int{}
		c := 0
		cur := &sl.head
		lvl := sl.level
		for {
			cc := c + cur.links[lvl].skip
			if cc > half {
				update[lvl] = cur
				lvl--
				if lvl+1 <= sl.level {
					updateOffset[lvl] = updateOffset[lvl+1]
				}
			} else if cc == half {
				update[lvl] = cur
				cur = cur.links[lvl].pnext
				c = cc
				break
			} else {
				c = cc
				if cur.links[lvl].skip > 0 {
					updateOffset[lvl] += cur.links[lvl].skip
					cur = cur.links[lvl].pnext
				} else {
					update[lvl] = cur
					lvl--
					if lvl+1 <= sl.level {
						updateOffset[lvl] = updateOffset[lvl+1]
					}
				}
			}
		}

		//fmt.Println(c, updateOffset)

		//将cur提升为最高level节点
		for lvl := sl.level; lvl >= 0; lvl-- {
			if update[lvl].links[lvl].pnext == cur {
				break
			} else {
				next := update[lvl].links[lvl].pnext

				next.links[lvl].pprev = cur
				cur.links[lvl].pnext = next

				if update[lvl].links[lvl].skip > 0 {
					nextOffset := updateOffset[lvl] + update[lvl].links[lvl].skip
					cur.links[lvl].skip = nextOffset - c
				} else {
					cur.links[lvl].skip = 0
				}

				update[lvl].links[lvl].pnext = cur
				cur.links[lvl].pprev = update[lvl]
				update[lvl].links[lvl].skip = c - updateOffset[lvl]
			}
		}

		if !sl.checkLink() {
			panic("1234")
		}

		//分裂
		saveLast := [maxLevel]*node{}
		for i := 0; i <= sl.level; i++ {
			saveLast[i] = sl.tail.links[i].pprev

			pre := cur.links[i].pprev
			pre.links[i].skip = 0
			pre.links[i].pnext = &sl.tail
			sl.tail.links[i].pprev = pre
		}

		o := newSkipLists(0)
		o.level = sl.level
		o.size = sl.size - c + 1
		sl.size = c - 1

		for i := 0; i <= o.level; i++ {

			last := saveLast[i]
			o.head.links[i].pnext = cur
			cur.links[i].pprev = &o.head
			o.head.links[i].skip = 1
			last.links[i].pnext = &o.tail
			o.tail.links[i].pprev = last
			last.links[i].skip = 0
		}

		cur = o.head.links[0].pnext
		for cur != &o.tail {
			cur.sl = o
			cur = cur.links[0].pnext
		}

		for sl.level > 0 && sl.head.links[sl.level].pnext == &sl.tail {
			sl.level--
		}

		return o

	} else {
		return nil
	}
}

/*
func (sl *skiplists) split() *skiplists {
	half := sl.size / 2
	if half >= 1 {
		c := 0
		cur := &sl.head

		for {
			if cur.links[sl.level].pnext == &sl.tail {
				break
			} else {
				c += cur.links[sl.level].skip
				cur = cur.links[sl.level].pnext
				if c >= half {
					break
				}
			}
		}

		if cur.links[0].pprev == &sl.head {
			return nil
		}

		saveLast := [maxLevel]*node{}
		for i := 0; i <= sl.level; i++ {
			saveLast[i] = sl.tail.links[i].pprev

			pre := cur.links[i].pprev
			pre.links[i].skip = 0
			pre.links[i].pnext = &sl.tail
			sl.tail.links[i].pprev = pre
		}

		o := newSkipLists(0)
		o.level = sl.level
		o.size = sl.size - c + 1
		sl.size = c - 1

		for i := 0; i <= o.level; i++ {

			last := saveLast[i]
			o.head.links[i].pnext = cur
			cur.links[i].pprev = &o.head
			o.head.links[i].skip = 1
			last.links[i].pnext = &o.tail
			o.tail.links[i].pprev = last
			last.links[i].skip = 0
		}

		cur = o.head.links[0].pnext
		for cur != &o.tail {
			cur.sl = o
			cur = cur.links[0].pnext
		}

		for sl.level > 0 && sl.head.links[sl.level].pnext == &sl.tail {
			sl.level--
		}

		return o
	} else {
		fmt.Println("split failed", sl.size, sl.level)
		return nil
	}
}
*/
