package dag

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/golang/protobuf/jsonpb"
)

var Marshaler = &jsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: false,
	OrigName:     true,
}

type Node struct {
	task     func(context.Context, io.Writer) error
	rollback func(context.Context, io.Writer) error

	DependsOn []string `yaml:"depends_on"`
	// IgnoreError bool     `yaml:"ignore_error"`
}

func NewNode() *Node {
	return nil
}

func (n Node) Do(ctx context.Context, out io.Writer) error {
	return n.task(ctx, out)
}

func (n Node) Rollback(ctx context.Context, out io.Writer) error {
	return n.rollback(ctx, out)
}

type DAG struct {
	nodes map[string]*Node
}

func (d *DAG) getInverseIndex() (map[string][]string, error) {
	children := make(map[string][]string)
	for k := range d.nodes {
		children[k] = make([]string, 0)
	}

	for k, v := range d.nodes {
		for _, n := range v.DependsOn {
			if _, ok := d.nodes[n]; !ok {
				return nil, fmt.Errorf("Depended task '%s' do not exist", n)
			}

			children[n] = append(children[n], k)
		}
	}

	return children, nil
}

// topological sort
// 実際遅いけどもういいや O(E^2 + V)
// 副作用をなくしたい
func (d *DAG) Check() error {
	result := 0

	depending := make(map[string]int)
	for k := range d.nodes {
		depending[k] = len(d.nodes[k].DependsOn)
	}

	children, err := d.getInverseIndex()
	if err != nil {
		return err
	}

	s := make([]string, 0, len(d.nodes))
	for k := range d.nodes {
		if depending[k] == 0 {
			s = append(s, k)
			result++
		}
	}

	for len(s) != 0 {
		n := s[len(s)-1]
		s = s[:len(s)-1]

		for _, c := range children[n] {
			depending[c]--
			if depending[c] == 0 {
				s = append(s, c)
				result++
			}
		}
	}

	if result != len(d.nodes) {
		return fmt.Errorf("This request is not DAG")
	}

	return nil
}

type ActionResult struct {
	Name string
	Err  error
}

// 出力で時間を出したほうがよさそう
func (d *DAG) Do(ctx context.Context, out io.Writer) error {
	if err := d.Check(); err != nil {
		return err
	}

	children, err := d.getInverseIndex()
	if err != nil {
		return err
	}

	errorList := []string{
		"some tasks are failed",
	}

	depending := make(map[string]int)
	for k := range d.nodes {
		depending[k] = len(d.nodes[k].DependsOn)
	}

	resultChan := make(chan ActionResult, 100)
	wg := new(sync.WaitGroup)
	total := len(d.nodes)
	done := make([]string, 0, total)

	taskCtx := context.Background()
	runTask := func(taskName string) {
		wg.Add(1)
		defer wg.Done()

		err := d.nodes[taskName].Do(taskCtx, out)
		resultChan <- ActionResult{
			Name: taskName,
			Err:  err,
		}
	}

	canceled := false
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		<-ctxWithCancel.Done()

		// すでにリクエストしたタスクの終了を待つ
		wg.Wait()
		close(resultChan)
	}()

	for k := range d.nodes {
		if depending[k] == 0 {
			go runTask(k)
		}
	}

	for r := range resultChan {
		if r.Err == nil {
			done = append(done, r.Name)
		} else {
			errorList = append(errorList, fmt.Sprintf("%s: %s", r.Name, r.Err))
		}

		if !canceled {
			if r.Err != nil {
				canceled = true
				cancel()

			} else {
				// queueing
				for _, n := range children[r.Name] {
					depending[n]--
					if depending[n] == 0 {
						go runTask(n)
					}
				}
			}

			// Successful completion
			if len(done) == total {
				close(resultChan)
			}
		}
	}

	// TODO: rollback
	if canceled {
		depending := make(map[string]int)
		for n := range d.nodes {
			depending[n] = len(children[n])
		}
		for n := range d.nodes {
			for _, k := range done {
				if n == k {
					continue
				}
			}

			for _, k := range d.nodes[n].DependsOn {
				depending[k]--
			}
		}

		resultChan := make(chan ActionResult, 100)
		wg := new(sync.WaitGroup)
		total := len(done)
		rollbacked := 0

		taskCtx := context.Background()
		rollbackTask := func(taskName string) {
			wg.Add(1)
			defer wg.Done()

			err := d.nodes[taskName].Rollback(taskCtx, out)
			resultChan <- ActionResult{
				Name: taskName,
				Err:  err,
			}
		}

		for _, k := range done {
			if depending[k] == 0 {
				go rollbackTask(k)
			}
		}

		for r := range resultChan {
			rollbacked++

			if r.Err != nil {
				errorList = append(errorList, fmt.Sprintf("rollback/%s: %s", r.Name, r.Err))
			}

			// Successful completion
			if rollbacked == total {
				close(resultChan)
			}

			// queueing
			for _, n := range d.nodes[r.Name].DependsOn {
				depending[n]--
				if depending[n] == 0 {
					go rollbackTask(n)
				}
			}
		}
	}

	if len(errorList) != 1 {
		return errors.New(strings.Join(errorList, "\n  "))
	}

	return nil
}
