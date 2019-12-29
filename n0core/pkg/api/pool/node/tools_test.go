package node

import (
	"testing"

	"n0st.ac/n0stack/n0proto.go/budget/v0"
	"n0st.ac/n0stack/n0proto.go/pool/v0"
)

func TestIsLockedForDeletion(t *testing.T) {
	cases := []struct {
		Name string
		Node *ppool.Node
		Res  bool
	}{
		{
			Name: "nothing",
			Node: &ppool.Node{},
			Res:  false,
		},
		{
			Name: "not lock by compute",
			Node: &ppool.Node{
				ReservedComputes: map[string]*pbudget.Compute{
					"not lock1": &pbudget.Compute{
						Annotations: map[string]string{
							AnnotationComputeDisableDeletionLock: "true",
						},
					},
					"not lock2": &pbudget.Compute{
						Annotations: map[string]string{
							AnnotationComputeDisableDeletionLock: "true",
						},
					},
				},
			},
			Res: false,
		},
		{
			Name: "lock by compute",
			Node: &ppool.Node{
				ReservedComputes: map[string]*pbudget.Compute{
					"lock": &pbudget.Compute{},
				},
			},
			Res: true,
		},
		{
			Name: "lock after not lock by compute",
			Node: &ppool.Node{
				ReservedComputes: map[string]*pbudget.Compute{
					"not lock": &pbudget.Compute{
						Annotations: map[string]string{
							AnnotationComputeDisableDeletionLock: "true",
						},
					},
					"lock": &pbudget.Compute{},
				},
			},
			Res: true,
		},
		{
			Name: "not lock by storage",
			Node: &ppool.Node{
				ReservedStorages: map[string]*pbudget.Storage{
					"not lock1": &pbudget.Storage{
						Annotations: map[string]string{
							AnnotationStorageDisableDeletionLock: "true",
						},
					},
					"not lock2": &pbudget.Storage{
						Annotations: map[string]string{
							AnnotationStorageDisableDeletionLock: "true",
						},
					},
				},
			},
			Res: false,
		},
		{
			Name: "lock by storage",
			Node: &ppool.Node{
				ReservedStorages: map[string]*pbudget.Storage{
					"lock": &pbudget.Storage{},
				},
			},
			Res: true,
		},
		{
			Name: "lock after not lock by storage",
			Node: &ppool.Node{
				ReservedStorages: map[string]*pbudget.Storage{
					"not lock": &pbudget.Storage{
						Annotations: map[string]string{
							AnnotationStorageDisableDeletionLock: "true",
						},
					},
					"lock": &pbudget.Storage{},
				},
			},
			Res: true,
		},
	}

	for _, c := range cases {
		if IsLockedForDeletion(c.Node) != c.Res {
			t.Errorf("[%s] IsLockedForDeletion return value was wrong: Node=%v", c.Name, c.Node)
		}
	}
}
