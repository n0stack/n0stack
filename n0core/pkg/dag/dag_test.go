package dag

import (
	"testing"

	"google.golang.org/grpc"
)

// func TestHoge(t *testing.T) {
// 	a := struct {
// 		Hoge string
// 		Foo  string
// 	}{
// 		"hage",
// 		"bar",
// 	}

// 	va := reflect.ValueOf(a)
// 	buf, err := json.Marshal(va.Interface())
// 	if err != nil {
// 		t.Errorf("err: %+v", err.Error())
// 	}

// 	t.Errorf("buf: %s", string(buf))
// }

func TestIsDAG(t *testing.T) {
	cases := []struct {
		name    string
		actions map[string]*Task
		result  bool
	}{
		{
			"loop",
			map[string]*Task{
				"g1": &Task{
					DependOn: []string{
						"g2",
					},
				},
				"g2": &Task{
					DependOn: []string{
						"g3",
					},
				},
				"g3": &Task{
					DependOn: []string{
						"g1",
					},
				},
			},
			false,
		},
		{
			"liner",
			map[string]*Task{
				"g1": &Task{
					DependOn: []string{
						"g2",
					},
				},
				"g2": &Task{
					DependOn: []string{
						"g3",
					},
				},
				"g3": &Task{},
			},
			true,
		},
	}

	for _, c := range cases {
		if IsDAG(c.actions) != c.result {
			t.Errorf("[%s]", c.name)
		}
	}
}

func TestDoDAG(t *testing.T) {
	endpoint := "localhost:20180"
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect API: endpoint=%s, err=%s", endpoint, err.Error())
	}
	defer conn.Close()

	cases := []struct {
		name string
		task *Task
	}{
		{
			"loop",
			&Task{
				ResourceType: "Node",
				Action:       "GetNode",
				Args:         map[string]interface{}{"Name": "mock-node"},
				DependOn: []string{
					"g3",
				},
			},
		},
	}

	for _, c := range cases {
		c.task.Do(conn)
	}
}

// func TestDoDAG(t *testing.T) {
// 	endpoint := "localhost:20180"
// 	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatalf("Failed to connect API: endpoint=%s, err=%s", endpoint, err.Error())
// 	}
// 	defer conn.Close()

// 	cases := []struct {
// 		name    string
// 		actions map[string]*Task
// 		result  bool
// 	}{
// 		{
// 			"loop",
// 			map[string]*Task{
// 				"g1": &Task{
// 					DependOn: []string{
// 						"g2",
// 					},
// 				},
// 				"g2": &Task{
// 					DependOn: []string{
// 						"g3",
// 					},
// 				},
// 				"g3": &Task{
// 					DependOn: []string{
// 						"g1",
// 					},
// 				},
// 			},
// 			false,
// 		},
// 		{
// 			"liner",
// 			map[string]*Task{
// 				"g1": &Task{
// 					DependOn: []string{
// 						"g2",
// 					},
// 				},
// 				"g2": &Task{
// 					DependOn: []string{
// 						"g3",
// 					},
// 				},
// 				"g3": &Task{},
// 			},
// 			true,
// 		},
// 	}

// 	for _, c := range cases {
// 		if DoDAG(c.actions) != c.result {
// 			t.Errorf("[%s]", c.name)
// 		}
// 	}
// }
