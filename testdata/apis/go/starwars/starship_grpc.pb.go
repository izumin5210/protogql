// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package starwars_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// StarshipServiceClient is the client API for StarshipService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StarshipServiceClient interface {
	GetStarship(ctx context.Context, in *GetStarshipRequest, opts ...grpc.CallOption) (*Starship, error)
}

type starshipServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStarshipServiceClient(cc grpc.ClientConnInterface) StarshipServiceClient {
	return &starshipServiceClient{cc}
}

func (c *starshipServiceClient) GetStarship(ctx context.Context, in *GetStarshipRequest, opts ...grpc.CallOption) (*Starship, error) {
	out := new(Starship)
	err := c.cc.Invoke(ctx, "/testapi.starwars.StarshipService/GetStarship", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StarshipServiceServer is the server API for StarshipService service.
// All implementations must embed UnimplementedStarshipServiceServer
// for forward compatibility
type StarshipServiceServer interface {
	GetStarship(context.Context, *GetStarshipRequest) (*Starship, error)
	mustEmbedUnimplementedStarshipServiceServer()
}

// UnimplementedStarshipServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStarshipServiceServer struct {
}

func (UnimplementedStarshipServiceServer) GetStarship(context.Context, *GetStarshipRequest) (*Starship, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStarship not implemented")
}
func (UnimplementedStarshipServiceServer) mustEmbedUnimplementedStarshipServiceServer() {}

// UnsafeStarshipServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StarshipServiceServer will
// result in compilation errors.
type UnsafeStarshipServiceServer interface {
	mustEmbedUnimplementedStarshipServiceServer()
}

func RegisterStarshipServiceServer(s grpc.ServiceRegistrar, srv StarshipServiceServer) {
	s.RegisterService(&_StarshipService_serviceDesc, srv)
}

func _StarshipService_GetStarship_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStarshipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StarshipServiceServer).GetStarship(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/testapi.starwars.StarshipService/GetStarship",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StarshipServiceServer).GetStarship(ctx, req.(*GetStarshipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StarshipService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "testapi.starwars.StarshipService",
	HandlerType: (*StarshipServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStarship",
			Handler:    _StarshipService_GetStarship_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "starwars/starship.proto",
}
