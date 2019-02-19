/*
Example

	func SomeEndpoint(ctx context.COntext, req SomeRequest) error {
		// validation

		datastore.Lock(req.Key)
		defer datastore.Unlock(req.Key)

		tx := transaction.Begin()
		defer tx.Rollback()

		// API Process

		tx.Done()

		return nil
	}
*/
package datastore
