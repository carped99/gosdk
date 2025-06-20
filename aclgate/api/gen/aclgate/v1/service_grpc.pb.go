// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: aclgate/v1/service.proto

package aclgatev1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AclGateService_Check_FullMethodName       = "/aclgate.v1.AclGateService/Check"
	AclGateService_BatchCheck_FullMethodName  = "/aclgate.v1.AclGateService/BatchCheck"
	AclGateService_Mutate_FullMethodName      = "/aclgate.v1.AclGateService/Mutate"
	AclGateService_StreamCheck_FullMethodName = "/aclgate.v1.AclGateService/StreamCheck"
	AclGateService_List_FullMethodName        = "/aclgate.v1.AclGateService/List"
	AclGateService_Audit_FullMethodName       = "/aclgate.v1.AclGateService/Audit"
)

// AclGateServiceClient is the client API for AclGateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 권한 확인 서비스
type AclGateServiceClient interface {
	// 단건 권한 확인
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	// 다건 권한 확인
	BatchCheck(ctx context.Context, in *BatchCheckRequest, opts ...grpc.CallOption) (*BatchCheckResponse, error)
	Mutate(ctx context.Context, in *MutateRequest, opts ...grpc.CallOption) (*MutateResponse, error)
	// StreamCheck streams permission check in real-time
	StreamCheck(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamCheckRequest, StreamCheckResponse], error)
	// 권한 목록 조회
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	// 감사 로그 조회
	Audit(ctx context.Context, in *AuditRequest, opts ...grpc.CallOption) (*AuditResponse, error)
}

type aclGateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAclGateServiceClient(cc grpc.ClientConnInterface) AclGateServiceClient {
	return &aclGateServiceClient{cc}
}

func (c *aclGateServiceClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, AclGateService_Check_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aclGateServiceClient) BatchCheck(ctx context.Context, in *BatchCheckRequest, opts ...grpc.CallOption) (*BatchCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchCheckResponse)
	err := c.cc.Invoke(ctx, AclGateService_BatchCheck_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aclGateServiceClient) Mutate(ctx context.Context, in *MutateRequest, opts ...grpc.CallOption) (*MutateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MutateResponse)
	err := c.cc.Invoke(ctx, AclGateService_Mutate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aclGateServiceClient) StreamCheck(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamCheckRequest, StreamCheckResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &AclGateService_ServiceDesc.Streams[0], AclGateService_StreamCheck_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StreamCheckRequest, StreamCheckResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AclGateService_StreamCheckClient = grpc.BidiStreamingClient[StreamCheckRequest, StreamCheckResponse]

func (c *aclGateServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, AclGateService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aclGateServiceClient) Audit(ctx context.Context, in *AuditRequest, opts ...grpc.CallOption) (*AuditResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuditResponse)
	err := c.cc.Invoke(ctx, AclGateService_Audit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AclGateServiceServer is the server API for AclGateService service.
// All implementations must embed UnimplementedAclGateServiceServer
// for forward compatibility.
//
// 권한 확인 서비스
type AclGateServiceServer interface {
	// 단건 권한 확인
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	// 다건 권한 확인
	BatchCheck(context.Context, *BatchCheckRequest) (*BatchCheckResponse, error)
	Mutate(context.Context, *MutateRequest) (*MutateResponse, error)
	// StreamCheck streams permission check in real-time
	StreamCheck(grpc.BidiStreamingServer[StreamCheckRequest, StreamCheckResponse]) error
	// 권한 목록 조회
	List(context.Context, *ListRequest) (*ListResponse, error)
	// 감사 로그 조회
	Audit(context.Context, *AuditRequest) (*AuditResponse, error)
	mustEmbedUnimplementedAclGateServiceServer()
}

// UnimplementedAclGateServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAclGateServiceServer struct{}

func (UnimplementedAclGateServiceServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedAclGateServiceServer) BatchCheck(context.Context, *BatchCheckRequest) (*BatchCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchCheck not implemented")
}
func (UnimplementedAclGateServiceServer) Mutate(context.Context, *MutateRequest) (*MutateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mutate not implemented")
}
func (UnimplementedAclGateServiceServer) StreamCheck(grpc.BidiStreamingServer[StreamCheckRequest, StreamCheckResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StreamCheck not implemented")
}
func (UnimplementedAclGateServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedAclGateServiceServer) Audit(context.Context, *AuditRequest) (*AuditResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Audit not implemented")
}
func (UnimplementedAclGateServiceServer) mustEmbedUnimplementedAclGateServiceServer() {}
func (UnimplementedAclGateServiceServer) testEmbeddedByValue()                        {}

// UnsafeAclGateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AclGateServiceServer will
// result in compilation errors.
type UnsafeAclGateServiceServer interface {
	mustEmbedUnimplementedAclGateServiceServer()
}

func RegisterAclGateServiceServer(s grpc.ServiceRegistrar, srv AclGateServiceServer) {
	// If the following call pancis, it indicates UnimplementedAclGateServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AclGateService_ServiceDesc, srv)
}

func _AclGateService_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AclGateServiceServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AclGateService_Check_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AclGateServiceServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AclGateService_BatchCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AclGateServiceServer).BatchCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AclGateService_BatchCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AclGateServiceServer).BatchCheck(ctx, req.(*BatchCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AclGateService_Mutate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AclGateServiceServer).Mutate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AclGateService_Mutate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AclGateServiceServer).Mutate(ctx, req.(*MutateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AclGateService_StreamCheck_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AclGateServiceServer).StreamCheck(&grpc.GenericServerStream[StreamCheckRequest, StreamCheckResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type AclGateService_StreamCheckServer = grpc.BidiStreamingServer[StreamCheckRequest, StreamCheckResponse]

func _AclGateService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AclGateServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AclGateService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AclGateServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AclGateService_Audit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuditRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AclGateServiceServer).Audit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AclGateService_Audit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AclGateServiceServer).Audit(ctx, req.(*AuditRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AclGateService_ServiceDesc is the grpc.ServiceDesc for AclGateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AclGateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aclgate.v1.AclGateService",
	HandlerType: (*AclGateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _AclGateService_Check_Handler,
		},
		{
			MethodName: "BatchCheck",
			Handler:    _AclGateService_BatchCheck_Handler,
		},
		{
			MethodName: "Mutate",
			Handler:    _AclGateService_Mutate_Handler,
		},
		{
			MethodName: "List",
			Handler:    _AclGateService_List_Handler,
		},
		{
			MethodName: "Audit",
			Handler:    _AclGateService_Audit_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamCheck",
			Handler:       _AclGateService_StreamCheck_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "aclgate/v1/service.proto",
}
