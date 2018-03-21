package kvm

import (
	"context"

	"github.com/n0stack/n0core/lib"
	n0stack "github.com/n0stack/proto"
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

// Apply スペックを元にステートレスに適用する
func (a *Agent) Apply(ctx context.Context, spec *vm.Spec) (*n0stack.Notification, error) {
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

		if n := k.runVM(spec); !n.Success {
			return n, nil
		}
	}

	if n := k.connectQMP(); !n.Success {
		return n, nil
	}
	k.qmp.Connect()
	defer k.qmp.Disconnect()

	// if n := k.listVolumes(); !n.Success {
	// 	return n, nil
	// }

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		return nil, nil
	}
	defer conn.Close()

	var volumes []volumeURL
	for _, v := range spec.Volume {
		cli := volume.NewStandardClient(conn)

		r, err := cli.Apply(context.Background(), v) // これはすでにあるvolumeを再利用したい時に危険である
		if err != nil {
			return nil, nil
		}

		volumes = append(volumes, volumeURL{
			id:  r.Spec.Model.Id,
			url: r.Status.Url,
		})
	}

	if n := k.attachVolumes(volumes); !n.Success {
		return n, nil
	}

	// k.attachNIC()

	// Applyした時に毎回ブートしなければいけないわけではない
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

	n = k.kill()
	if !n.Success {
		return n, nil
	}

	// remove volumeはしない、なぜならデータが消えるとまずいから
	// remove networkはする

	// remove workdir

	return lib.MakeNotification("Delete", true, ""), nil
}
