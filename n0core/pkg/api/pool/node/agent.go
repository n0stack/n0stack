package node

import (
	"context"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"code.cloudfoundry.org/bytefmt"
	"github.com/cenkalti/backoff"
	"github.com/n0stack/n0stack/n0proto/pool/v0"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

	mem, err := GetTotalMemory()
	if err != nil {
		return err
	}

	aip, err := net.ResolveIPAddr("ip", advertiseAddress)
	if err != nil {
		return err
	}

	ar := &ppool.ApplyNodeRequest{
		Name:          name,
		Address:       aip.String(),
		IpmiAddress:   GetIpmiAddress(),
		Serial:        GetSerial(),
		CpuMilliCores: GetTotalCPUMilliCores() * 1000,
		MemoryBytes:   mem,
		StorageBytes:  100. * bytefmt.GIGABYTE,
	}

	n, err := cli.GetNode(context.Background(), &ppool.GetNodeRequest{Name: name})
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return err
		}
	} else {
		log.Printf("[INFO] Get res: '%+v'", n)
		ar.Version = n.Version
		ar.Annotations = n.Annotations
	}

	log.Printf("[INFO] Apply req: '%+v'", ar)

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	err = backoff.Retry(func() error {
		n, err = cli.ApplyNode(context.Background(), ar)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Applied Node to APi on registerNodeToAPI, Node:%v", n)

		return nil
	}, b)
	if err != nil {
		return err
	}

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
