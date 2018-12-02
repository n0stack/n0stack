package transaction

import "log"

func WrapRollbackError(err error) {
	if err != nil {
		log.Printf("[CRITICAL] Failed to rollback: err=\n%s", err.Error())
	}
}
