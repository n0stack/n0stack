package node

import (
	"log"

	"github.com/hashicorp/memberlist"
	"github.com/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/proto.go/pool/v0"
)

type NodeAPIEventDelegate struct {
	ds datastore.Datastore
}

func (a NodeAPIEventDelegate) NotifyJoin(n *memberlist.Node) {
	node := &ppool.Node{}
	if err := a.ds.Get(n.Name, node); err != nil {
		// log.Printf()
		return
	}
	if node == nil {
		return
	}

	if node.Status == nil {
		node.Status = &ppool.NodeStatus{}
	}
	node.Status.State = ppool.NodeStatus_Ready

	if err := a.ds.Apply(node.Metadata.Name, node); err != nil {
		// log.Printf()
		return
	}
	log.Printf("[INFO] On NotifyJoin, applied Node:%v", node)

	return
}

func (a NodeAPIEventDelegate) NotifyLeave(n *memberlist.Node) {
	node := &ppool.Node{}
	if err := a.ds.Get(n.Name, node); err != nil {
		// log.Printf()
		return
	}
	if node == nil {
		return
	}

	if node.Status == nil {
		node.Status = &ppool.NodeStatus{}
	}
	node.Status.State = ppool.NodeStatus_NotReady

	if err := a.ds.Apply(node.Metadata.Name, node); err != nil {
		// log.Printf()
		return
	}
	log.Printf("[INFO] On NotifyLeave, applied Node:%v", node)

	return
}

func (a NodeAPIEventDelegate) NotifyUpdate(n *memberlist.Node) {
	log.Print("NotifyUpdate 検証 #82")
}
