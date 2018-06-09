package node

import (
	"log"

	"github.com/hashicorp/memberlist"
	"github.com/n0stack/n0core/datastore"
	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
)

type NodeAPIEventDelegate struct {
	ds datastore.Datastore
}

func (a NodeAPIEventDelegate) NotifyJoin(n *memberlist.Node) {
	node := &pprovisioning.Node{}

	if err := a.ds.Get(n.Name, node); err != nil {
		return
	}
	if node.Metadata == nil {
		return
	}

	if node.Status == nil {
		node.Status = &pprovisioning.NodeStatus{}
	}
	node.Status.State = pprovisioning.NodeStatus_Ready

	if err := a.ds.Apply(node.Metadata.Name, node); err != nil {
		return
	}
	log.Printf("[INFO] On NotifyJoin, applied Node:%v", node)

	return
}

func (a NodeAPIEventDelegate) NotifyLeave(n *memberlist.Node) {
	node := &pprovisioning.Node{}

	if err := a.ds.Get(n.Name, node); err != nil {
		return
	}
	if node.Metadata == nil {
		return
	}

	if node.Status == nil {
		node.Status = &pprovisioning.NodeStatus{}
	}
	node.Status.State = pprovisioning.NodeStatus_NotReady

	if err := a.ds.Apply(node.Metadata.Name, node); err != nil {
		return
	}
	log.Printf("[INFO] On NotifyLeave, applied Node:%v", node)

	return
}

func (a NodeAPIEventDelegate) NotifyUpdate(n *memberlist.Node) {
	log.Print("NotifyUpdate 検証 #82")
}
