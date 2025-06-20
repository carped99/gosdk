package aclgate

import (
	"context"
	"errors"
	"testing"

	v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockAclGateServiceClient is a mock implementation of v1.AclGateServiceClient
type MockAclGateServiceClient struct {
	mock.Mock
}

func (m *MockAclGateServiceClient) Check(ctx context.Context, in *v1.CheckRequest, opts ...grpc.CallOption) (*v1.CheckResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.CheckResponse), args.Error(1)
}

func (m *MockAclGateServiceClient) BatchCheck(ctx context.Context, in *v1.BatchCheckRequest, opts ...grpc.CallOption) (*v1.BatchCheckResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.BatchCheckResponse), args.Error(1)
}

func (m *MockAclGateServiceClient) Mutate(ctx context.Context, in *v1.MutateRequest, opts ...grpc.CallOption) (*v1.MutateResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.MutateResponse), args.Error(1)
}

func (m *MockAclGateServiceClient) StreamCheck(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[v1.StreamCheckRequest, v1.StreamCheckResponse], error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(grpc.BidiStreamingClient[v1.StreamCheckRequest, v1.StreamCheckResponse]), args.Error(1)
}

// MockClientConn is a mock implementation of grpc.ClientConnInterface
type MockClientConn struct {
	mock.Mock
}

func (m *MockClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	mockArgs := m.Called(ctx, method, args, reply, opts)
	return mockArgs.Error(0)
}

func (m *MockClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	mockArgs := m.Called(ctx, desc, method, opts)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(grpc.ClientStream), mockArgs.Error(1)
}

func TestNewAclGateService(t *testing.T) {
	tests := []struct {
		name    string
		conn    grpc.ClientConnInterface
		wantErr bool
	}{
		{
			name:    "successful creation",
			conn:    &MockClientConn{},
			wantErr: false,
		},
		{
			name:    "nil connection",
			conn:    nil,
			wantErr: false, // 현재 구현에서는 nil 체크가 없음
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewAclGateService(tt.conn)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.Implements(t, (*AclGateService)(nil), service)
			}
		})
	}
}

