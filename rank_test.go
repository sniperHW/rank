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
	//"strings"
	"testing"
	"time"
)

func TestBenchmarkRank2(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 50000000
	idRange := 10000000
	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			item := r.getRankItem(uint64(idx))
			var score int
			if nil == item {
				score = rand.Int() % 1000000
			} else {
				score = item.score + rand.Int()%10000
			}
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg), len(r.id2Item))
	}
}

func TestBenchmarkRank(t *testing.T) {
	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 10000000
	idRange := 10000000

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
		fmt.Println(len(r.id2Item), len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())

		low := []string{}

		total := 0

		for _, v := range r.spans {
			total += v.count
			if v.count < len(r.id2Item)/len(r.spans) {
				low = append(low, fmt.Sprintf("(%d,%d)", v.idx, v.count))
			}
		}

		assert.Equal(t, len(r.id2Item), total)

		fmt.Println("< ", len(low))
		//fmt.Println(strings.Join(low, ","))

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
		fmt.Println(len(r.id2Item), len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())

		low := []string{}

		total := 0

		for _, v := range r.spans {
			total += v.count
			if v.count < len(r.id2Item)/len(r.spans) {
				low = append(low, fmt.Sprintf("(%d,%d)", v.idx, v.count))
			}
		}

		assert.Equal(t, len(r.id2Item), total)

		fmt.Println("< ", len(low))

		//fmt.Println(strings.Join(low, ","))

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
		fmt.Println(len(r.id2Item), len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())

		low := []string{}

		total := 0

		for _, v := range r.spans {
			total += v.count
			if v.count < len(r.id2Item)/len(r.spans) {
				low = append(low, fmt.Sprintf("(%d,%d)", v.idx, v.count))
			}
		}

		assert.Equal(t, len(r.id2Item), total)
		fmt.Println("< ", len(low))

		//fmt.Println(strings.Join(low, ","))

	}

	{

		testCount := 10000000

		bar := progressbar.New(int(testCount))
		beg := time.Now()
		for i := 0; i < testCount; i++ {
			r.shrink(0, nil)
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

	{
		var r *Rank = NewRank()
		for i := 0; i < 100; i++ {
			idx := i + 1
			score := i + 1
			r.UpdateScore(uint64(idx), score)
		}

		for i := 200; i < 300; i++ {
			idx := i + 1
			score := i + 1
			r.UpdateScore(uint64(idx), score)
		}

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
		lastItem := lastC.tail.pprev

		assert.Equal(t, len(r.id2Item), r.GetExactRank(uint64(lastItem.id)))

		firstC := r.spans[0]
		firstItem := firstC.head.pnext
		assert.Equal(t, 1, r.GetExactRank(uint64(firstItem.id)))

		assert.Equal(t, 100, r.GetPercentRank(uint64(firstItem.id)))

		assert.Equal(t, 0, r.GetPercentRank(uint64(lastItem.id)))

		assert.Equal(t, 100, len(r.spans))

		assert.Equal(t, true, r.Check())

	}

	{
		var r *Rank = NewRank()

		for i := 0; i < 100000; i++ {
			idx := (rand.Int() % 100000) + 1
			score := rand.Int() % 1000000
			r.UpdateScore(uint64(idx), score)
		}

		assert.Equal(t, true, r.Check())

		lastC := r.spans[len(r.spans)-1]
		lastItem := lastC.tail.pprev

		//assert.Equal(t, len(r.id2Item), r.GetExactRank(uint64(lastItem.id)))

		firstC := r.spans[0]
		firstItem := firstC.head.pnext
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
