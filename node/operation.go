package node

import (
	"io/ioutil"

	"github.com/satori/go.uuid"
)

func GetNodeUUID() (*uuid.UUID, error) {
	data, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return nil, err
	}

	i, err := uuid.FromString(string(data[:36]))
	if err != nil {
		return nil, err
	}

	return &i, nil
}
