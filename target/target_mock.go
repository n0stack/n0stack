package target

import "github.com/n0stack/n0core/model"

type MockTarget struct {
	PatchApply        func(m model.AbstractModel) (string, bool)
	PatchManagingType func() string
	// PatchInitialize   func() (string, bool)
	// Patchtest         func() (string, bool)
}

func (mock MockTarget) Apply(m model.AbstractModel) (string, bool) {
	return mock.PatchApply(m)
}

// func (mock MockTarget) Initialize() (string, bool) {
// 	return mock.PatchInitialize()
// }

// func (mock MockTarget) test() (string, bool) {
// 	return mock.Patchtest()
// }

func (mock MockTarget) ManagingType() string {
	return mock.PatchManagingType()
}
