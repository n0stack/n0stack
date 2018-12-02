package virtualmachine

import "github.com/n0stack/n0stack/n0core/pkg/api/pool/node"

func init() {
	go UpMockAgent(node.MockNodeIP)
}
