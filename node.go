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

	//拷贝，返回副本，不然会有并发问题（watcher协程会增删节点）
	//增加节点还好说，删除节点会把底层切片索引位置的节点与最后一个节点交换，并将最后一个节点位置置为nil
	return append(m.nodes[:0:0], m.nodes...)
}
