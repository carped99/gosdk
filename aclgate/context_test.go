package aclgate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAclGateService is a mock implementation of AclGateService for testing
type MockAclGateService struct {
	mock.Mock
}

func (m *MockAclGateService) Check(ctx context.Context, req CheckRequest) (bool, error) {
	args := m.Called(ctx, req)
	return args.Bool(0), args.Error(1)
}

func (m *MockAclGateService) Write(ctx context.Context, tuples []Tuple) (bool, error) {
	args := m.Called(ctx, tuples)
	return args.Bool(0), args.Error(1)
}

func (m *MockAclGateService) Delete(ctx context.Context, tuples []Tuple) (bool, error) {
	args := m.Called(ctx, tuples)
	return args.Bool(0), args.Error(1)
}

func (m *MockAclGateService) DeleteResource(ctx context.Context, resourceType, resourceId string) (bool, error) {
	args := m.Called(ctx, resourceType, resourceId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAclGateService) DeleteSubject(ctx context.Context, subjectType, subjectId string) (bool, error) {
	args := m.Called(ctx, subjectType, subjectId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAclGateService) BatchCheck(ctx context.Context, requests []CheckRequest) ([]BatchCheckResult, error) {
	args := m.Called(ctx, requests)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]BatchCheckResult), args.Error(1)
}

func (m *MockAclGateService) Mutate(ctx context.Context, writes []Tuple, deletes []Tuple) (bool, error) {
	args := m.Called(ctx, writes, deletes)
	return args.Bool(0), args.Error(1)
}

func (m *MockAclGateService) StreamCheck(ctx context.Context) (StreamCheckClient, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(StreamCheckClient), args.Error(1)
}

func (m *MockAclGateService) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListResponse), args.Error(1)
}

func (m *MockAclGateService) Audit(ctx context.Context, req AuditRequest) (*AuditResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AuditResponse), args.Error(1)
}

func TestNewContext(t *testing.T) {
	mockService := &MockAclGateService{}
	ctx := context.Background()

	// Test creating new context with service
	newCtx := NewContext(ctx, mockService)

	// Verify service is stored in context
	retrievedService, err := FromContext(newCtx)
	assert.NoError(t, err)
	assert.Equal(t, mockService, retrievedService)

	// Verify original context is not modified
	_, err = FromContext(ctx)
	assert.Error(t, err)
}

func TestFromContext(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedError bool
		expectedNil   bool
	}{
		{
			name: "service found in context",
			setupContext: func() context.Context {
				mockService := &MockAclGateService{}
				return NewContext(context.Background(), mockService)
			},
			expectedError: false,
			expectedNil:   false,
		},
		{
			name: "service not found in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedError: true,
			expectedNil:   true,
		},
		{
			name: "wrong type in context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), aclServiceKey, "not a service")
			},
			expectedError: true,
			expectedNil:   true,
		},
		{
			name: "nil context",
			setupContext: func() context.Context {
				return nil
			},
			expectedError: true,
			expectedNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			service, err := FromContext(ctx)

			if tt.expectedError {
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

func TestMustFromContext(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func() context.Context
		shouldPanic  bool
	}{
		{
			name: "service found - no panic",
			setupContext: func() context.Context {
				mockService := &MockAclGateService{}
				return NewContext(context.Background(), mockService)
			},
			shouldPanic: false,
		},
		{
			name: "service not found - panic",
			setupContext: func() context.Context {
				return context.Background()
			},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()

			if tt.shouldPanic {
				assert.Panics(t, func() {
					MustFromContext(ctx)
				})
			} else {
				assert.NotPanics(t, func() {
					service := MustFromContext(ctx)
					assert.NotNil(t, service)
					assert.Implements(t, (*AclGateService)(nil), service)
				})
			}
		})
	}
}

func TestWithContext_Middleware(t *testing.T) {
	mockService := &MockAclGateService{}
	middleware := WithContext(mockService)

	// Create a test handler that checks if service is in context
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		// Verify service is available in context
		service, err := FromContext(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, mockService, service)

		w.WriteHeader(http.StatusOK)
	})

	// Create wrapped handler
	wrappedHandler := middleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify handler was called
	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify original request context is not modified
	_, err := FromContext(req.Context())
	assert.Error(t, err)
}

func TestCheckPermission(t *testing.T) {
	tests := []struct {
		name           string
		request        CheckRequest
		setupContext   func() context.Context
		mockResult     bool
		mockError      error
		expectedResult bool
		expectedError  bool
	}{
		{
			name:    "successful permission check",
			request: CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "read"}},
			setupContext: func() context.Context {
				mockService := &MockAclGateService{}
				mockService.On("Check", mock.Anything, CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "read"}}).
					Return(true, nil)
				return NewContext(context.Background(), mockService)
			},
			mockResult:     true,
			mockError:      nil,
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:    "permission denied",
			request: CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "write"}},
			setupContext: func() context.Context {
				mockService := &MockAclGateService{}
				mockService.On("Check", mock.Anything, CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "write"}}).
					Return(false, nil)
				return NewContext(context.Background(), mockService)
			},
			mockResult:     false,
			mockError:      nil,
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:    "service error",
			request: CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "read"}},
			setupContext: func() context.Context {
				mockService := &MockAclGateService{}
				mockService.On("Check", mock.Anything, CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "read"}}).
					Return(false, assert.AnError)
				return NewContext(context.Background(), mockService)
			},
			mockResult:     false,
			mockError:      assert.AnError,
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:    "service not in context",
			request: CheckRequest{Tuple: Tuple{ResourceType: "document", ResourceId: "doc123", SubjectType: "user", SubjectId: "user456", Relation: "read"}},
			setupContext: func() context.Context {
				return context.Background()
			},
			mockResult:     false,
			mockError:      nil,
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()

			result, err := CheckPermission(ctx, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestContextKey_Uniqueness(t *testing.T) {
	// Test that context key is unique and doesn't conflict with other keys
	ctx := context.Background()

	// Add different values with different keys
	ctx1 := context.WithValue(ctx, "key1", "value1")
	ctx2 := context.WithValue(ctx1, aclServiceKey, "value2")
	ctx3 := context.WithValue(ctx2, "key3", "value3")

	// Verify values are stored separately
	assert.Equal(t, "value1", ctx3.Value("key1"))
	assert.Equal(t, "value2", ctx3.Value(aclServiceKey))
	assert.Equal(t, "value3", ctx3.Value("key3"))
}

func TestContextKey_TypeSafety(t *testing.T) {
	// Test that context key type prevents accidental collisions
	ctx := context.Background()

	// Add string key with same value as aclServiceKey
	stringKey := "acl_service"
	ctx1 := context.WithValue(ctx, stringKey, "string value")

	// Add aclServiceKey
	mockService := &MockAclGateService{}
	ctx2 := NewContext(ctx1, mockService)

	// Verify they are different
	assert.Equal(t, "string value", ctx2.Value(stringKey))

	service, err := FromContext(ctx2)
	assert.NoError(t, err)
	assert.Equal(t, mockService, service)
}

func TestWithContext_MultipleMiddleware(t *testing.T) {
	mockService1 := &MockAclGateService{}
	mockService2 := &MockAclGateService{}

	middleware1 := WithContext(mockService1)
	middleware2 := WithContext(mockService2)

	// Create a test handler
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		// Should get the last service added (mockService2)
		service, err := FromContext(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, mockService2, service)

		w.WriteHeader(http.StatusOK)
	})

	// Apply multiple middleware
	wrappedHandler := middleware1(middleware2(testHandler))

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify handler was called
	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, w.Code)
}
