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
	value int
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
		sl.max = sl.head.links[0].pnext.value
		sl.min = sl.tail.links[0].pprev.value
	} else {
		sl.max = 0
		sl.min = 0
	}
}

func (sl *skiplists) show() {

	fmt.Println("max", sl.max, "min", sl.min, sl.head.links[0].pnext.value, sl.tail.links[0].pprev.value)

	for i := 0; i <= sl.level; i++ {
		cur := sl.head.links[i].pnext
		s := []string{}
		s = append(s, fmt.Sprintf("head skip:%d", sl.head.links[i].skip))
		for &sl.tail != cur {
			s = append(s, fmt.Sprintf("(%d,%d,%d)", cur.key, cur.value, cur.links[i].skip))
			cur = cur.links[i].pnext
		}
		fmt.Println("level", i+1, strings.Join(s, ","))
	}
	fmt.Println("--------------------------------------------------------------------------")
}

func (sl *skiplists) randomLevel() int {
	lvl := 0
	for rand.Float32() < 0.1 && lvl < maxLevel-1 {
		lvl++
	}
	return lvl
}

func (sl *skiplists) check(v int) int {
	c := 0
	cur := sl.head.links[0].pnext
	for &sl.tail != cur {
		if cur.value > v {
			return -1
		}
		v = cur.value
		c++
		cur = cur.links[0].pnext
	}

	if c != sl.size {
		//fmt.Println("check failed2", c, sl.size)
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

func (sl *skiplists) insert(n *node, update [maxLevel]*node, updateOffset [maxLevel]int) int {

	offset0 := updateOffset[0] + 1 //新节点在level1的位置

	lvl := sl.randomLevel()
	if lvl > sl.level {
		for i := sl.level + 1; i <= lvl; i++ {
			update[i] = &sl.head
			updateOffset[i] = 1
		}
		sl.level = lvl
	}

	x := n

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
					sl.show()
					panic("3")
				}
			} else {
				x.links[i].skip = 0
			}

		} else {

			if update[i].links[i].pnext != &sl.tail {
				update[i].links[i].skip++
			}

		}
	}

	sl.size++
	n.sl = sl

	return offset0 - 1
}

func (sl *skiplists) InsertFront(n *node) {

	update := [maxLevel]*node{}     //n插入到update后面，n的offset=updateOffset[0] + 1
	updateOffset := [maxLevel]int{} //update节点的offset

	for i := 0; i < maxLevel; i++ {
		update[i] = &sl.head
		updateOffset[i] = 1
	}

	sl.insert(n, update, updateOffset)

}

//由上层确保n不在skiplists里
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

		for &sl.tail != x.links[i].pnext && x.links[i].pnext.key < key {
			updateOffset[i] += x.links[i].skip
			x = x.links[i].pnext
			sl.step++
		}

		update[i] = x
	}

	return sl.insert(n, update, updateOffset)
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

	for x != nil {
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
	}

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

	for sl.level >= 0 && sl.head.links[sl.level].pnext == &sl.tail {
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
