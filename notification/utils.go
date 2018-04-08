package notification

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	notification "github.com/n0stack/proto.go/notification/v0"
)

func Notify(n *notification.Notification) {
	fmt.Printf("%v\n", n)

	if !n.Succeeded {
		panic(n)
	}
}

// MakeNotification return Notification message with time.
// In the future, this method will be hooked some functions such as storing notifications for database.
// Hard-code string without some logic for making searching line easily, when you set arguments for operation. (ゴミ英語、grep検索を容易にするために関数呼び出しをするときにはtypoがあってもいいからstringをハードコードしろってこと、別にtypoがあってもバグにはならないし検索するときには問題ないため)
func MakeNotification(operation string, success bool, description string) *notification.Notification {
	return &notification.Notification{
		Operation:   operation,
		Succeeded:   success,
		Description: description,
		NotifyTime:  ptypes.TimestampNow(),
	}
}
