package main

import (
	"fmt"

	"github.com/hashicorp/memberlist"
)

func main() {
	/* Create the initial memberlist from a safe configuration.
	   Please reference the godoc for other default config types.
	   http://godoc.org/github.com/hashicorp/memberlist#Config
	*/
	c1 := memberlist.DefaultLocalConfig()
	c1.Name = "hogehoge"
	l1, err := memberlist.Create(c1)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	c2 := memberlist.DefaultLocalConfig()
	c2.Name = "hogehoge"
	l2, err := memberlist.Create(c2)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	// Join an existing cluster by specifying at least one known member.
	// _, err = list.Join([]string{"hoge"})
	// if err != nil {
	// 	panic("Failed to join cluster: " + err.Error())
	// }

	// Ask for members of the cluster
	for _, member := range l1.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}

	for _, member := range l2.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}

	// Continue doing whatever you need, memberlist will maintain membership
	// information in the background. Delegates can be used for receiving
	// events when members join or leave.
}
