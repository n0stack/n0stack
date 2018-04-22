package node

import (
	"context"
	"fmt"

	"github.com/n0stack/proto.go/metadata/v0"
	"github.com/n0stack/proto.go/node/v0"

	"github.com/satori/go.uuid"

	"github.com/n0stack/n0core/notification"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/memberlist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Agent struct {
	id *uuid.UUID
	pmetadata.Metadata
	pnode.Status

	list *memberlist.Memberlist
}

// TODO: notification全般

func GetAgent(addr string, port int, starter string, name string, labels map[string]string) (*Agent, error) {
	a := &Agent{
		Metadata: pmetadata.Metadata{
			Name:   name,
			Labels: labels,
		},
		Status: pnode.Status{
			Connection: &pnode.Connection{},
		},
	}

	var err error
	a.id, err = GetNodeUUID()
	if err != nil {
		panic(err)
	}

	a.Connection.Port = uint32(port)

	if addr == "" || addr == "0.0.0.0" {
		ip, err := sockaddr.GetPrivateIP()
		if err != nil {
			return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
		}
		if ip == "" {
			return nil, fmt.Errorf("No private IP address found, and explicit IP not provided")
		}

		a.Connection.Address = ip
	} else {
		a.Connection.Address = addr
	}

	c := memberlist.DefaultLANConfig()
	c.Name = a.id.String()
	c.AdvertiseAddr = a.Connection.Address
	c.AdvertisePort = int(a.Connection.Port)

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

// すべてのnodeに対して `GetNode` を行うことでListを生成するため非常に遅い
func (a *Agent) ListNodes(context.Context, *pnode.ListNodesRequest) (res *pnode.ListNodesResponse, errRes error) {
	res = &pnode.ListNodesResponse{}

	for _, m := range a.list.Members() {
		dest := fmt.Sprintf("%s:%d", m.Addr.String(), m.Port)
		conn, err := grpc.Dial(dest, grpc.WithInsecure())
		if err != nil {
			return nil, nil
		}
		defer conn.Close()

		c := pnode.NewNodeServiceClient(conn)

		r, err := c.GetNode(context.Background(), &pnode.GetNodeRequest{})
		if err != nil {
			notification.Notify(notification.MakeNotification("List.GetNode", false, "Could not get from "+m.Name+" "+dest))
		}

		res.Nodes = append(res.Nodes, r) // これは遅いので `a.list.Members()` の長さを取得してから最初にmakeするべき
	}

	return
}

// 接続先をメモリ上のmemberlistを線形探索することで決定し、1回 `GetNode` をリクエストする
// ローカルの情報についてはTCPコネクションを生成しないので高速
func (a *Agent) GetNode(ctx context.Context, req *pnode.GetNodeRequest) (res *pnode.Node, errRes error) {
	if req.Id == "" {
		res = &pnode.Node{
			Id:       a.id.String(),
			Metadata: &a.Metadata,
			Status:   &a.Status,
		}
		return
	}

	var err error
	id, err := uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	if id.String() == a.id.String() {
		res = &pnode.Node{
			Id:       a.id.String(),
			Metadata: &a.Metadata,
			Status:   &a.Status,
		}
		return
	}

	c, err := a.GetConnection(context.Background(), &pnode.GetConnectionRequest{NodeId: id.String()})
	if err != nil {
		errRes = err
		return
	}

	dest := fmt.Sprintf("%s:%d", c.Address, c.Port)
	conn, err := grpc.Dial(dest, grpc.WithInsecure())
	if err != nil {
		notification.Notify(notification.MakeNotification("List.GetNode", false, "Could not connect to node "+id.String()+": "+err.Error()))
	}
	defer conn.Close()

	cli := pnode.NewNodeServiceClient(conn)

	res, err = cli.GetNode(context.Background(), &pnode.GetNodeRequest{Id: id.String()})
	if err != nil {
		notification.Notify(notification.MakeNotification("List.GetNode", false, "Could not get from "+id.String()+" "+dest+": "+err.Error()))
	}

	return
}

// メモリ上のmemberlistを線形探索するため高速
func (a *Agent) GetConnection(ctx context.Context, req *pnode.GetConnectionRequest) (res *pnode.Connection, errRes error) {
	if req.NodeId == "" {
		res = a.Connection
		return
	}

	var err error
	id, err := uuid.FromString(req.NodeId)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.NodeId)
		return
	}

	if id.String() == a.id.String() { // 冗長な気がする
		res = a.Connection
		return
	}

	for _, m := range a.list.Members() {
		if m.Name == id.String() {
			res.Address = m.Addr.String()
			res.Port = uint32(m.Port)
			return
		}
	}

	errRes = grpc.Errorf(codes.NotFound, "message:Node do not exist")
	return
}
