package dag

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	raceutil "n0st.ac/n0stack/n0core/pkg/util/race"
)

const ResultChanLength = 10

type Task struct {
	task     func(context.Context, io.Writer) error
	rollback func(context.Context, io.Writer) error

	DependsOn []string `yaml:"depends_on"`
	// IgnoreError bool     `yaml:"ignore_error"`
}

func (n Task) Run(ctx context.Context, out io.Writer) error {
	return n.task(ctx, out)
}

func (n Task) Rollback(ctx context.Context, out io.Writer) error {
	return n.rollback(ctx, out)
}

type DAG struct {
	Tasks map[string]*Task `yaml:"tasks"`
}

func (d *DAG) AddTask(name string, task, rollback func(context.Context, io.Writer) error, dependsOn []string) error {
	for _, t := range dependsOn {
		if _, ok := d.Tasks[t]; !ok {
			return fmt.Errorf("Depending task '%s' do not exist", t)
		}
	}

	if _, ok := d.Tasks[name]; !ok {
		return fmt.Errorf("The task name '%s' is duplicated", name)
	}

	d.Tasks[name] = &Task{
		task:      task,
		rollback:  rollback,
		DependsOn: dependsOn,
	}

	return nil
}

func (d DAG) getInverseIndex() (map[string][]string, error) {
	children := make(map[string][]string)
	for k := range d.Tasks {
		children[k] = make([]string, 0)
	}

	for k, v := range d.Tasks {
		for _, n := range v.DependsOn {
			if _, ok := d.Tasks[n]; !ok {
				return nil, fmt.Errorf("Depended task '%s' do not exist", n)
			}

			children[n] = append(children[n], k)
		}
	}

	return children, nil
}

// topological sort
// 実際遅いけどもういいや O(E^2 + V)
func (d DAG) Check() error {
	result := 0

	depending := make(map[string]int)
	for k := range d.Tasks {
		depending[k] = len(d.Tasks[k].DependsOn)
	}

	children, err := d.getInverseIndex()
	if err != nil {
		return err
	}

	s := make([]string, 0, len(d.Tasks))
	for k := range d.Tasks {
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

	if result != len(d.Tasks) {
		return fmt.Errorf("This request is not DAG")
	}

	return nil
}

type TaskErrors struct {
	RunErrors      map[string]error
	RollbackErrors map[string]error
}

func (d TaskErrors) Error() string {
	errorList := make([]string, 1, 1+len(d.RunErrors)+len(d.RollbackErrors))
	errorList[0] = "some tasks are failed"

	for k, v := range d.RunErrors {
		errorList = append(errorList, fmt.Sprintf("  %s: %s", k, v.Error()))
	}
	for k, v := range d.RollbackErrors {
		errorList = append(errorList, fmt.Sprintf("  rollback/%s: %s", k, v.Error()))
	}

	if len(errorList) != 1 {
		return strings.Join(errorList, "\n  ")
	}

	return ""
}

type ActionResult struct {
	Name string
	Err  error
}

// 出力で時間を出したほうがよさそう
func (d DAG) Do(ctx context.Context, out io.Writer) error {
	if err := d.Check(); err != nil {
		return err
	}

	children, err := d.getInverseIndex()
	if err != nil {
		return err
	}

	mout := raceutil.NewLockedWriter(out)

	done, runErrs := d.run(ctx, children, mout)

	rollbackErrs := map[string]error(nil)
	if runErrs != nil {
		rollbackErrs = d.rollback(done, children, mout)
	}

	if runErrs != nil || rollbackErrs != nil {
		return &TaskErrors{
			RunErrors:      runErrs,
			RollbackErrors: rollbackErrs,
		}
	}

	return nil
}

func (d DAG) run(ctx context.Context, children map[string][]string, out io.Writer) ([]string, map[string]error) {
	depending := make(map[string]int)
	for k := range d.Tasks {
		depending[k] = len(d.Tasks[k].DependsOn)
	}

	resultChan := make(chan ActionResult, ResultChanLength)
	wg := new(sync.WaitGroup)
	total := len(d.Tasks)
	done := make([]string, 0, total)
	errs := make(map[string]error)

	taskCtx := context.Background()
	runTask := func(taskName string) {
		wg.Add(1)
		defer wg.Done()

		err := d.Tasks[taskName].Run(taskCtx, out)
		resultChan <- ActionResult{
			Name: taskName,
			Err:  err,
		}
	}

	canceled := false
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-ctxWithCancel.Done()

		// すでにリクエストしたタスクの終了を待つ
		wg.Wait()
		close(resultChan)
	}()

	for k := range d.Tasks {
		if depending[k] == 0 {
			go runTask(k)
		}
	}

	for r := range resultChan {
		if r.Err == nil {
			done = append(done, r.Name)
		} else {
			errs[r.Name] = r.Err
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
				cancel()
			}
		}
	}

	if len(errs) != 0 {
		return done, errs
	}
	return done, nil
}

func (d DAG) rollback(done []string, children map[string][]string, out io.Writer) map[string]error {
	depending := make(map[string]int)
	for n := range d.Tasks {
		depending[n] = len(children[n])
	}
	for n := range d.Tasks {
		f := false
		for _, k := range done {
			f = f || n == k
		}
		if f {
			continue
		}

		for _, k := range d.Tasks[n].DependsOn {
			depending[k]--
		}
	}

	resultChan := make(chan ActionResult, ResultChanLength)
	wg := new(sync.WaitGroup)
	total := len(done)
	rollbacked := 0
	errs := make(map[string]error)

	taskCtx := context.Background()
	rollbackTask := func(taskName string) {
		wg.Add(1)
		defer wg.Done()

		err := d.Tasks[taskName].Rollback(taskCtx, out)
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
			errs[r.Name] = r.Err
		}

		// Successful completion
		if rollbacked == total {
			close(resultChan)
		}

		// queueing
		for _, n := range d.Tasks[r.Name].DependsOn {
			depending[n]--
			if depending[n] == 0 {
				go rollbackTask(n)
			}
		}
	}

	if len(errs) != 0 {
		return errs
	}
	return nil
}
