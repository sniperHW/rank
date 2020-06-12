package rank

//go test -covermode=count -v -coverprofile=coverage.out -run=. -cpuprofile=rank.p
//go tool cover -html=coverage.out
//go tool pprof rank.p
//go test -v -run=^$ -bench BenchmarkRank -count 10

import (
	//"fmt"
	//"github.com/schollz/progressbar"
	//"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	//"time"
)

func TestSkipLists(t *testing.T) {

	for c := 0; c < 10000; c++ {

		sl1 := newSkipLists(0)

		{

			nodes := []*node{}
			for i := 0; i < rand.Int()%100; i++ {
				nodes = append(nodes, &node{
					key:   rand.Int() % 10000,
					value: uint64(i + i),
				})
			}

			for _, v := range nodes {
				sl1.InsertNode(v)
			}

		}

		sl2 := newSkipLists(0)

		{

			nodes := []*node{}
			for i := 0; i < rand.Int()%100; i++ {
				nodes = append(nodes, &node{
					key:   rand.Int() % 10000,
					value: uint64(i + i),
				})
			}

			for _, v := range nodes {
				sl2.InsertNode(v)
			}
		}

		sl1.merge(sl2)

		if !sl1.checkLink() {
			panic("checkLink")
		}
	}

}
