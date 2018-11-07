package provisioning

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Create のときに Node を指定したい時に利用
const AnnotationRequestNodeName = "n0core/provisioning/request_node_name"

// TODO: こいつらは外出しして、共有ライブラリとして使えるようにする
func WrapRollbackError(err error) {
	if err != nil {
		log.Printf("[CRITICAL] Failed to rollback: err=\n%s", err.Error())
	}
}

// WrapGrpcErrorf returns grpc.Errorf
// in the case of 'Internal', logging message because the server has failed
func WrapGrpcErrorf(c codes.Code, format string, a ...interface{}) error {
	err := grpc.Errorf(c, format, a...)

	if c == codes.Internal {
		log.Printf("[WARNING] "+format, a...)
	}

	return err
}
