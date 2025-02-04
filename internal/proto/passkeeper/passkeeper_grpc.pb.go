// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v4.25.0
// source: passkeeper.proto

package passkeeperv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	PassKeeper_AddEntity_FullMethodName    = "/auth.PassKeeper/AddEntity"
	PassKeeper_UpdateEntity_FullMethodName = "/auth.PassKeeper/UpdateEntity"
	PassKeeper_DeleteEntity_FullMethodName = "/auth.PassKeeper/DeleteEntity"
	PassKeeper_ListEntities_FullMethodName = "/auth.PassKeeper/ListEntities"
	PassKeeper_UploadFile_FullMethodName   = "/auth.PassKeeper/UploadFile"
	PassKeeper_DownloadFile_FullMethodName = "/auth.PassKeeper/DownloadFile"
)

// PassKeeperClient is the client API for PassKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PassKeeperClient interface {
	// AddEntity adds a new entity.
	AddEntity(ctx context.Context, in *AddEntityRequest, opts ...grpc.CallOption) (*AddEntityResponse, error)
	// UpdateEntity updates the entity.
	UpdateEntity(ctx context.Context, in *UpdateEntityRequest, opts ...grpc.CallOption) (*UpdateEntityResponse, error)
	// DeleteEntity deletes the entity.
	DeleteEntity(ctx context.Context, in *DeleteEntityRequest, opts ...grpc.CallOption) (*DeleteEntityResponse, error)
	// ListEntities return all entities list.
	ListEntities(ctx context.Context, in *ListEntitiesRequest, opts ...grpc.CallOption) (*ListEntitiesResponse, error)
	// UploadFile uploads file to the server.
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (PassKeeper_UploadFileClient, error)
	// DownloadFile downloads file from the server.
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (PassKeeper_DownloadFileClient, error)
}

type passKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewPassKeeperClient(cc grpc.ClientConnInterface) PassKeeperClient {
	return &passKeeperClient{cc}
}

