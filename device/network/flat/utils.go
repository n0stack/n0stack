package flat

import (
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
)

func (f flat) getBridgeName(id uuid.UUID) string {
	i := strings.Split(id.String(), "-")
	return fmt.Sprintf("nbr%s", i[0])
}
