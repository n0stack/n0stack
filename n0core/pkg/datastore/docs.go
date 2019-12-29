// Example:
//   var ds Datastore
//   v, err := ds.Get(key)
//   v, err = ds.Apply(key, value, v)
//   // ds.Apply(key, value, v-1) will be failed
//   err = ds.Delete(key, v)
package datastore
