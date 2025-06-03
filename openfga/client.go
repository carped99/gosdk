package openfga

import (
	"context"
	"fmt"
	v1 "github.com/carped99/gosdk/openfga/gen/openfga/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewOpenFGAServiceClient(conn grpc.ClientConnInterface, storeId, modelId string) v1.OpenFGAServiceClient {
	serviceClient := v1.NewOpenFGAServiceClient(conn)

	return &openFGAServiceClientAdaptor{
		client:  serviceClient,
		storeId: storeId,
		modelId: modelId,
	}
}

// validateConnection validates the client connection
func validateConnection(ctx context.Context, conn grpc.ClientConnInterface) error {
	// 헬스 체크
	healthClient := grpc_health_v1.NewHealthClient(conn)
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("OpenFGA server is not serving: %s", resp.Status)
	}
	return nil
}

// GetStoreById validates if the given store ID exists and is accessible
func GetStoreById(client v1.OpenFGAServiceClient, storeId string) (*v1.GetStoreResponse, error) {
	if storeId == "" {
		return nil, fmt.Errorf("store ID cannot be empty")
	}

	return client.GetStore(context.Background(), &v1.GetStoreRequest{
		StoreId: storeId,
	})
}

// GetModelById validates if the given store ID exists and is accessible
func GetModelById(client v1.OpenFGAServiceClient, storeId string, modelId string) (*v1.ReadAuthorizationModelResponse, error) {
	return client.ReadAuthorizationModel(context.Background(), &v1.ReadAuthorizationModelRequest{
		StoreId: storeId,
		Id:      modelId,
	})
}

func ListModels(client v1.OpenFGAServiceClient, storeId string) ([]*v1.AuthorizationModel, error) {
	res, err := client.ReadAuthorizationModels(context.Background(), &v1.ReadAuthorizationModelsRequest{
		StoreId: storeId,
	})

	if err != nil {
		return nil, err
	}

	return res.GetAuthorizationModels(), err
}
