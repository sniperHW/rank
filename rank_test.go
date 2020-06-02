package rank

//go test -covermode=count -v -coverprofile=coverage.out -run=. -cpuprofile=rank.p
//go tool cover -html=coverage.out
//go tool pprof rank.p
//go test -v -run=^$ -bench BenchmarkRank -count 10

import (
	"fmt"
	"github.com/schollz/progressbar"
	"github.com/stretchr/testify/assert"
	"math/rand"
	//"net/http"
	//_ "net/http/pprof"
	"testing"
	"time"
)

/*
func init() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

}*/

func TestSkipList(t *testing.T) {
	/*{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 1; i <= 100; i++ {
			n := &node{
				id:    uint64(i),
				score: i,
			}
			sl.push_back(n)
			nodes = append(nodes, n)
		}
		sl.show()

		for k, v := range nodes {
			assert.Equal(t, k+1, sl.getRank(v))
		}
	}

	{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 100; i >= 1; i-- {
			n := &node{
				id:    uint64(i),
				score: i,
			}
			sl.push_front(n)
			nodes = append(nodes, n)
		}
		sl.show()

		for k, v := range nodes {
			assert.Equal(t, len(nodes)-k, sl.getRank(v))
		}
	}*/

	/*{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 1; i <= 10; i++ {
			n := &node{
				id:    uint64(i),
				score: i,
			}
			sl.insert(n)
			nodes = append(nodes, n)
		}
		sl.show()
	}

	{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 10; i >= 1; i-- {
			n := &node{
				id:    uint64(i),
				score: i,
			}
			sl.insert(n)
			nodes = append(nodes, n)
		}
		sl.show()
	}*/

	/*{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 10; i >= 1; i-- {
			n := &node{
				id:    uint64(i),
				score: i,
			}
			sl.push_front(n)
			nodes = append(nodes, n)
		}
		sl.show()

		sl.pop_front()
		sl.show()

		sl.pop_back()
		sl.show()

	}*/

	/*{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 1; i <= 20; i++ {
			nodes = append(nodes, &node{
				id:    uint64(i),
				score: i,
			})
		}

		sl.push_back(nodes[0])
		sl.push_back(nodes[1])
		sl.push_back(nodes[2])
		sl.show()

		sl.remove(nodes[1])
		sl.show()

	}

	{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 1; i <= 20; i++ {
			nodes = append(nodes, &node{
				id:    uint64(i),
				score: i,
			})
		}

		sl.push_back(nodes[0])
		sl.push_back(nodes[1])
		sl.push_back(nodes[2])
		sl.show()

		sl.remove(nodes[0])
		sl.show()

	}*/

	{
		sl := newSkipList(0)

		nodes := []*node{}

		for i := 1; i <= 20; i++ {
			nodes = append(nodes, &node{
				id:    uint64(i),
				score: i,
			})
			sl.push_back(nodes[i-1])
		}

		//sl.push_back(nodes[0])
		//sl.push_back(nodes[1])
		//sl.push_back(nodes[2])
		sl.show()

		sl.remove(nodes[10])
		sl.show()

		//sl.remove(nodes[0])
		//sl.show()

	}

	/*
		{
			for i := 0; i < 10; i++ {

				sl := newSkipList(0)
				nodes := []*node{}
				for i := 1; i <= 100; i++ {
					n := &node{
						id:    uint64(i),
						score: i,
					}
					nodes = append(nodes, n)
				}

				shuffle := func() {
					i := rand.Int() % len(nodes)
					j := rand.Int() % len(nodes)
					nodes[j], nodes[i] = nodes[i], nodes[j]
				}

				for i := 0; i < 100; i++ {
					shuffle()
				}

				for _, v := range nodes {
					sl.insert(v)
				}

				sl.show()

				for _, v := range nodes {
					assert.Equal(t, len(nodes)-v.score+1, sl.getRank(v))
				}

			}

		}
	*/
}

func TestBenchmarkRank(t *testing.T) {
	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 5000000
	idRange := 5000000

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int() % 1000000
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans), r.getToplinkCount()/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{
		testCount := 10000000
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int() % 1000000
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans), r.getToplinkCount()/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		testCount := 10000000

		bar := progressbar.New(int(testCount))
		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := (rand.Int() % len(r.id2Item)) + 1
			item := r.id2Item[uint64(idx)]
			score := rand.Int() % 10000
			score = item.score + score
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans), r.getToplinkCount()/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		testCount := 10000000

		bar := progressbar.New(int(testCount))
		beg := time.Now()
		for i := 0; i < testCount; i++ {
			r.shrink(0)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			r.GetPercentRank(uint64(idx))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
	}

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			r.GetExactRank(uint64(idx))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
	}

}

func TestRank(t *testing.T) {
	fmt.Println("TestRank")

	/*{
		var r *Rank = NewRank()
		for i := 0; i < 100; i++ {
			idx := i + 1
			score := i + 1
			r.UpdateScore(uint64(idx), score)
		}

		assert.Equal(t, true, r.Check())

		for i := 200; i < 300; i++ {
			idx := i + 1
			score := i + 1
			r.UpdateScore(uint64(idx), score)
		}

		assert.Equal(t, true, r.Check())

		r.Show()

		fmt.Println(r.UpdateScore(uint64(150), 150))

		r.UpdateScore(uint64(150), 10)

		r.Show()

		assert.Equal(t, true, r.Check())

		fmt.Println(r.UpdateScore(uint64(150), 10))

		r.Reset()

	}

	{
		var r *Rank = NewRank()

		for i := 0; i < 10000; i++ {
			idx := i + 1
			score := rand.Int()
			r.UpdateScore(uint64(idx), score)
		}

		assert.Equal(t, true, r.Check())

		lastC := r.spans[len(r.spans)-1]
		lastItem := lastC.tail.links[0].pprev

		assert.Equal(t, len(r.id2Item), r.GetExactRank(uint64(lastItem.id)))

		firstC := r.spans[0]
		firstItem := firstC.head.links[0].pnext
		assert.Equal(t, 1, r.GetExactRank(uint64(firstItem.id)))

		assert.Equal(t, 100, r.GetPercentRank(uint64(firstItem.id)))

		assert.Equal(t, 0, r.GetPercentRank(uint64(lastItem.id)))

		assert.Equal(t, 100, len(r.spans))

		assert.Equal(t, true, r.Check())

	}*/

	{
		var r *Rank = NewRank()

		for i := 0; i < 100000; i++ {
			idx := (rand.Int() % 100000) + 1
			score := rand.Int() % 1000000
			fmt.Println("i", i, idx, score)
			r.UpdateScore(uint64(idx), score)
		}

		assert.Equal(t, true, r.Check())

		lastC := r.spans[len(r.spans)-1]
		lastItem := lastC.tail.links[0].pprev

		//assert.Equal(t, len(r.id2Item), r.GetExactRank(uint64(lastItem.id)))

		firstC := r.spans[0]
		firstItem := firstC.head.links[0].pnext
		assert.Equal(t, 1, r.GetExactRank(uint64(firstItem.id)))

		assert.Equal(t, 100, r.GetPercentRank(uint64(firstItem.id)))

		assert.Equal(t, 0, r.GetPercentRank(uint64(lastItem.id)))

		assert.Equal(t, true, r.Check())

	}

}

func BenchmarkRank(b *testing.B) {
	var r *Rank = NewRank()
	for i := 0; i < b.N; i++ {
		idx := (i % 1000000) + 1
		score := rand.Int()
		r.UpdateScore(uint64(idx), score)
	}
}
