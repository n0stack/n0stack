package node

import (
	"context"
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/n0stack/n0core/notification"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/memberlist"
	pnode "github.com/n0stack/proto.go/node/v0"
	"google.golang.org/grpc"
)

type Agent struct {
	id   *uuid.UUID
	list *memberlist.Memberlist
}

// TODO: notification全般

func GetAgent(addr string, port int, starter string) (*Agent, error) {
	a := &Agent{}
	var err error

	a.id, err = GetNodeUUID()
	if err != nil {
		panic(err)
	}

	c := memberlist.DefaultLANConfig()
	c.Name = a.id.String()
	c.AdvertisePort = port

	if c.AdvertiseAddr == "" || c.AdvertiseAddr == "0.0.0.0" {
		ip, err := sockaddr.GetPrivateIP()
		if err != nil {
			return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
		}
		if ip == "" {
			return nil, fmt.Errorf("No private IP address found, and explicit IP not provided")
		}

		c.AdvertiseAddr = ip
	} else {
		c.AdvertiseAddr = addr
	}

	a.list, err = memberlist.Create(c)
	if err != nil {
		// panic("Failed to create memberlist: " + err.Error())
		return nil, err
	}

	if starter != "" {
		_, err := a.list.Join([]string{starter})
		if err != nil {
			// panic("Failed to join cluster: " + err.Error())
			return nil, err
		}
	}

	return a, nil
}

func (a *Agent) List(context.Context, *pnode.ListRequest) (res *pnode.ListResponse, errRes error) {
	res = &pnode.ListResponse{}

	for _, m := range a.list.Members() {
		dest := fmt.Sprintf("%s:%d", m.Addr.String(), m.Port)
		conn, err := grpc.Dial(dest, grpc.WithInsecure())
		if err != nil {
			return nil, nil
		}
		defer conn.Close()

		c := pnode.NewNodeServiceClient(conn)

		r, err := c.Get(context.Background(), &pnode.GetRequest{})
		if err != nil {
			notification.Notify(notification.MakeNotification("List.GetNode", false, "Could not get from "+m.Name+" "+dest))
		}

		res.Nodes = append(res.Nodes, r)
	}

	return
}

func (a *Agent) Get(context.Context, *pnode.GetRequest) (res *pnode.Node, errRes error) {
	res = &pnode.Node{Id: a.id.String()}

	return
}
