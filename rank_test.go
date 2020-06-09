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
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func init() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

}

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
				score = item.value + rand.Int()%10000
			}
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg), len(r.id2Item))
		//fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
	}
}

func TestBenchmarkRank1(t *testing.T) {
	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 10000000
	idRange := 10000000

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int()%1000000 + 1
			//fmt.Println(idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		//fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{
		testCount := 10000000
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int()%1000000 + 1
			//fmt.Println(idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		//fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		testCount := 10000000

		bar := progressbar.New(int(testCount))
		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := (rand.Int() % len(r.id2Item)) + 1
			item := r.id2Item[uint64(idx)]
			score := rand.Int()%10000 + 1
			score = item.value + score
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		//fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	/*	{

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
	*/
}

func TestRank(t *testing.T) {
	fmt.Println("TestRank")

	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 200
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

}

func BenchmarkRank(b *testing.B) {
	var r *Rank = NewRank()
	for i := 0; i < b.N; i++ {
		idx := (i % 1000000) + 1
		score := rand.Int()
		r.UpdateScore(uint64(idx), score)
	}
}
