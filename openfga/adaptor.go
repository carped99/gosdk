package openfga

import (
	"context"
	v1 "github.com/carped99/gosdk/openfga/gen/openfga/v1"
	"google.golang.org/grpc"
)

type openFGAServiceClientAdaptor struct {
	client  v1.OpenFGAServiceClient
	storeId string
	modelId string
}

func (s *openFGAServiceClientAdaptor) Read(ctx context.Context, in *v1.ReadRequest, opts ...grpc.CallOption) (*v1.ReadResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.Read(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) Write(ctx context.Context, in *v1.WriteRequest, opts ...grpc.CallOption) (*v1.WriteResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.Write(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) Check(ctx context.Context, in *v1.CheckRequest, opts ...grpc.CallOption) (*v1.CheckResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.Check(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) BatchCheck(ctx context.Context, in *v1.BatchCheckRequest, opts ...grpc.CallOption) (*v1.BatchCheckResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.BatchCheck(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) Expand(ctx context.Context, in *v1.ExpandRequest, opts ...grpc.CallOption) (*v1.ExpandResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.Expand(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ReadAuthorizationModels(ctx context.Context, in *v1.ReadAuthorizationModelsRequest, opts ...grpc.CallOption) (*v1.ReadAuthorizationModelsResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.ReadAuthorizationModels(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ReadAuthorizationModel(ctx context.Context, in *v1.ReadAuthorizationModelRequest, opts ...grpc.CallOption) (*v1.ReadAuthorizationModelResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.ReadAuthorizationModel(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) WriteAuthorizationModel(ctx context.Context, in *v1.WriteAuthorizationModelRequest, opts ...grpc.CallOption) (*v1.WriteAuthorizationModelResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.WriteAuthorizationModel(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) WriteAssertions(ctx context.Context, in *v1.WriteAssertionsRequest, opts ...grpc.CallOption) (*v1.WriteAssertionsResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.WriteAssertions(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ReadAssertions(ctx context.Context, in *v1.ReadAssertionsRequest, opts ...grpc.CallOption) (*v1.ReadAssertionsResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.ReadAssertions(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ReadChanges(ctx context.Context, in *v1.ReadChangesRequest, opts ...grpc.CallOption) (*v1.ReadChangesResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.ReadChanges(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) CreateStore(ctx context.Context, in *v1.CreateStoreRequest, opts ...grpc.CallOption) (*v1.CreateStoreResponse, error) {
	return s.client.CreateStore(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) UpdateStore(ctx context.Context, in *v1.UpdateStoreRequest, opts ...grpc.CallOption) (*v1.UpdateStoreResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.UpdateStore(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) DeleteStore(ctx context.Context, in *v1.DeleteStoreRequest, opts ...grpc.CallOption) (*v1.DeleteStoreResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.DeleteStore(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) GetStore(ctx context.Context, in *v1.GetStoreRequest, opts ...grpc.CallOption) (*v1.GetStoreResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.GetStore(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ListStores(ctx context.Context, in *v1.ListStoresRequest, opts ...grpc.CallOption) (*v1.ListStoresResponse, error) {
	return s.client.ListStores(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) StreamedListObjects(ctx context.Context, in *v1.StreamedListObjectsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[v1.StreamedListObjectsResponse], error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.StreamedListObjects(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ListObjects(ctx context.Context, in *v1.ListObjectsRequest, opts ...grpc.CallOption) (*v1.ListObjectsResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.ListObjects(ctx, in, opts...)
}

func (s *openFGAServiceClientAdaptor) ListUsers(ctx context.Context, in *v1.ListUsersRequest, opts ...grpc.CallOption) (*v1.ListUsersResponse, error) {
	if in != nil && in.StoreId == "" {
		in.StoreId = s.storeId
	}
	return s.client.ListUsers(ctx, in, opts...)
}
