package model

import (
	"github.com/satori/go.uuid"
)

type (
	// Model is mapped to express graph data structure with _Dependency.
	//
	// About meta:
	// 	There are a lot of way to collaborate other service.
	//
	// 	#### Example 1: n0gateway
	// 	以下のようなリクエストをユーザーが行うと `9c2476ab-dc1e-4904-b8a4-6d991fdc7770` のUUIDに関連付けられているLBサービスにportが参加する。
	//
	// 	```yaml
	// 	type: resuorce/nic
	// 	state: running
	// 	meta:
	// 	n0stack/n0gateway/join: 9c2476ab-dc1e-4904-b8a4-6d991fdc7770
	// 	```
	//
	// 	n0gatewayとしては `/api/spec` を監視していればサービスディスカバリを用意に実装することができる。
	//
	// Args:
	//
	// 	id: UUID  default: generate uuid
	// 	type:
	// 	state:
	// 	meta:
	// 	dependencies: List of dependency to
	//
	// Example:
	// 	>>> new_vm = Model("resource/vm/kvm", "running")
	// 	>>> new_disk = Model("resource/volume/local", "claimed")
	// 	>>> new_vm.add_dependency(new_disk,
	// 							  "n0stack/n0core/resource/vm/attachments",
	// 							  {"n0stack/n0core/resource/vm/boot_priority": "1"})
	//
	// TODO:
	// 	- dependencyの2重定義ができないようにしたい
	Model struct {
		ID           uuid.UUID
		Type         string
		State        string // enumにしたい
		Name         string
		Meta         map[string]string
		Dependencies Dependencies
	}

	// AbstractModel is abstract interface for all models to use as Model.
	AbstractModel interface {
		GetModel() *Model
	}

	// Args:
	// 	model: Model which is depended.
	// 	label: A word about the relationship of dependence.
	// 	property: Additional options to explain the relationship of dependence.
	//
	// TODO:
	// 	- labelを書き込み可能にするか否か
	Dependency struct {
		Model    Model
		Label    string
		Property map[string]string
	}

	Dependencies []Dependency
)

func (d Dependencies) SelectWithLabel(l string) (Dependencies, error) {
	return d, nil
}

func (d Dependencies) SelectWithModelType(t string) (Dependencies, error) {
	return d, nil
}
