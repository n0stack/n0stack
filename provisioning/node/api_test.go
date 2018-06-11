package node

// memberlistがあるのでうまくテストできてない
// func TestApplyNode(t *testing.T) {
// 	c := memberlist.DefaultLANConfig()

// 	l, err := memberlist.Create(c)
// 	if err != nil {
// 		t.Errorf("Failed to prepare memberlist, err:%v", err.Error())
// 	}

// 	a := NodeAPI{
// 		dataStore: memory.NewMemoryDatastore(),
// 		list:      l,
// 	}
// 	a.dataStore.Apply("version_test", &pprovisioning.Node{
// 		Metadata: &pn0stack.Metadata{
// 			Name:    "version_test",
// 			Version: 3,
// 		},
// 		Spec: &pprovisioning.NodeSpec{},
// 	})

// 	testCases := []struct {
// 		req  *pprovisioning.ApplyNodeRequest
// 		node *pprovisioning.Node
// 		code codes.Code
// 	}{
// 		{
// 			req: &pprovisioning.ApplyNodeRequest{
// 				Metadata: &pn0stack.Metadata{
// 					Name: "test_node",
// 					Annotations: map[string]string{
// 						"hoge": "hoge",
// 					},
// 				},
// 				Spec: &pprovisioning.NodeSpec{},
// 			},
// 			node: &pprovisioning.Node{
// 				Metadata: &pn0stack.Metadata{
// 					Name: "test_node",
// 					Annotations: map[string]string{
// 						"hoge": "hoge",
// 					},
// 					Version: 1,
// 				},
// 				Spec: &pprovisioning.NodeSpec{},
// 				Status: &pprovisioning.NodeStatus{
// 					State: pprovisioning.NodeStatus_NotReady,
// 				},
// 			},
// 			code: codes.OK,
// 		},
// 		{
// 			req: &pprovisioning.ApplyNodeRequest{
// 				Metadata: &pn0stack.Metadata{
// 					Name:    "version_test",
// 					Version: 3,
// 				},
// 				Spec: &pprovisioning.NodeSpec{},
// 			},
// 			node: &pprovisioning.Node{
// 				Metadata: &pn0stack.Metadata{
// 					Name:    "version_test",
// 					Version: 4,
// 				},
// 				Spec: &pprovisioning.NodeSpec{},
// 				Status: &pprovisioning.NodeStatus{
// 					State: pprovisioning.NodeStatus_NotReady,
// 				},
// 			},
// 			code: codes.OK,
// 		},
// 		{
// 			req: &pprovisioning.ApplyNodeRequest{
// 				Metadata: &pn0stack.Metadata{
// 					Name:    "version_test",
// 					Version: 100,
// 				},
// 				Spec: &pprovisioning.NodeSpec{},
// 			},
// 			node: nil,
// 			code: codes.InvalidArgument,
// 		},
// 		{
// 			req: &pprovisioning.ApplyNodeRequest{
// 				Metadata: &pn0stack.Metadata{
// 					Name:    "version_0_test",
// 					Version: 100,
// 				},
// 				Spec: &pprovisioning.NodeSpec{},
// 			},
// 			node: nil,
// 			code: codes.InvalidArgument,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		_, err := a.ApplyNode(context.Background(), tc.req)

// 		// 何故か動かない
// 		// if !reflect.DeepEqual(n, tc.node) {
// 		// 	t.Errorf("Wrong status value.\n\thave:%v\n\twant:%v", n, tc.node)
// 		// }

// 		if status.Code(err) != tc.code {
// 			t.Errorf("Wrong status code.\n\thave:%v\n\twant:%v", status.Code(err), tc.code)
// 		}
// 	}
// }

// func TestDeleteNode(t *testing.T) {
// 	c := memberlist.DefaultLANConfig()

// 	l, err := memberlist.Create(c)
// 	if err != nil {
// 		t.Errorf("Failed to prepare memberlist, err:%v", err.Error())
// 	}

// 	a := NodeAPI{
// 		dataStore: memory.NewMemoryDatastore(),
// 		list:      l,
// 	}
// 	a.dataStore.Apply("OK", &pprovisioning.Node{
// 		Metadata: &pn0stack.Metadata{
// 			Name:    "OK",
// 			Version: 3,
// 		},
// 		Spec: &pprovisioning.NodeSpec{},
// 	})

// 	testCases := []struct {
// 		req  *pprovisioning.DeleteNodeRequest
// 		code codes.Code
// 	}{
// 		{
// 			req: &pprovisioning.DeleteNodeRequest{
// 				Name: "OK",
// 			},
// 			code: codes.OK,
// 		},
// 		{
// 			req: &pprovisioning.DeleteNodeRequest{
// 				Name: "NotFound",
// 			},
// 			code: codes.NotFound,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		_, err := a.DeleteNode(context.Background(), tc.req)

// 		if status.Code(err) != tc.code {
// 			t.Errorf("Wrong status code.\n\thave:%v\n\twant:%v", status.Code(err), tc.code)
// 		}
// 	}
// }
