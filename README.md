### Leader election implementations in Go.

#### Zookeeper based leader election.

ZooKeeper is a distributed, high-performance coordination service that are used in many distributed applications. It provides name service, distributed locking, configuration, synchronization off-the-shelf. It's designed in a very simple way and easy to use.

- It has a file-system alike 'znode' tree structure, every name (znode) in ZooKeeper is a sequence of path that linked by "/", e.g. /election/node_1, /locks/lock1 and etc.
- unlike file system, each node can have information stored with it, also it has ephemeral nodes (when the client that create this node is gone (checked via heartbeat), the node will be deleted automatically).
- It supports "watches", as a client can set a watch on a particular znode, which will be triggered when the znode is changed or deleted. When a watch is triggered the client receives a packet saying that the znode has changed.
- It supports sequential creation of nodes, that means two clients are allowed to create the same znode, in this mode, and the system will guarantee that the nodes get created will be in sequential order.

##### ZooKeeper Guarantees

- Sequential Consistency - Updates from a client will be applied in the order that they were sent.
- Atomicity - Updates either succeed or fail. No partial results.
- Single System Image - A client will see the same view of the service regardless of the server that it connects to.
- Reliability - Once an update has been applied, it will persist from that time forward until a client overwrites the update.
- Timeliness - The clients view of the system is guaranteed to be up-to-date within a certain time bound.


##### ZooKeeper based leader election
Leverage on the ZooKeeper features, we can design the following leader election algorithm:

1. When a node is on-line, it creates a znode in ZooKeeper namely /election/node_ in sequential and ephemeral node.
2. It then get all children in the znode /election.
3. It will check whether itself is the smallest /election/node_, if it is, then it's the leader, if it's not, then it's the follower.
4. If it's the follower, and let it be the znode /election/node_i, then it sets a watch on znode /election/node_j, where j < i and j is the largest amongst the rest.
5. When the node receives a notification its watched znode is deleted, repeat the process from 2.

The implementation is in the function `func (n *Node) ElectLeaderByZK() error`
