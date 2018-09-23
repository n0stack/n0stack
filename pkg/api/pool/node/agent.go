package node

import (
	"context"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/v0"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IPMIを持っていない場合が考えられるので、とりあえずエラーハンドリングはしていない
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

// Serialが取得できなくても動作に問題はないため、エラーハンドリングはしていない
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

func GetTotalMemory() (uint64, error) {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to call syscall 'sysinfo'")
	}

	return si.Totalram / uint64(si.Unit), nil
}

func GetTotalCPUMilliCores() uint32 {
	return uint32(runtime.NumCPU())
}

// TODO: エラーハンドリング適当
func RegisterNodeToAPI(name, advertiseAddress, api string) error {
	conn, err := grpc.Dial(api, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := ppool.NewNodeServiceClient(conn)

	n, err := cli.GetNode(context.Background(), &ppool.GetNodeRequest{Name: name})
	var ar *ppool.ApplyNodeRequest
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return err
		}

		ar = &ppool.ApplyNodeRequest{
			Metadata: &pn0stack.Metadata{
				Name: name,
			},
			Spec: &ppool.NodeSpec{},
		}
	} else {
		ar = &ppool.ApplyNodeRequest{
			Metadata: n.Metadata,
			Spec:     n.Spec,
		}
	}

	mem, err := GetTotalMemory()
	if err != nil {
		return err
	}

	ar.Spec.Address = advertiseAddress
	ar.Spec.IpmiAddress = GetIpmiAddress()
	ar.Spec.Serial = GetSerial()
	ar.Spec.CpuMilliCores = GetTotalCPUMilliCores()
	ar.Spec.MemoryBytes = mem

	n, err = cli.ApplyNode(context.Background(), ar)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Applied Node to APi on registerNodeToAPI, Node:%v", n)

	return nil
}

// func joinNodeToMemberlist(name, advertiseAddress, api string) (*memberlist.Memberlist, error) {
// 	c := memberlist.DefaultLANConfig()
// 	c.Name = name
// 	c.AdvertiseAddr = advertiseAddress
// 	// c.AdvertisePort = int(a.Connection.Port)

// 	list, err := memberlist.Create(c)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = list.Join([]string{api})
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("[INFO] Join Node to memberlist on joinNodeToMemberlist, LocalNode:%v", list.LocalNode())

// 	return list, nil
// }

// func JoinNode(name, advertiseAddress, apiAddress string, apiPort int) error {
// 	addr, err := net.ResolveIPAddr("ip", advertiseAddress)
// 	if err != nil {
// 		return err
// 	}

// 	if err := registerNodeToAPI(name, addr.String(), fmt.Sprintf("%s:%d", apiAddress, apiPort)); err != nil {
// 		return err
// 	}

// 	list, err := joinNodeToMemberlist(name, addr.String(), apiAddress)
// 	if err != nil {
// 		return err
// 	}

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)
// 	go func() {
// 		for _ = range c {
// 			if err := LeaveNode(list); err != nil {
// 				log.Fatalf("Failed to LeaveNode, err:%s", err.Error())
// 			}
// 			os.Exit(0)
// 		}
// 	}()

// 	return nil
// }

// func LeaveNode(list *memberlist.Memberlist) error {
// 	list.Leave(3 * time.Second)

// 	return nil
// }
