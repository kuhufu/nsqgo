package nsqgo

import (
	"errors"
	"github.com/nsqio/go-nsq"
	"log"
	"github.com/kuhufu/nsqgo/watcher"
	"time"
)

type Option func(p *Options)

type Producer struct {
	opts *Options

	reconnectChan chan *Node
	nodes         *NodeManager
	watcher       watcher.Watcher
	exitC         chan struct{}
}

func NewProducer(opts ...Option) (*Producer, error) {
	p := &Producer{
		opts:          NewOptions(),
		reconnectChan: make(chan *Node),
		nodes:         NewNodes(),
	}

	for _, opt := range opts {
		opt(p.opts)
	}

	log.Printf("working on %v mode", p.opts.mode)

	p.runWatcher()

	go p.reconnect()

	return p, nil
}

func (p *Producer) Publish(topic string, data []byte) error {
	var err error

	for retry := 3; retry > 0; retry-- {
		nodes := filter(p.nodes.GetAll())
		node := p.opts.selectStrategy(nodes)

		if node == nil {
			return errors.New("无可用node")
		}

		log.Printf("publish on node: %v", node.addr)
		err = node.p.Publish(topic, data)
		if err != nil {
			log.Printf("publish error: %v: %v", node.p, err)
			node.Mark(err)
			p.reconnectChan <- node
			continue
		}
		break
	}
	return err
}

func filter(nodes []*Node) []*Node {
	var res []*Node
	for _, node := range nodes {
		if node.err == nil {
			res = append(res, node)
		}
	}

	return res
}

func (p *Producer) reconnect() {
	for {
		select {
		case node := <-p.reconnectChan:
			go func() {
				for {
					select {
					case <-node.exitC:
						log.Printf("退出节点：%v 的重连", node.addr)
						return
					default:
						time.Sleep(p.opts.reconnectInterval)
						log.Printf("重连节点：%v", node.addr)
						err := node.p.Ping()
						if err == nil {
							node.Mark(nil)
							return
						}
					}
				}
			}()
		case <-p.exitC:
			log.Println("reconnect goroutine exit")
			return
		}
	}
}

func (p *Producer) runWatcher() {
	switch p.opts.mode {
	case "nsqd":
		p.watcher = watcher.NewNsqdWatcher(p.opts.addrs...)
	case "nsqlookupd":
		p.watcher = watcher.NewNsqlookupdWatcher(p.opts.addrs...)
	}

	go p.watcher.Watch()
	go func() {
		for {
			select {
			case events := <-p.watcher.Event():
				for _, event := range events {
					switch event.Action {
					case watcher.ActionAdd:
						producer, err := nsq.NewProducer(event.Addr, p.opts.nsqConfig)
						if err != nil {
							log.Printf("添加节点出错：%v", err)
							continue
						}
						p.nodes.Add(NewNode(event.Addr, producer))
						log.Printf("添加节点：%v", event.Addr)
					case watcher.ActionDel:
						node := p.nodes.Del(event.Addr)
						node.Close()
						log.Printf("删除节点：%v", event.Addr)
					}
				}
			case <-p.watcher.Exit():
				log.Println("watcher exit")
				return
			}
		}
	}()
}
