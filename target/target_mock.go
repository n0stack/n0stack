package target

import "github.com/n0stack/n0core/model"

type MockTarget struct {
	PatchManagingType func() string
	PatchOperations   func(state, task string) ([]func(n model.AbstractModel) (string, bool, string), bool)
	// PatchInitialize   func() (string, bool)
	// Patchtest         func() (string, bool)
}

func (mock MockTarget) Operations(state, task string) ([]func(n model.AbstractModel) (string, bool, string), bool) {
	return mock.PatchOperations(state, task)
}

func (mock MockTarget) ManagingType() string {
	return mock.PatchManagingType()
}

// func (mock MockTarget) Initialize() (string, bool) {
// 	return mock.PatchInitialize()
// }

// func (mock MockTarget) test() (string, bool) {
// 	return mock.Patchtest()
// }
