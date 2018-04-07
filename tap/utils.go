package tap

import (
	"fmt"
	"strings"
)

func (f flat) getBridgeName() string {
	return "nbr-flat"
}

func (f flat) getTapName() string {
	i := strings.Split(f.id.String(), "-")
	return fmt.Sprintf("ntap%s", i[0])
}
