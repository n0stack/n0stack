package dag

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/n0stack/n0stack/n0proto/deployment/v0"
	"github.com/n0stack/n0stack/n0proto/pool/v0"
	"github.com/n0stack/n0stack/n0proto/provisioning/v0"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

var Marshaler = &jsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: false,
	OrigName:     true,
}

type Task struct {
	Type        string      `yaml:"type"`
	Action      string      `yaml:"action"`
	Args        interface{} `yaml:"args"`
	DependsOn   []string    `yaml:"depends_on"`
	IgnoreError bool        `yaml:"ignore_error"`
	// Rollback []*Task `yaml:"rollback"`

	child   []string
	depends int
}

// referenced by https://stackoverflow.com/questions/40737122/convert-yaml-to-json-without-struct
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}

	return i
}

// return response JSON bytes
func (a Task) Do(conn *grpc.ClientConn) (proto.Message, error) {
	var grpcCliType reflect.Type
	var grpcCliValue reflect.Value

	// TODO :生成自動化
	switch a.Type {
	case "node", "Node":
		grpcCliType = reflect.TypeOf(ppool.NewNodeServiceClient(conn))
		grpcCliValue = reflect.ValueOf(ppool.NewNodeServiceClient(conn))
	case "network", "Network":
		grpcCliType = reflect.TypeOf(ppool.NewNetworkServiceClient(conn))
		grpcCliValue = reflect.ValueOf(ppool.NewNetworkServiceClient(conn))
	case "block_storage", "BlockStorage":
		grpcCliType = reflect.TypeOf(pprovisioning.NewBlockStorageServiceClient(conn))
		grpcCliValue = reflect.ValueOf(pprovisioning.NewBlockStorageServiceClient(conn))
	case "virtual_machine", "VirtualMachine":
		grpcCliType = reflect.TypeOf(pprovisioning.NewVirtualMachineServiceClient(conn))
		grpcCliValue = reflect.ValueOf(pprovisioning.NewVirtualMachineServiceClient(conn))
	case "image", "Image":
		grpcCliType = reflect.TypeOf(pdeployment.NewImageServiceClient(conn))
		grpcCliValue = reflect.ValueOf(pdeployment.NewImageServiceClient(conn))
	case "flavor", "Flavor":
		grpcCliType = reflect.TypeOf(pdeployment.NewFlavorServiceClient(conn))
		grpcCliValue = reflect.ValueOf(pdeployment.NewFlavorServiceClient(conn))
	default:
		return nil, fmt.Errorf("Resource type '%s' do not exist", a.Type)
	}

	fnt, ok := grpcCliType.MethodByName(a.Action)
	if !ok {
		return nil, fmt.Errorf("Resource type '%s' do not have action '%s'", a.Type, a.Action)
	}

	// 1st arg is instance, 2nd is context.Background()
	// TODO: 何かがおかしい、 argsElem is "**SomeMessage", so use argsElem.Elem() in Call
	argsType := fnt.Type.In(2)
	argsElem := reflect.New(argsType)
	if a.Args == nil {
		a.Args = make(map[string]interface{})
	}
	buf, err := json.Marshal(convert(a.Args))
	if err != nil {
		return nil, fmt.Errorf("Args is invalid, set fields of message '%s' err=%s", argsType.String(), err.Error())
	}
	if err := json.Unmarshal(buf, argsElem.Interface()); err != nil {
		return nil, fmt.Errorf("Args is invalid, set fields of message '%s' err=%s", argsType.String(), err.Error())
	}

	out := fnt.Func.Call([]reflect.Value{grpcCliValue, reflect.ValueOf(context.Background()), argsElem.Elem()})
	if err, _ := out[1].Interface().(error); err != nil {
		return nil, fmt.Errorf("got error response: %s", err.Error())
	}

	return out[0].Interface().(proto.Message), nil
}

// topological sort
// 実際遅いけどもういいや O(E^2 + V)
func CheckDAG(tasks map[string]*Task) error {
	result := 0

	for k, _ := range tasks {
		tasks[k].child = make([]string, 0)
		tasks[k].depends = len(tasks[k].DependsOn)
	}

	for k, v := range tasks {
		for _, d := range v.DependsOn {
			if _, ok := tasks[d]; !ok {
				return fmt.Errorf("Depended task '%s' do not exist", d)
			}

			tasks[d].child = append(tasks[d].child, k)
		}
	}

	s := make([]string, 0, len(tasks))
	for k, v := range tasks {
		if v.depends == 0 {
			s = append(s, k)
			result++
		}
	}

	for len(s) != 0 {
		n := s[len(s)-1]
		s = s[:len(s)-1]

		for _, c := range tasks[n].child {
			tasks[c].depends--
			if tasks[c].depends == 0 {
				s = append(s, c)
				result++
			}
		}
	}

	if result != len(tasks) {
		return fmt.Errorf("This request is not DAG")
	}

	return nil
}

type ActionResult struct {
	Name string
	Res  proto.Message
	Err  error
}

// 出力で時間を出したほうがよさそう
func DoDAG(tasks map[string]*Task, out io.Writer, conn *grpc.ClientConn) bool {
	for k, _ := range tasks {
		tasks[k].child = make([]string, 0)
		tasks[k].depends = len(tasks[k].DependsOn)
	}

	for k, v := range tasks {
		for _, d := range v.DependsOn {
			tasks[d].child = append(tasks[d].child, k)
		}
	}

	resultChan := make(chan ActionResult, 100)
	wg := new(sync.WaitGroup)
	total := len(tasks)
	done := 0

	doTask := func(taskName string) {
		defer wg.Done()

		result, err := tasks[taskName].Do(conn)
		resultChan <- ActionResult{
			Name: taskName,
			Res:  result,
			Err:  err,
		}
	}

	for k, v := range tasks {
		if v.depends == 0 {
			wg.Add(1)
			fmt.Fprintf(out, "---> Task '%s' is started\n", k)
			log.Printf("[DEBUG] Task '%s' is started: %+v", k, v)
			go doTask(k)
		}
	}

	failed := false
	for r := range resultChan {
		done++

		if r.Err != nil {
			fmt.Fprintf(out, "---> [ %d/%d ] Task '%s' is failed: %s\n", done, total, r.Name, r.Err.Error())

			if !tasks[r.Name].IgnoreError && !failed {
				failed = true

				// すでにリクエストしたタスクの終了を待つ
				fmt.Fprintf(out, "---> Wait to finish requested tasks\n")
				go func() {
					wg.Wait()
					close(resultChan)
				}()
			}
		} else {
			res, _ := Marshaler.MarshalToString(r.Res)

			if failed {
				fmt.Fprintf(out, "---> [ %d/%d ] Task '%s', which was requested until failed, is finished\n--- Response ---\n%s\n", done, total, r.Name, res)
			} else {
				fmt.Fprintf(out, "---> [ %d/%d ] Task '%s' is finished\n--- Response ---\n%s\n", done, total, r.Name, res)

				// queueing
				for _, d := range tasks[r.Name].child {
					tasks[d].depends--
					if tasks[d].depends == 0 {
						wg.Add(1)
						fmt.Fprintf(out, "---> Task '%s' is started\n", d)
						log.Printf("[DEBUG] Task '%s' is started: %+v", d, tasks[d])
						go doTask(d)
					}
				}
			}
		}

		if !failed && done == total {
			close(resultChan)
		}
	}

	if failed {
		// TODO: rollback

		return false
	}

	return true
}
