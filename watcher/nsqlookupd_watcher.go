package watcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type nsqlookupdWatcher struct {
	addrs     []string
	ticker    *time.Ticker
	addr      string
	eventChan chan []Event
	exitChan  chan struct{}
}

const (
	ActionAdd = "ADD"
	ActionDel = "DEL"
)

type Action string

type Event struct {
	Action Action
	Addr   string
}

func NewNsqlookupdWatcher(addr ...string) *nsqlookupdWatcher {
	return &nsqlookupdWatcher{
		addr:      addr[0],
		ticker:    time.NewTicker(time.Second * 5),
		eventChan: make(chan []Event),
	}
}

func (w *nsqlookupdWatcher) Watch() {
	httpAddr := fmt.Sprintf("http://%v/nodes", w.addr)
	encounterErr := false
	for {
		addrs, err := fetchNodes(httpAddr)
		if err != nil {
			encounterErr = true
			log.Println("获取节点信息出错: ", err)
		}

		if err == nil {
			if encounterErr {
				encounterErr = false
				log.Println("上一次获取节点信息发生错误，睡眠30s，然后重新获取")
				//因为重启nsqlookupd之后，nsqlookupd需要过一段时间才能从nsqd得到节点信息，
				//如果直接获取就会导致节点为空（不过这样做有意义吗，毕竟重启nsqlookupd是很少见的情况）
				time.Sleep(time.Second * 30)
				continue
			}

			events := w.translateToEvent(addrs)

			w.addrs = addrs

			if len(events) != 0 {
				w.eventChan <- events
			}
		}

		select {
		case <-w.ticker.C:

		}
	}
}

func (w *nsqlookupdWatcher) Event() <-chan []Event {
	return w.eventChan
}

func (n *nsqlookupdWatcher) Exit() <-chan struct{} {
	return n.exitChan
}

func (n *nsqlookupdWatcher) Close() {
	close(n.exitChan)
}

func (w *nsqlookupdWatcher) translateToEvent(addrs []string) []Event {
	newAddrs := sub(addrs, w.addrs)
	delAddrs := sub(w.addrs, addrs)
	res := make([]Event, 0, len(newAddrs)+len(delAddrs))

	for _, addr := range delAddrs {
		res = append(res, Event{
			Action: ActionDel,
			Addr:   addr,
		})
	}

	for _, addr := range newAddrs {
		res = append(res, Event{
			Action: ActionAdd,
			Addr:   addr,
		})
	}

	return res
}

type nsqProducer struct {
	RemoteAddress string `json:"remote_address"`
	TcpPort       int    `json:"tcp_port"`
	HttpPort      int    `json:"http_port"`
}

func fetchNodes(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Producers []*nsqProducer `json:"producers"`
	}
	err = json.Unmarshal(bytes, &data)

	infos := data.Producers
	addrs := make([]string, 0, len(infos))
	for _, info := range infos {
		host, _, _ := net.SplitHostPort(info.RemoteAddress)
		addrs = append(addrs, net.JoinHostPort(host, strconv.Itoa(info.TcpPort)))
	}

	return addrs, err
}
