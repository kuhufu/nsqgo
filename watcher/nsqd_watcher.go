package watcher

type nsqdWatcher struct {
	addrs     []string
	eventChan chan []Event
	exitChan  chan struct{}
}

func NewNsqdWatcher(addrs ...string) *nsqdWatcher {
	return &nsqdWatcher{
		addrs:     addrs,
		eventChan: make(chan []Event),
		exitChan:  make(chan struct{}),
	}
}

func (n *nsqdWatcher) Watch() {
	var events []Event

	for _, addr := range n.addrs {
		events = append(events, Event{
			Action: ActionAdd,
			Addr:   addr,
		})
	}

	n.eventChan <- events
	n.Close()
}

func (n *nsqdWatcher) Event() <-chan []Event {
	return n.eventChan
}

func (n *nsqdWatcher) Exit() <-chan struct{} {
	return n.exitChan
}

func (n *nsqdWatcher) Close() {
	close(n.exitChan)
}
