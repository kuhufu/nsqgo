package nsqgo

import (
	"github.com/nsqio/go-nsq"
	"sync"
)

type Node struct {
	addr  string
	p     *nsq.Producer
	err   error
	exitC chan struct{}
}

func NewNode(addr string, producer *nsq.Producer) *Node {
	return &Node{
		addr:  addr,
		p:     producer,
		exitC: make(chan struct{}),
	}
}

func (n *Node) Mark(err error) {
	n.err = err
}

func (n *Node) Close() {
	close(n.exitC)
	go n.p.Stop()
}

type NodeManager struct {
	nodes []*Node
	mu    *sync.RWMutex
}

func NewNodes() *NodeManager {
	return &NodeManager{
		mu: &sync.RWMutex{},
	}
}

func (m *NodeManager) Add(node *Node) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nodes = append(m.nodes, node)
}

func (m *NodeManager) Del(addr string) *Node {
	m.mu.Lock()
	defer m.mu.Unlock()

	nodes := m.nodes
	var target *Node

	for i, node := range nodes {
		if node.addr == addr {
			target = node
			nodes[i] = nodes[len(nodes)-1]
			nodes[len(nodes)-1] = nil
			m.nodes = nodes[:len(nodes)-1]
			break
		}
	}

	return target
}

func (m *NodeManager) GetAll() []*Node {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.nodes
}
