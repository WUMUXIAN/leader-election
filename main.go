package main

import (
	"sync"

	"github.com/WUMUXIAN/leader-election/common"
)

func main() {
	// create a node instance.
	node := common.Node{Role: common.NodeRoleUnknown}

	err := node.ConnectToZK([]string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"})
	if err != nil {
		panic(err)
	}
	defer node.CloseZK()

	err = node.ElectLeaderByZK()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