func (c *passKeeperClient) AddEntity(ctx context.Context, in *AddEntityRequest, opts ...grpc.CallOption) (*AddEntityResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddEntityResponse)
	err := c.cc.Invoke(ctx, PassKeeper_AddEntity_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passKeeperClient) UpdateEntity(ctx context.Context, in *UpdateEntityRequest, opts ...grpc.CallOption) (*UpdateEntityResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateEntityResponse)
	err := c.cc.Invoke(ctx, PassKeeper_UpdateEntity_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passKeeperClient) DeleteEntity(ctx context.Context, in *DeleteEntityRequest, opts ...grpc.CallOption) (*DeleteEntityResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteEntityResponse)
	err := c.cc.Invoke(ctx, PassKeeper_DeleteEntity_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passKeeperClient) ListEntities(ctx context.Context, in *ListEntitiesRequest, opts ...grpc.CallOption) (*ListEntitiesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListEntitiesResponse)
	err := c.cc.Invoke(ctx, PassKeeper_ListEntities_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passKeeperClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (PassKeeper_UploadFileClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PassKeeper_ServiceDesc.Streams[0], PassKeeper_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &passKeeperUploadFileClient{ClientStream: stream}
	return x, nil
}

type PassKeeper_UploadFileClient interface {
	Send(*UploadFileRequest) error
	CloseAndRecv() (*UploadFileResponse, error)
	grpc.ClientStream
}

type passKeeperUploadFileClient struct {
	grpc.ClientStream
}

func (x *passKeeperUploadFileClient) Send(m *UploadFileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *passKeeperUploadFileClient) CloseAndRecv() (*UploadFileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *passKeeperClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (PassKeeper_DownloadFileClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &PassKeeper_ServiceDesc.Streams[1], PassKeeper_DownloadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &passKeeperDownloadFileClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PassKeeper_DownloadFileClient interface {
	Recv() (*DownloadFileResponse, error)
	grpc.ClientStream
}

type passKeeperDownloadFileClient struct {
	grpc.ClientStream
}

func (x *passKeeperDownloadFileClient) Recv() (*DownloadFileResponse, error) {
	m := new(DownloadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PassKeeperServer is the server API for PassKeeper service.
// All implementations must embed UnimplementedPassKeeperServer
// for forward compatibility
type PassKeeperServer interface {
	// AddEntity adds a new entity.
	AddEntity(context.Context, *AddEntityRequest) (*AddEntityResponse, error)
	// UpdateEntity updates the entity.
	UpdateEntity(context.Context, *UpdateEntityRequest) (*UpdateEntityResponse, error)
	// DeleteEntity deletes the entity.
	DeleteEntity(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error)
	// ListEntities return all entities list.
	ListEntities(context.Context, *ListEntitiesRequest) (*ListEntitiesResponse, error)
	// UploadFile uploads file to the server.
	UploadFile(PassKeeper_UploadFileServer) error
	// DownloadFile downloads file from the server.
	DownloadFile(*DownloadFileRequest, PassKeeper_DownloadFileServer) error
	mustEmbedUnimplementedPassKeeperServer()
}

// UnimplementedPassKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedPassKeeperServer struct {
}

func (UnimplementedPassKeeperServer) AddEntity(context.Context, *AddEntityRequest) (*AddEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEntity not implemented")
}
func (UnimplementedPassKeeperServer) UpdateEntity(context.Context, *UpdateEntityRequest) (*UpdateEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEntity not implemented")
}
func (UnimplementedPassKeeperServer) DeleteEntity(context.Context, *DeleteEntityRequest) (*DeleteEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEntity not implemented")
}
func (UnimplementedPassKeeperServer) ListEntities(context.Context, *ListEntitiesRequest) (*ListEntitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEntities not implemented")
}
func (UnimplementedPassKeeperServer) UploadFile(PassKeeper_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedPassKeeperServer) DownloadFile(*DownloadFileRequest, PassKeeper_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedPassKeeperServer) mustEmbedUnimplementedPassKeeperServer() {}

// UnsafePassKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PassKeeperServer will
// result in compilation errors.
type UnsafePassKeeperServer interface {
	mustEmbedUnimplementedPassKeeperServer()
}

func RegisterPassKeeperServer(s grpc.ServiceRegistrar, srv PassKeeperServer) {
	s.RegisterService(&PassKeeper_ServiceDesc, srv)
}

func _PassKeeper_AddEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PassKeeperServer).AddEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PassKeeper_AddEntity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PassKeeperServer).AddEntity(ctx, req.(*AddEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PassKeeper_UpdateEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PassKeeperServer).UpdateEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PassKeeper_UpdateEntity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PassKeeperServer).UpdateEntity(ctx, req.(*UpdateEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PassKeeper_DeleteEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PassKeeperServer).DeleteEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PassKeeper_DeleteEntity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PassKeeperServer).DeleteEntity(ctx, req.(*DeleteEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PassKeeper_ListEntities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PassKeeperServer).ListEntities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PassKeeper_ListEntities_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PassKeeperServer).ListEntities(ctx, req.(*ListEntitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PassKeeper_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PassKeeperServer).UploadFile(&passKeeperUploadFileServer{ServerStream: stream})
}

type PassKeeper_UploadFileServer interface {
	SendAndClose(*UploadFileResponse) error
	Recv() (*UploadFileRequest, error)
	grpc.ServerStream
}

type passKeeperUploadFileServer struct {
	grpc.ServerStream
}

func (x *passKeeperUploadFileServer) SendAndClose(m *UploadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *passKeeperUploadFileServer) Recv() (*UploadFileRequest, error) {
	m := new(UploadFileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PassKeeper_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PassKeeperServer).DownloadFile(m, &passKeeperDownloadFileServer{ServerStream: stream})
}

type PassKeeper_DownloadFileServer interface {
	Send(*DownloadFileResponse) error
	grpc.ServerStream
}

type passKeeperDownloadFileServer struct {
	grpc.ServerStream
}

func (x *passKeeperDownloadFileServer) Send(m *DownloadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

// PassKeeper_ServiceDesc is the grpc.ServiceDesc for PassKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PassKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.PassKeeper",
	HandlerType: (*PassKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddEntity",
			Handler:    _PassKeeper_AddEntity_Handler,
		},
		{
			MethodName: "UpdateEntity",
			Handler:    _PassKeeper_UpdateEntity_Handler,
		},
		{
			MethodName: "DeleteEntity",
			Handler:    _PassKeeper_DeleteEntity_Handler,
		},
		{
			MethodName: "ListEntities",
			Handler:    _PassKeeper_ListEntities_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _PassKeeper_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _PassKeeper_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "passkeeper.proto",
}
