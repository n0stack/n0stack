package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/hashicorp/memberlist"

	"github.com/n0stack/proto.go/v0"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

func GetIpmiAddress() string {
	out, err := exec.Command("ipmitool", "lan", "print").Output()
	if err != nil {
		return ""
	}

	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "IP Address              :") { // これで正しいのかよくわからず、要テスト
			s := strings.Split(l, " ")
			return s[len(s)-1]
		}
	}

	return ""
}

func GetSerial() string {
	out, err := exec.Command("dmidecode", "-t", "system").Output()
	if err != nil {
		return ""
	}

	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "Serial Number:") {
			s := strings.Split(l, " ")
			return s[len(s)-1]
		}
	}

	return ""
}

func registerNodeToAPI(name, advertiseAddress, api string) error {
	conn, err := grpc.Dial(api, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := pprovisioning.NewNodeServiceClient(conn)

	n, err := cli.GetNode(context.Background(), &pprovisioning.GetNodeRequest{Name: name})
	var ar *pprovisioning.ApplyNodeRequest
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return err
		}
		ar = &pprovisioning.ApplyNodeRequest{
			Metadata: &pn0stack.Metadata{
				Name: name,
			},
			Spec: &pprovisioning.NodeSpec{},
		}
	} else {
		ar = &pprovisioning.ApplyNodeRequest{
			Metadata: n.Metadata,
			Spec:     n.Spec,
		}
	}

	ar.Spec.Address = advertiseAddress
	// ar.Spec.Endpoints =
	ar.Spec.IpmiAddress = GetIpmiAddress()
	ar.Spec.Serial = GetSerial()

	n, err = cli.ApplyNode(context.Background(), ar)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Applied Node to APi on registerNodeToAPI, Node:%v", n)

	return nil
}

func joinNodeToMemberlist(name, advertiseAddress, api string) error {
	c := memberlist.DefaultLANConfig()
	c.Name = name
	c.AdvertiseAddr = advertiseAddress
	// c.AdvertisePort = int(a.Connection.Port)

	list, err := memberlist.Create(c)
	if err != nil {
		return err
	}

	_, err = list.Join([]string{api})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Join Node to memberlist on joinNodeToMemberlist, LocalNode:%v", list.LocalNode())

	return nil
}

func JoinNode(name, advertiseAddress, apiAddress string, apiPort int) error {
	addr, err := net.ResolveIPAddr("ip", advertiseAddress)
	if err != nil {
		return err
	}

	if err := registerNodeToAPI(name, addr.String(), fmt.Sprintf("%s:%d", apiAddress, apiPort)); err != nil {
		return err
	}

	if err := joinNodeToMemberlist(name, addr.String(), apiAddress); err != nil {
		return err
	}

	return nil
}

func LeaveNode(name, api string) error {
	conn, err := grpc.Dial(api, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := pprovisioning.NewNodeServiceClient(conn)

	_, err = cli.DeleteNode(context.Background(), &pprovisioning.DeleteNodeRequest{
		Name: name,
	})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Node from API")

	return nil
}
