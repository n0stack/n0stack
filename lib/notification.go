package lib

import (
	"github.com/golang/protobuf/ptypes"

	n0stack "github.com/n0stack/proto"
)

// MakeNotification return Notification message with time.
// In the future, this method will be hooked some functions such as storing notifications for database.
// Hard-code string without some logic for making searching line easily, when you set arguments for operation. (ゴミ英語、grep検索を容易にするために関数呼び出しをするときにはtypoがあってもいいからstringをハードコードしろってこと、別にtypoがあってもバグにはならないし検索するときには問題ないため)
func MakeNotification(operation string, success bool, description string) *n0stack.Notification {
	return &n0stack.Notification{
		Operation:   operation,
		Success:     success,
		Description: description,
		NotifiedAt:  ptypes.TimestampNow(),
	}
}
