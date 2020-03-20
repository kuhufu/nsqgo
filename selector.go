package nsqgo

import (
	"math/rand"
	"sync"
)

type Selector func(nodes []*Node) *Node

func Random() Selector {
	return func(nodes []*Node) *Node {
		if len(nodes) == 0 {
			return nil
		}
		//随机选择
		idx := rand.Int() % len(nodes)

		return nodes[idx]
	}
}

func RoundRobin() Selector {
	i := 0
	mu := &sync.Mutex{}
	return func(nodes []*Node) *Node {
		if len(nodes) == 0 {
			return nil
		}

		mu.Lock()
		i++
		mu.Unlock()

		idx := i % len(nodes)
		return nodes[idx]
	}
}
