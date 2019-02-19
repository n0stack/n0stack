package transaction

import (
	"fmt"
	"log"

	"github.com/cenkalti/backoff"
)

type RollbackTask struct {
	Name string
	Func func() error
}

type Transaction struct {
	stack []*RollbackTask
}

func Begin() *Transaction {
	t := &Transaction{}
	t.stack = make([]*RollbackTask, 0)

	return &Transaction{}
}

func (tx *Transaction) PushRollback(name string, f func() error) {
	tx.stack = append(tx.stack, &RollbackTask{
		Name: name,
		Func: f,
	})
}

func (tx *Transaction) PopRollback() *RollbackTask {
	l := len(tx.stack)
	if l == 0 {
		return nil
	}

	ret := tx.stack[l-1]
	tx.stack = tx.stack[:l-1]

	return ret
}

func (tx *Transaction) Rollback() error {
	errMes := ""

	for r := tx.PopRollback(); r != nil; r = tx.PopRollback() {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
		err := backoff.Retry(r.Func, b)
		if err != nil {
			errMes = fmt.Sprintf("  [%s] %s\n%s", r.Name, err.Error(), errMes)
		}
	}

	if errMes != "" {
		return fmt.Errorf(errMes)
	}
	return nil
}

func (tx *Transaction) RollbackWithLog() {
	if err := tx.Rollback(); err != nil {
		log.Printf("[CRITICAL] Failed to rollback: err=\n%s", err.Error())
	}
}

func (tx *Transaction) Commit() {
	tx.stack = make([]*RollbackTask, 0)
}
