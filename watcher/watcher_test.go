package watcher

import (
	"fmt"
	"testing"
)

func TestFetchNodeInfos(t *testing.T) {
	addrs, err := fetchNodes("http://localhost:4161/nodes")
	if err != nil {
		t.Error(err)
	}

	t.Log(addrs)
}

func TestSub(t *testing.T) {
	left := []string{"1", "2", "3"}
	right := []string{"1", "2", "3", "4"}
	fmt.Println(sub(left, right))
}
