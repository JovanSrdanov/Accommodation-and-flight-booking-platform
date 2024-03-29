// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: rating_service.proto

package rating

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RatingServiceClient is the client API for RatingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RatingServiceClient interface {
	RateAccommodation(ctx context.Context, in *RateAccommodationRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	GetRatingForAccommodation(ctx context.Context, in *RatingForAccommodationRequest, opts ...grpc.CallOption) (*RatingForAccommodationResponse, error)
	DeleteRatingForAccommodation(ctx context.Context, in *RatingForAccommodationRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	RateHost(ctx context.Context, in *RateHostRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	GetRatingForHost(ctx context.Context, in *RatingForHostRequest, opts ...grpc.CallOption) (*RatingForHostResponse, error)
	DeleteRatingForHost(ctx context.Context, in *RatingForHostRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	GetRecommendedAccommodations(ctx context.Context, in *RecommendedAccommodationsRequest, opts ...grpc.CallOption) (*RecommendedAccommodationsResponse, error)
	CalculateRatingForHost(ctx context.Context, in *RatingForHostRequest, opts ...grpc.CallOption) (*CalculateRatingForHostResponse, error)
	CalculateRatingForAccommodation(ctx context.Context, in *RatingForAccommodationRequest, opts ...grpc.CallOption) (*CalculateRatingForAccommodationResponse, error)
	GetRatingGuestGaveHost(ctx context.Context, in *GetRatingGuestGaveHostRequest, opts ...grpc.CallOption) (*GetRatingGuestGaveHostResponse, error)
	GetRatingGuestGaveAccommodation(ctx context.Context, in *GetRatingGuestGaveAccommodationRequest, opts ...grpc.CallOption) (*GetRatingGuestGaveAccommodationResponse, error)
}

type ratingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRatingServiceClient(cc grpc.ClientConnInterface) RatingServiceClient {
	return &ratingServiceClient{cc}
}

func (c *ratingServiceClient) RateAccommodation(ctx context.Context, in *RateAccommodationRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/RateAccommodation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) GetRatingForAccommodation(ctx context.Context, in *RatingForAccommodationRequest, opts ...grpc.CallOption) (*RatingForAccommodationResponse, error) {
	out := new(RatingForAccommodationResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/GetRatingForAccommodation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) DeleteRatingForAccommodation(ctx context.Context, in *RatingForAccommodationRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/DeleteRatingForAccommodation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) RateHost(ctx context.Context, in *RateHostRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/RateHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) GetRatingForHost(ctx context.Context, in *RatingForHostRequest, opts ...grpc.CallOption) (*RatingForHostResponse, error) {
	out := new(RatingForHostResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/GetRatingForHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) DeleteRatingForHost(ctx context.Context, in *RatingForHostRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/DeleteRatingForHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) GetRecommendedAccommodations(ctx context.Context, in *RecommendedAccommodationsRequest, opts ...grpc.CallOption) (*RecommendedAccommodationsResponse, error) {
	out := new(RecommendedAccommodationsResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/GetRecommendedAccommodations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) CalculateRatingForHost(ctx context.Context, in *RatingForHostRequest, opts ...grpc.CallOption) (*CalculateRatingForHostResponse, error) {
	out := new(CalculateRatingForHostResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/CalculateRatingForHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) CalculateRatingForAccommodation(ctx context.Context, in *RatingForAccommodationRequest, opts ...grpc.CallOption) (*CalculateRatingForAccommodationResponse, error) {
	out := new(CalculateRatingForAccommodationResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/CalculateRatingForAccommodation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) GetRatingGuestGaveHost(ctx context.Context, in *GetRatingGuestGaveHostRequest, opts ...grpc.CallOption) (*GetRatingGuestGaveHostResponse, error) {
	out := new(GetRatingGuestGaveHostResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/GetRatingGuestGaveHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ratingServiceClient) GetRatingGuestGaveAccommodation(ctx context.Context, in *GetRatingGuestGaveAccommodationRequest, opts ...grpc.CallOption) (*GetRatingGuestGaveAccommodationResponse, error) {
	out := new(GetRatingGuestGaveAccommodationResponse)
	err := c.cc.Invoke(ctx, "/rating.RatingService/GetRatingGuestGaveAccommodation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RatingServiceServer is the server API for RatingService service.
// All implementations must embed UnimplementedRatingServiceServer
// for forward compatibility
type RatingServiceServer interface {
	RateAccommodation(context.Context, *RateAccommodationRequest) (*EmptyResponse, error)
	GetRatingForAccommodation(context.Context, *RatingForAccommodationRequest) (*RatingForAccommodationResponse, error)
	DeleteRatingForAccommodation(context.Context, *RatingForAccommodationRequest) (*SimpleResponse, error)
	RateHost(context.Context, *RateHostRequest) (*EmptyResponse, error)
	GetRatingForHost(context.Context, *RatingForHostRequest) (*RatingForHostResponse, error)
	DeleteRatingForHost(context.Context, *RatingForHostRequest) (*SimpleResponse, error)
	GetRecommendedAccommodations(context.Context, *RecommendedAccommodationsRequest) (*RecommendedAccommodationsResponse, error)
	CalculateRatingForHost(context.Context, *RatingForHostRequest) (*CalculateRatingForHostResponse, error)
	CalculateRatingForAccommodation(context.Context, *RatingForAccommodationRequest) (*CalculateRatingForAccommodationResponse, error)
	GetRatingGuestGaveHost(context.Context, *GetRatingGuestGaveHostRequest) (*GetRatingGuestGaveHostResponse, error)
	GetRatingGuestGaveAccommodation(context.Context, *GetRatingGuestGaveAccommodationRequest) (*GetRatingGuestGaveAccommodationResponse, error)
	mustEmbedUnimplementedRatingServiceServer()
}

// UnimplementedRatingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRatingServiceServer struct {
}

func (UnimplementedRatingServiceServer) RateAccommodation(context.Context, *RateAccommodationRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RateAccommodation not implemented")
}
func (UnimplementedRatingServiceServer) GetRatingForAccommodation(context.Context, *RatingForAccommodationRequest) (*RatingForAccommodationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRatingForAccommodation not implemented")
}
func (UnimplementedRatingServiceServer) DeleteRatingForAccommodation(context.Context, *RatingForAccommodationRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRatingForAccommodation not implemented")
}
func (UnimplementedRatingServiceServer) RateHost(context.Context, *RateHostRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RateHost not implemented")
}
func (UnimplementedRatingServiceServer) GetRatingForHost(context.Context, *RatingForHostRequest) (*RatingForHostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRatingForHost not implemented")
}
func (UnimplementedRatingServiceServer) DeleteRatingForHost(context.Context, *RatingForHostRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRatingForHost not implemented")
}
func (UnimplementedRatingServiceServer) GetRecommendedAccommodations(context.Context, *RecommendedAccommodationsRequest) (*RecommendedAccommodationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecommendedAccommodations not implemented")
}
func (UnimplementedRatingServiceServer) CalculateRatingForHost(context.Context, *RatingForHostRequest) (*CalculateRatingForHostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateRatingForHost not implemented")
}
func (UnimplementedRatingServiceServer) CalculateRatingForAccommodation(context.Context, *RatingForAccommodationRequest) (*CalculateRatingForAccommodationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateRatingForAccommodation not implemented")
}
func (UnimplementedRatingServiceServer) GetRatingGuestGaveHost(context.Context, *GetRatingGuestGaveHostRequest) (*GetRatingGuestGaveHostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRatingGuestGaveHost not implemented")
}
func (UnimplementedRatingServiceServer) GetRatingGuestGaveAccommodation(context.Context, *GetRatingGuestGaveAccommodationRequest) (*GetRatingGuestGaveAccommodationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRatingGuestGaveAccommodation not implemented")
}
func (UnimplementedRatingServiceServer) mustEmbedUnimplementedRatingServiceServer() {}

// UnsafeRatingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RatingServiceServer will
// result in compilation errors.
type UnsafeRatingServiceServer interface {
	mustEmbedUnimplementedRatingServiceServer()
}

func RegisterRatingServiceServer(s grpc.ServiceRegistrar, srv RatingServiceServer) {
	s.RegisterService(&RatingService_ServiceDesc, srv)
}

func _RatingService_RateAccommodation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateAccommodationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).RateAccommodation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/RateAccommodation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).RateAccommodation(ctx, req.(*RateAccommodationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_GetRatingForAccommodation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatingForAccommodationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).GetRatingForAccommodation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/GetRatingForAccommodation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).GetRatingForAccommodation(ctx, req.(*RatingForAccommodationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_DeleteRatingForAccommodation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatingForAccommodationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).DeleteRatingForAccommodation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/DeleteRatingForAccommodation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).DeleteRatingForAccommodation(ctx, req.(*RatingForAccommodationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_RateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).RateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/RateHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).RateHost(ctx, req.(*RateHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_GetRatingForHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatingForHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).GetRatingForHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/GetRatingForHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).GetRatingForHost(ctx, req.(*RatingForHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_DeleteRatingForHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatingForHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).DeleteRatingForHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/DeleteRatingForHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).DeleteRatingForHost(ctx, req.(*RatingForHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_GetRecommendedAccommodations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecommendedAccommodationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).GetRecommendedAccommodations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/GetRecommendedAccommodations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).GetRecommendedAccommodations(ctx, req.(*RecommendedAccommodationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_CalculateRatingForHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatingForHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).CalculateRatingForHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/CalculateRatingForHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).CalculateRatingForHost(ctx, req.(*RatingForHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_CalculateRatingForAccommodation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RatingForAccommodationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).CalculateRatingForAccommodation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/CalculateRatingForAccommodation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).CalculateRatingForAccommodation(ctx, req.(*RatingForAccommodationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_GetRatingGuestGaveHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRatingGuestGaveHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).GetRatingGuestGaveHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/GetRatingGuestGaveHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).GetRatingGuestGaveHost(ctx, req.(*GetRatingGuestGaveHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RatingService_GetRatingGuestGaveAccommodation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRatingGuestGaveAccommodationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RatingServiceServer).GetRatingGuestGaveAccommodation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rating.RatingService/GetRatingGuestGaveAccommodation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RatingServiceServer).GetRatingGuestGaveAccommodation(ctx, req.(*GetRatingGuestGaveAccommodationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RatingService_ServiceDesc is the grpc.ServiceDesc for RatingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RatingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rating.RatingService",
	HandlerType: (*RatingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RateAccommodation",
			Handler:    _RatingService_RateAccommodation_Handler,
		},
		{
			MethodName: "GetRatingForAccommodation",
			Handler:    _RatingService_GetRatingForAccommodation_Handler,
		},
		{
			MethodName: "DeleteRatingForAccommodation",
			Handler:    _RatingService_DeleteRatingForAccommodation_Handler,
		},
		{
			MethodName: "RateHost",
			Handler:    _RatingService_RateHost_Handler,
		},
		{
			MethodName: "GetRatingForHost",
			Handler:    _RatingService_GetRatingForHost_Handler,
		},
		{
			MethodName: "DeleteRatingForHost",
			Handler:    _RatingService_DeleteRatingForHost_Handler,
		},
		{
			MethodName: "GetRecommendedAccommodations",
			Handler:    _RatingService_GetRecommendedAccommodations_Handler,
		},
		{
			MethodName: "CalculateRatingForHost",
			Handler:    _RatingService_CalculateRatingForHost_Handler,
		},
		{
			MethodName: "CalculateRatingForAccommodation",
			Handler:    _RatingService_CalculateRatingForAccommodation_Handler,
		},
		{
			MethodName: "GetRatingGuestGaveHost",
			Handler:    _RatingService_GetRatingGuestGaveHost_Handler,
		},
		{
			MethodName: "GetRatingGuestGaveAccommodation",
			Handler:    _RatingService_GetRatingGuestGaveAccommodation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rating_service.proto",
}
