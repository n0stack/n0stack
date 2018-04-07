package kvm

import (
	"context"
	"fmt"

	"github.com/n0stack/n0core/lib"
	n0stack "github.com/n0stack/proto"
	"github.com/n0stack/proto/device/tap/v0"
	"github.com/n0stack/proto/device/vm"
	"github.com/n0stack/proto/device/volume"
	"google.golang.org/grpc"
)

type Agent struct {
	// DB *gorm.DB
}

const (
	modelType = "device/vm/kvm"
)

func notify(n *n0stack.Notification) {
	if !n.Success {
		panic(n)
	}

	fmt.Printf("%v\n", n)
}

// Apply スペックを元にステートレスに適用する
func (a *Agent) Apply(ctx context.Context, spec *vm.Spec) (n *n0stack.Notification, errRes error) {
	// var (
	// 	n *n0stack.Notification
	// 	k kvm
	// )
	// r = &vm.Response{
	// 	Notification: n,
	// 	Spec:         &k.Spec,
	// 	Status:       &k.Status,
	// }

	k, n := getVM(spec.Device.Model)
	if !n.Success {
		return n, nil
	}

	// if vm is not running
	if k.args == nil {
		// check CPU usage
		// check Memory usage

		notify(k.runVM(spec))
	}

	notify(k.connectQMP())
	k.qmp.Connect()
	defer k.qmp.Disconnect()

	// if n := k.getStatus(); !n.Success {
	// 	return n, nil
	// }

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		return nil, nil
	}
	defer conn.Close()

	for i, v := range spec.Volumes {
		cli := volume.NewStandardClient(conn)

		r, err := cli.Apply(context.Background(), v) // これはすでにあるvolumeを再利用したい時に危険である
		if err != nil {
			return nil, nil
		}

		notify(k.attachVolume(&volumeURL{
			id:  r.Spec.Model.Id,
			url: r.Status.Url,
		}, i+1))
	}

	for _, n := range spec.Nics {
		cli := tap.NewStandardClient(conn)

		r, err := cli.Apply(context.Background(), &tap.ApplyRequest{
			Model: n.Model,
			Spec:  n.Tap,
		})
		if err != nil {
			return nil, nil
		}

		notify(k.attachNIC(&nic{
			id:  r.Model.Id,
			tap: r.Status.Tap,
			mac: n.HwAddr.Address,
		}))
	}

	// Applyした時に毎回ブートしなければいけないわけではない
	// TODO: プロセスを始めたらbootする
	if n := k.bootVM(); !n.Success {
		return n, nil
	}

	return lib.MakeNotification("Apply", true, ""), nil
}

func (a *Agent) Delete(ctx context.Context, model *n0stack.Model) (*n0stack.Notification, error) {
	k, n := getVM(model)
	if !n.Success {
		return n, nil
	}

	// if vm is not running
	if k.args == nil {
		return lib.MakeNotification("Delete", true, "Process is not existing"), nil
	}

	// vmのシャットダウン

	n = k.kill()
	if !n.Success {
		return n, nil
	}

	// remove volumeはしない、なぜならデータが消えるとまずいから
	// remove networkはする

	// remove workdir

	return lib.MakeNotification("Delete", true, ""), nil
}
