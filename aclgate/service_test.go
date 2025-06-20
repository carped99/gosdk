package aclgate

import (
	"context"
	"errors"
	"testing"

	v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func (m *MockAclGateServiceClient) List(ctx context.Context, in *v1.ListRequest, opts ...grpc.CallOption) (*v1.ListResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.ListResponse), args.Error(1)
}

func (m *MockAclGateServiceClient) Audit(ctx context.Context, in *v1.AuditRequest, opts ...grpc.CallOption) (*v1.AuditResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.AuditResponse), args.Error(1)
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
			req := CheckRequest{
				Tuple: Tuple{
					ResourceType: tt.resourceType,
					ResourceId:   tt.resourceId,
					SubjectType:  tt.subjectType,
					SubjectId:    tt.subjectId,
					Relation:     tt.relation,
				},
			}
			result, err := service.Check(context.Background(), req)

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
	req := CheckRequest{
		Tuple: Tuple{
			ResourceType: resourceType,
			ResourceId:   resourceId,
			SubjectType:  subjectType,
			SubjectId:    subjectId,
			Relation:     relation,
		},
	}
	_, err := service.Check(context.Background(), req)

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
	req := CheckRequest{
		Tuple: Tuple{
			ResourceType: "document",
			ResourceId:   "doc123",
			SubjectType:  "user",
			SubjectId:    "user456",
			Relation:     "read",
		},
	}
	_, err := service.Check(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, ctx, capturedContext)
	assert.Equal(t, "test-value", capturedContext.Value("test-key"))

	mockClient.AssertExpectations(t)
}

