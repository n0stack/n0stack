/*
Example

	func SomeEndpoint(ctx context.COntext, req SomeRequest) error {
		// validation

		tx := transaction.Begin()
		defer tx.Rollback()

		// API Process
		if err := Process(); err != nil {
			return err
		}
		tx.PushRollback("Process", func() error {
			return InverseProcess()
		})

		tx.Done()

		return nil
	}
*/
package transaction