func TestAclGateService_Check(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   string
		resourceId     string
		subjectType    string
		subjectId      string
		relation       string
		mockResponse   *v1.CheckResponse
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name:         "successful permission check - allowed",
			resourceType: "document",
			resourceId:   "doc123",
			subjectType:  "user",
			subjectId:    "user456",
			relation:     "read",
			mockResponse: &v1.CheckResponse{
				Allowed: true,
				Reason:  "User has read permission",
			},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:         "successful permission check - denied",
			resourceType: "document",
			resourceId:   "doc123",
			subjectType:  "user",
			subjectId:    "user456",
			relation:     "write",
			mockResponse: &v1.CheckResponse{
				Allowed: false,
				Reason:  "User does not have write permission",
			},
			mockError:      nil,
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:           "gRPC error",
			resourceType:   "document",
			resourceId:     "doc123",
			subjectType:    "user",
			subjectId:      "user456",
			relation:       "read",
			mockResponse:   nil,
			mockError:      errors.New("gRPC connection error"),
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:           "empty parameters",
			resourceType:   "",
			resourceId:     "",
			subjectType:    "",
			subjectId:      "",
			relation:       "",
			mockResponse:   &v1.CheckResponse{Allowed: false},
			mockError:      nil,
			expectedResult: false,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockAclGateServiceClient{}

			// Setup mock expectations
			mockClient.On("Check", mock.Anything, mock.MatchedBy(func(req *v1.CheckRequest) bool {
				return req.Tuple.Resource.Type == tt.resourceType &&
					req.Tuple.Resource.Id == tt.resourceId &&
					req.Tuple.Subject.Type == tt.subjectType &&
					req.Tuple.Subject.Id == tt.subjectId &&
					req.Tuple.Relation.Name == tt.relation
			}), mock.Anything).Return(tt.mockResponse, tt.mockError)

			// Create service with mock client
			service := &aclGateService{
				client: mockClient,
			}

			// Execute test
			result, err := service.Check(context.Background(), tt.resourceType, tt.resourceId, tt.subjectType, tt.subjectId, tt.relation)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Verify mock was called
			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_MustCheck(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   string
		resourceId     string
		subjectType    string
		subjectId      string
		relation       string
		mockResponse   *v1.CheckResponse
		mockError      error
		expectedResult bool
	}{
		{
			name:           "successful check - allowed",
			resourceType:   "document",
			resourceId:     "doc123",
			subjectType:    "user",
			subjectId:      "user456",
			relation:       "read",
			mockResponse:   &v1.CheckResponse{Allowed: true},
			mockError:      nil,
			expectedResult: true,
		},
		{
			name:           "successful check - denied",
			resourceType:   "document",
			resourceId:     "doc123",
			subjectType:    "user",
			subjectId:      "user456",
			relation:       "write",
			mockResponse:   &v1.CheckResponse{Allowed: false},
			mockError:      nil,
			expectedResult: false,
		},
		{
			name:           "error case - returns false",
			resourceType:   "document",
			resourceId:     "doc123",
			subjectType:    "user",
			subjectId:      "user456",
			relation:       "read",
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockAclGateServiceClient{}

			// Setup mock expectations
			mockClient.On("Check", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			// Create service with mock client
			service := &aclGateService{
				client: mockClient,
			}

			// Execute test
			result := service.MustCheck(context.Background(), tt.resourceType, tt.resourceId, tt.subjectType, tt.subjectId, tt.relation)

			// Assertions
			assert.Equal(t, tt.expectedResult, result)

			// Verify mock was called
			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_Check_RequestValidation(t *testing.T) {
	mockClient := &MockAclGateServiceClient{}

	// Setup mock to capture the request
	var capturedRequest *v1.CheckRequest
	mockClient.On("Check", mock.Anything, mock.Anything, mock.Anything).Return(
		&v1.CheckResponse{Allowed: true}, nil,
	).Run(func(args mock.Arguments) {
		capturedRequest = args.Get(1).(*v1.CheckRequest)
	})

	service := &aclGateService{client: mockClient}

	// Test parameters
	resourceType := "document"
	resourceId := "doc123"
	subjectType := "user"
	subjectId := "user456"
	relation := "read"

	// Execute
	_, err := service.Check(context.Background(), resourceType, resourceId, subjectType, subjectId, relation)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Tuple)
	assert.NotNil(t, capturedRequest.Tuple.Resource)
	assert.NotNil(t, capturedRequest.Tuple.Subject)
	assert.NotNil(t, capturedRequest.Tuple.Relation)

	assert.Equal(t, resourceType, capturedRequest.Tuple.Resource.Type)
	assert.Equal(t, resourceId, capturedRequest.Tuple.Resource.Id)
	assert.Equal(t, subjectType, capturedRequest.Tuple.Subject.Type)
	assert.Equal(t, subjectId, capturedRequest.Tuple.Subject.Id)
	assert.Equal(t, relation, capturedRequest.Tuple.Relation.Name)

	mockClient.AssertExpectations(t)
}

func TestAclGateService_InterfaceCompliance(t *testing.T) {
	// Test that aclGateService implements AclGateService interface
	var _ AclGateService = (*aclGateService)(nil)
}

func TestAclGateService_ContextPropagation(t *testing.T) {
	mockClient := &MockAclGateServiceClient{}

	// Setup mock to capture context
	var capturedContext context.Context
	mockClient.On("Check", mock.Anything, mock.Anything, mock.Anything).Return(
		&v1.CheckResponse{Allowed: true}, nil,
	).Run(func(args mock.Arguments) {
		capturedContext = args.Get(0).(context.Context)
	})

	service := &aclGateService{client: mockClient}

	// Create a context with a value
	ctx := context.WithValue(context.Background(), "test-key", "test-value")

	// Execute
	_, err := service.Check(ctx, "document", "doc123", "user", "user456", "read")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, ctx, capturedContext)
	assert.Equal(t, "test-value", capturedContext.Value("test-key"))

	mockClient.AssertExpectations(t)
}
