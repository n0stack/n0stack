package provisioning

import (
	"os"

	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

func getTestNodeAPI() (ppool.NodeServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("NODE_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20181"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return ppool.NewNodeServiceClient(conn), conn, nil
}

func getTestNetworkAPI() (ppool.NetworkServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("NETWORK_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20182"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return ppool.NewNetworkServiceClient(conn), conn, nil
}

func getTestBlockStorageAPI() (pprovisioning.BlockStorageServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("BLOCK_STORAGE_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20183"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return pprovisioning.NewBlockStorageServiceClient(conn), conn, nil
}
