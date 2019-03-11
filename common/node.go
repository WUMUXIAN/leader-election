package common

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/WUMUXIAN/go-common-utils/slice"
	"github.com/samuel/go-zookeeper/zk"
)

// NodeRoleType defines the role of a node
type NodeRoleType int

// Enum the node role type
const (
	_ NodeRoleType = iota
	NodeRoleUnknown
	NodeRoleLeader
	NodeRoleFollower
)

var (
	zkElectionRoot = "/election"
)

// Node defines a node
type Node struct {
	Role         NodeRoleType
	zkConnection *zk.Conn
	zkPath       string
}

// ConnectToZK connects node to ZK server.
func (n *Node) ConnectToZK(zkServers []string) (err error) {
	n.zkConnection, _, err = zk.Connect([]string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}, time.Second)
	return
}

// CloseZK closes the zk connection
func (n *Node) CloseZK() {
	if n.zkConnection != nil {
		n.zkConnection.Close()
	}
}

func (n *Node) ElectLeaderByZK() error {
	if n.zkConnection == nil {
		return errors.New("zk not connected")
	}

	// try to create a node.
	_, err := n.zkConnection.Create(zkElectionRoot, []byte{}, 0, zk.WorldACL(zk.PermAll))
	if err != nil && err.Error() != "zk: node already exists" {
		return err
	}

	// try to create a sequential ephemeral znode.
	n.zkPath, err = n.zkConnection.Create(zkElectionRoot+"/n_", []byte{}, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	n.zkPath = strings.Replace(n.zkPath, zkElectionRoot+"/", "", 1)
	fmt.Println("Path Created:", n.zkPath)
	err = n.checkRole()
	return err
}

func (n *Node) PrintRole() {
	switch n.Role {
	case NodeRoleLeader:
		fmt.Println("<------ I'm the leader ------->")
	case NodeRoleFollower:
		fmt.Println("<------ I'm the follower ------->")
	default:
		fmt.Println("<------ My role is unknown ------->")
	}
}

func (n *Node) checkRole() error {
	// get all the chidren.
	children, _, err := n.zkConnection.Children(zkElectionRoot)
	if err != nil {
		return err
	}
	// sort the chidren
	sort.Strings(children)
	fmt.Println("Children in the path now:", children)
	// If I am the leader.
	if children[0] == n.zkPath {
		n.Role = NodeRoleLeader
	} else {
		// I'm not the leader.
		n.Role = NodeRoleFollower
		// In the case where the node is a follower, it needs to watch first smaller node
		index := slice.IndexOfString(children, n.zkPath)
		// Watch it.
		_, _, event, err := n.zkConnection.GetW(zkElectionRoot + "/" + children[index-1])
		if err != nil {
			return err
		}
		n.watch(event)
	}
	n.PrintRole()
	return nil
}

func (n *Node) watch(event <-chan zk.Event) {
	go func() {
		select {
		case eve := <-event:
			// The znode is deleted, means the watched node is gone.
			if eve.Type == zk.EventNodeDeleted {
				n.checkRole()
				break
			}
		}
	}()

}
