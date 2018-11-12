package transaction

import "testing"

func TestTransaction(t *testing.T) {
	tx := Begin()

	called := false
	tx.PushRollback("test", func() error {
		called = true
		return nil
	})

	if err := tx.Rollback(); err != nil {
		t.Errorf("Make err nil: err=%s", err.Error())
	}
	if !called {
		t.Errorf("Rollback is not called")
	}
}

func TestTransactionManyTimes(t *testing.T) {
	tx := Begin()

	called := 0
	for i := 0; i < 10; i++ {
		tx.PushRollback("test", func() error {
			called++
			return nil
		})
	}

	if err := tx.Rollback(); err != nil {
		t.Errorf("Make err nil: err=%s", err.Error())
	}
	if called != 10 {
		t.Errorf("Rollback was not called 10 times")
	}
}
