package node

import (
	"context"
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

func JoinNode(name, advertiseAddress, api string) error {
	// register to API
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
			Metadata: &pn0stack.Metadata{},
			Spec:     &pprovisioning.NodeSpec{},
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

	// join to memberlist
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

	return nil
}
