package watcher

type Watcher interface {
	Watch()
	Event() <-chan []Event
	Exit() <-chan struct{}
	Close()
}
