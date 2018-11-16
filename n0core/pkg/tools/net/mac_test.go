package nettools

import "testing"

func TestGenerateHardwareAddress(t *testing.T) {
	result := "52:54:f5:8c:a4:f2"
	hw := GenerateHardwareAddress("hogehoge")
	if hw.String() != result {
		t.Errorf("Wrong hardware address\n\thave:%s\n\twant:%s", hw, result)
	}
}