func TestAclGateService_BatchCheck(t *testing.T) {
	tests := []struct {
		name           string
		requests       []CheckRequest
		mockResponse   *v1.BatchCheckResponse
		mockError      error
		expectedResult []BatchCheckResult
		expectedError  bool
	}{
		{
			name: "successful batch check",
			requests: []CheckRequest{
				CheckRequest{
					Tuple: Tuple{
						ResourceType: "document",
						ResourceId:   "doc1",
						SubjectType:  "user",
						SubjectId:    "user1",
						Relation:     "read",
					},
				},
				CheckRequest{
					Tuple: Tuple{
						ResourceType: "document",
						ResourceId:   "doc2",
						SubjectType:  "user",
						SubjectId:    "user1",
						Relation:     "write",
					},
				},
			},
			mockResponse: &v1.BatchCheckResponse{
				Results: []*v1.BatchCheckResult{
					{
						Request: &v1.CheckRequest{
							Tuple: &v1.Tuple{
								Resource: &v1.Resource{Type: "document", Id: "doc1"},
								Subject:  &v1.Subject{Type: "user", Id: "user1"},
								Relation: &v1.Relation{Name: "read"},
							},
						},
						Allowed: true,
					},
					{
						Request: &v1.CheckRequest{
							Tuple: &v1.Tuple{
								Resource: &v1.Resource{Type: "document", Id: "doc2"},
								Subject:  &v1.Subject{Type: "user", Id: "user1"},
								Relation: &v1.Relation{Name: "write"},
							},
						},
						Allowed: false,
					},
				},
			},
			mockError: nil,
			expectedResult: []BatchCheckResult{
				{Request: CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"}}, Allowed: true},
				{Request: CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc2", SubjectType: "user", SubjectId: "user1", Relation: "write"}}, Allowed: false},
			},
			expectedError: false,
		},
		{
			name:           "gRPC error",
			requests:       []CheckRequest{{Tuple: Tuple{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"}}},
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAclGateServiceClient{}
			mockClient.On("BatchCheck", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			service := &aclGateService{client: mockClient}

			result, err := service.BatchCheck(context.Background(), tt.requests)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				for i, expected := range tt.expectedResult {
					assert.Equal(t, expected.Request, result[i].Request)
					assert.Equal(t, expected.Allowed, result[i].Allowed)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_Mutate(t *testing.T) {
	tests := []struct {
		name           string
		writes         []Tuple
		deletes        []Tuple
		mockResponse   *v1.MutateResponse
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name: "successful mutation",
			writes: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
			},
			deletes: []Tuple{
				{ResourceType: "document", ResourceId: "doc2", SubjectType: "user", SubjectId: "user1", Relation: "write"},
			},
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "mutation failed",
			writes: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
			},
			deletes:        []Tuple{},
			mockResponse:   &v1.MutateResponse{Success: false},
			mockError:      nil,
			expectedResult: false,
			expectedError:  false,
		},
		{
			name: "gRPC error",
			writes: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
			},
			deletes:        []Tuple{},
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAclGateServiceClient{}
			mockClient.On("Mutate", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			service := &aclGateService{client: mockClient}

			result, err := service.Mutate(context.Background(), tt.writes, tt.deletes)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_StreamCheck(t *testing.T) {
	mockClient := &MockAclGateServiceClient{}

	// Create a mock stream
	mockStream := &MockBidiStreamingClient{}
	mockClient.On("StreamCheck", mock.Anything, mock.Anything).Return(mockStream, nil)

	service := &aclGateService{client: mockClient}

	stream, err := service.StreamCheck(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, stream)
	assert.Implements(t, (*StreamCheckClient)(nil), stream)

	mockClient.AssertExpectations(t)
}

// MockBidiStreamingClient is a mock implementation of grpc.BidiStreamingClient
type MockBidiStreamingClient struct {
	mock.Mock
}

func (m *MockBidiStreamingClient) Send(msg *v1.StreamCheckRequest) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockBidiStreamingClient) Recv() (*v1.StreamCheckResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.StreamCheckResponse), args.Error(1)
}

func (m *MockBidiStreamingClient) CloseSend() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBidiStreamingClient) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *MockBidiStreamingClient) Header() (metadata.MD, error) {
	args := m.Called()
	return args.Get(0).(metadata.MD), args.Error(1)
}

func (m *MockBidiStreamingClient) Trailer() metadata.MD {
	args := m.Called()
	return args.Get(0).(metadata.MD)
}

func (m *MockBidiStreamingClient) RecvMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockBidiStreamingClient) SendMsg(msg interface{}) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestStreamCheckClient(t *testing.T) {
	mockStream := &MockBidiStreamingClient{}
	client := &streamCheckClient{stream: mockStream}

	// Test Send
	request := CheckRequest{
		Tuple: Tuple{
			ResourceType: "document",
			ResourceId:   "doc1",
			SubjectType:  "user",
			SubjectId:    "user1",
			Relation:     "read",
		},
	}

	mockStream.On("Send", mock.Anything).Return(nil)
	err := client.Send(request)
	assert.NoError(t, err)

	// Test Recv
	expectedResponse := &v1.StreamCheckResponse{
		Allowed: true,
		Reason:  "User has permission",
		Error:   "",
	}
	mockStream.On("Recv").Return(expectedResponse, nil)

	response, err := client.Recv()
	assert.NoError(t, err)
	assert.Equal(t, true, response.Allowed)
	assert.Equal(t, "User has permission", response.Reason)
	assert.Equal(t, "", response.Error)

	// Test Close
	mockStream.On("CloseSend").Return(nil)
	err = client.Close()
	assert.NoError(t, err)

	mockStream.AssertExpectations(t)
}

func TestCheckRequest_Struct(t *testing.T) {
	req := CheckRequest{
		Tuple: Tuple{
			ResourceType: "document",
			ResourceId:   "doc123",
			SubjectType:  "user",
			SubjectId:    "user456",
			Relation:     "read",
		},
	}

	assert.Equal(t, "document", req.Tuple.ResourceType)
	assert.Equal(t, "doc123", req.Tuple.ResourceId)
	assert.Equal(t, "user", req.Tuple.SubjectType)
	assert.Equal(t, "user456", req.Tuple.SubjectId)
	assert.Equal(t, "read", req.Tuple.Relation)
}

func TestTuple_Struct(t *testing.T) {
	tuple := Tuple{
		ResourceType: "document",
		ResourceId:   "doc123",
		SubjectType:  "user",
		SubjectId:    "user456",
		Relation:     "read",
	}

	assert.Equal(t, "document", tuple.ResourceType)
	assert.Equal(t, "doc123", tuple.ResourceId)
	assert.Equal(t, "user", tuple.SubjectType)
	assert.Equal(t, "user456", tuple.SubjectId)
	assert.Equal(t, "read", tuple.Relation)
}

func TestBatchCheckResult_Struct(t *testing.T) {
	request := CheckRequest{
		Tuple: Tuple{
			ResourceType: "document",
			ResourceId:   "doc123",
			SubjectType:  "user",
			SubjectId:    "user456",
			Relation:     "read",
		},
	}

	result := BatchCheckResult{
		Request: request,
		Allowed: true,
		Error:   nil,
	}

	assert.Equal(t, request, result.Request)
	assert.Equal(t, true, result.Allowed)
	assert.Nil(t, result.Error)
}

func TestStreamCheckResponse_Struct(t *testing.T) {
	response := StreamCheckResponse{
		Allowed: true,
		Reason:  "User has permission",
		Error:   "",
	}

	assert.Equal(t, true, response.Allowed)
	assert.Equal(t, "User has permission", response.Reason)
	assert.Equal(t, "", response.Error)
}

func TestAclGateService_Write(t *testing.T) {
	tests := []struct {
		name           string
		tuples         []Tuple
		mockResponse   *v1.MutateResponse
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name: "successful write",
			tuples: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
				{ResourceType: "document", ResourceId: "doc2", SubjectType: "user", SubjectId: "user1", Relation: "write"},
			},
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "empty tuples",
			tuples:         []Tuple{},
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "gRPC error",
			tuples: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
			},
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAclGateServiceClient{}
			mockClient.On("Mutate", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			service := &aclGateService{client: mockClient}

			result, err := service.Write(context.Background(), tt.tuples)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_Delete(t *testing.T) {
	tests := []struct {
		name           string
		tuples         []Tuple
		mockResponse   *v1.MutateResponse
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name: "successful delete",
			tuples: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
				{ResourceType: "document", ResourceId: "doc2", SubjectType: "user", SubjectId: "user1", Relation: "write"},
			},
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "empty tuples",
			tuples:         []Tuple{},
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "gRPC error",
			tuples: []Tuple{
				{ResourceType: "document", ResourceId: "doc1", SubjectType: "user", SubjectId: "user1", Relation: "read"},
			},
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAclGateServiceClient{}
			mockClient.On("Mutate", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			service := &aclGateService{client: mockClient}

			result, err := service.Delete(context.Background(), tt.tuples)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_DeleteResource(t *testing.T) {
	tests := []struct {
		name           string
		resourceType   string
		resourceId     string
		mockResponse   *v1.MutateResponse
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name:           "successful resource deletion",
			resourceType:   "document",
			resourceId:     "doc123",
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "gRPC error",
			resourceType:   "document",
			resourceId:     "doc123",
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAclGateServiceClient{}
			mockClient.On("Mutate", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			service := &aclGateService{client: mockClient}

			result, err := service.DeleteResource(context.Background(), tt.resourceType, tt.resourceId)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestAclGateService_DeleteSubject(t *testing.T) {
	tests := []struct {
		name           string
		subjectType    string
		subjectId      string
		mockResponse   *v1.MutateResponse
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name:           "successful subject deletion",
			subjectType:    "user",
			subjectId:      "user123",
			mockResponse:   &v1.MutateResponse{Success: true},
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:           "gRPC error",
			subjectType:    "user",
			subjectId:      "user123",
			mockResponse:   nil,
			mockError:      errors.New("gRPC error"),
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAclGateServiceClient{}
			mockClient.On("Mutate", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)

			service := &aclGateService{client: mockClient}

			result, err := service.DeleteSubject(context.Background(), tt.subjectType, tt.subjectId)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
