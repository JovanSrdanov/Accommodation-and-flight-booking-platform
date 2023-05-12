// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.4
// source: reservation_service.proto

package reservation

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

const (
	ReservationService_CreateAvailability_FullMethodName         = "/reservation.ReservationService/CreateAvailability"
	ReservationService_GetAllMy_FullMethodName                   = "/reservation.ReservationService/GetAllMy"
	ReservationService_UpdatePriceAndDate_FullMethodName         = "/reservation.ReservationService/UpdatePriceAndDate"
	ReservationService_CreateReservation_FullMethodName          = "/reservation.ReservationService/CreateReservation"
	ReservationService_CreateAvailabilityBase_FullMethodName     = "/reservation.ReservationService/CreateAvailabilityBase"
	ReservationService_GetAllPendingReservations_FullMethodName  = "/reservation.ReservationService/GetAllPendingReservations"
	ReservationService_GetAllRejectedReservations_FullMethodName = "/reservation.ReservationService/GetAllRejectedReservations"
	ReservationService_RejectReservation_FullMethodName          = "/reservation.ReservationService/RejectReservation"
	ReservationService_AcceptReservation_FullMethodName          = "/reservation.ReservationService/AcceptReservation"
	ReservationService_CancelReservation_FullMethodName          = "/reservation.ReservationService/CancelReservation"
	ReservationService_GetAllReservationsForGuest_FullMethodName = "/reservation.ReservationService/GetAllReservationsForGuest"
)

// ReservationServiceClient is the client API for ReservationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReservationServiceClient interface {
	CreateAvailability(ctx context.Context, in *CreateAvailabilityRequest, opts ...grpc.CallOption) (*CreateAvailabilityResponse, error)
	GetAllMy(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllMyResponse, error)
	UpdatePriceAndDate(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateRequest, error)
	CreateReservation(ctx context.Context, in *CreateReservationRequest, opts ...grpc.CallOption) (*CreateReservationRequest, error)
	CreateAvailabilityBase(ctx context.Context, in *CreateAvailabilityBaseRequest, opts ...grpc.CallOption) (*EmptyRequest, error)
	GetAllPendingReservations(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllPendingReservationsResponse, error)
	GetAllRejectedReservations(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllRejectedReservationsResponse, error)
	RejectReservation(ctx context.Context, in *ChangeStatusRequest, opts ...grpc.CallOption) (*RejectReservationResponse, error)
	AcceptReservation(ctx context.Context, in *ChangeStatusRequest, opts ...grpc.CallOption) (*RejectReservationResponse, error)
	CancelReservation(ctx context.Context, in *ChangeStatusRequest, opts ...grpc.CallOption) (*RejectReservationResponse, error)
	GetAllReservationsForGuest(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllReservationsForGuestResponse, error)
}

type reservationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReservationServiceClient(cc grpc.ClientConnInterface) ReservationServiceClient {
	return &reservationServiceClient{cc}
}

func (c *reservationServiceClient) CreateAvailability(ctx context.Context, in *CreateAvailabilityRequest, opts ...grpc.CallOption) (*CreateAvailabilityResponse, error) {
	out := new(CreateAvailabilityResponse)
	err := c.cc.Invoke(ctx, ReservationService_CreateAvailability_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) GetAllMy(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllMyResponse, error) {
	out := new(GetAllMyResponse)
	err := c.cc.Invoke(ctx, ReservationService_GetAllMy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) UpdatePriceAndDate(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateRequest, error) {
	out := new(UpdateRequest)
	err := c.cc.Invoke(ctx, ReservationService_UpdatePriceAndDate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) CreateReservation(ctx context.Context, in *CreateReservationRequest, opts ...grpc.CallOption) (*CreateReservationRequest, error) {
	out := new(CreateReservationRequest)
	err := c.cc.Invoke(ctx, ReservationService_CreateReservation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) CreateAvailabilityBase(ctx context.Context, in *CreateAvailabilityBaseRequest, opts ...grpc.CallOption) (*EmptyRequest, error) {
	out := new(EmptyRequest)
	err := c.cc.Invoke(ctx, ReservationService_CreateAvailabilityBase_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) GetAllPendingReservations(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllPendingReservationsResponse, error) {
	out := new(GetAllPendingReservationsResponse)
	err := c.cc.Invoke(ctx, ReservationService_GetAllPendingReservations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) GetAllRejectedReservations(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllRejectedReservationsResponse, error) {
	out := new(GetAllRejectedReservationsResponse)
	err := c.cc.Invoke(ctx, ReservationService_GetAllRejectedReservations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) RejectReservation(ctx context.Context, in *ChangeStatusRequest, opts ...grpc.CallOption) (*RejectReservationResponse, error) {
	out := new(RejectReservationResponse)
	err := c.cc.Invoke(ctx, ReservationService_RejectReservation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) AcceptReservation(ctx context.Context, in *ChangeStatusRequest, opts ...grpc.CallOption) (*RejectReservationResponse, error) {
	out := new(RejectReservationResponse)
	err := c.cc.Invoke(ctx, ReservationService_AcceptReservation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) CancelReservation(ctx context.Context, in *ChangeStatusRequest, opts ...grpc.CallOption) (*RejectReservationResponse, error) {
	out := new(RejectReservationResponse)
	err := c.cc.Invoke(ctx, ReservationService_CancelReservation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationServiceClient) GetAllReservationsForGuest(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetAllReservationsForGuestResponse, error) {
	out := new(GetAllReservationsForGuestResponse)
	err := c.cc.Invoke(ctx, ReservationService_GetAllReservationsForGuest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReservationServiceServer is the server API for ReservationService service.
// All implementations must embed UnimplementedReservationServiceServer
// for forward compatibility
type ReservationServiceServer interface {
	CreateAvailability(context.Context, *CreateAvailabilityRequest) (*CreateAvailabilityResponse, error)
	GetAllMy(context.Context, *EmptyRequest) (*GetAllMyResponse, error)
	UpdatePriceAndDate(context.Context, *UpdateRequest) (*UpdateRequest, error)
	CreateReservation(context.Context, *CreateReservationRequest) (*CreateReservationRequest, error)
	CreateAvailabilityBase(context.Context, *CreateAvailabilityBaseRequest) (*EmptyRequest, error)
	GetAllPendingReservations(context.Context, *EmptyRequest) (*GetAllPendingReservationsResponse, error)
	GetAllRejectedReservations(context.Context, *EmptyRequest) (*GetAllRejectedReservationsResponse, error)
	RejectReservation(context.Context, *ChangeStatusRequest) (*RejectReservationResponse, error)
	AcceptReservation(context.Context, *ChangeStatusRequest) (*RejectReservationResponse, error)
	CancelReservation(context.Context, *ChangeStatusRequest) (*RejectReservationResponse, error)
	GetAllReservationsForGuest(context.Context, *EmptyRequest) (*GetAllReservationsForGuestResponse, error)
	mustEmbedUnimplementedReservationServiceServer()
}

// UnimplementedReservationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReservationServiceServer struct {
}

func (UnimplementedReservationServiceServer) CreateAvailability(context.Context, *CreateAvailabilityRequest) (*CreateAvailabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAvailability not implemented")
}
func (UnimplementedReservationServiceServer) GetAllMy(context.Context, *EmptyRequest) (*GetAllMyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMy not implemented")
}
func (UnimplementedReservationServiceServer) UpdatePriceAndDate(context.Context, *UpdateRequest) (*UpdateRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePriceAndDate not implemented")
}
func (UnimplementedReservationServiceServer) CreateReservation(context.Context, *CreateReservationRequest) (*CreateReservationRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReservation not implemented")
}
func (UnimplementedReservationServiceServer) CreateAvailabilityBase(context.Context, *CreateAvailabilityBaseRequest) (*EmptyRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAvailabilityBase not implemented")
}
func (UnimplementedReservationServiceServer) GetAllPendingReservations(context.Context, *EmptyRequest) (*GetAllPendingReservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPendingReservations not implemented")
}
func (UnimplementedReservationServiceServer) GetAllRejectedReservations(context.Context, *EmptyRequest) (*GetAllRejectedReservationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllRejectedReservations not implemented")
}
func (UnimplementedReservationServiceServer) RejectReservation(context.Context, *ChangeStatusRequest) (*RejectReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectReservation not implemented")
}
func (UnimplementedReservationServiceServer) AcceptReservation(context.Context, *ChangeStatusRequest) (*RejectReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptReservation not implemented")
}
func (UnimplementedReservationServiceServer) CancelReservation(context.Context, *ChangeStatusRequest) (*RejectReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelReservation not implemented")
}
func (UnimplementedReservationServiceServer) GetAllReservationsForGuest(context.Context, *EmptyRequest) (*GetAllReservationsForGuestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllReservationsForGuest not implemented")
}
func (UnimplementedReservationServiceServer) mustEmbedUnimplementedReservationServiceServer() {}

// UnsafeReservationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReservationServiceServer will
// result in compilation errors.
type UnsafeReservationServiceServer interface {
	mustEmbedUnimplementedReservationServiceServer()
}

func RegisterReservationServiceServer(s grpc.ServiceRegistrar, srv ReservationServiceServer) {
	s.RegisterService(&ReservationService_ServiceDesc, srv)
}

func _ReservationService_CreateAvailability_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAvailabilityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).CreateAvailability(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_CreateAvailability_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).CreateAvailability(ctx, req.(*CreateAvailabilityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_GetAllMy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).GetAllMy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_GetAllMy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).GetAllMy(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_UpdatePriceAndDate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).UpdatePriceAndDate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_UpdatePriceAndDate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).UpdatePriceAndDate(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_CreateReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReservationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).CreateReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_CreateReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).CreateReservation(ctx, req.(*CreateReservationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_CreateAvailabilityBase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAvailabilityBaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).CreateAvailabilityBase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_CreateAvailabilityBase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).CreateAvailabilityBase(ctx, req.(*CreateAvailabilityBaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_GetAllPendingReservations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).GetAllPendingReservations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_GetAllPendingReservations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).GetAllPendingReservations(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_GetAllRejectedReservations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).GetAllRejectedReservations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_GetAllRejectedReservations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).GetAllRejectedReservations(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_RejectReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).RejectReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_RejectReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).RejectReservation(ctx, req.(*ChangeStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_AcceptReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).AcceptReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_AcceptReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).AcceptReservation(ctx, req.(*ChangeStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_CancelReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).CancelReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_CancelReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).CancelReservation(ctx, req.(*ChangeStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReservationService_GetAllReservationsForGuest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServiceServer).GetAllReservationsForGuest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReservationService_GetAllReservationsForGuest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServiceServer).GetAllReservationsForGuest(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReservationService_ServiceDesc is the grpc.ServiceDesc for ReservationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReservationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "reservation.ReservationService",
	HandlerType: (*ReservationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAvailability",
			Handler:    _ReservationService_CreateAvailability_Handler,
		},
		{
			MethodName: "GetAllMy",
			Handler:    _ReservationService_GetAllMy_Handler,
		},
		{
			MethodName: "UpdatePriceAndDate",
			Handler:    _ReservationService_UpdatePriceAndDate_Handler,
		},
		{
			MethodName: "CreateReservation",
			Handler:    _ReservationService_CreateReservation_Handler,
		},
		{
			MethodName: "CreateAvailabilityBase",
			Handler:    _ReservationService_CreateAvailabilityBase_Handler,
		},
		{
			MethodName: "GetAllPendingReservations",
			Handler:    _ReservationService_GetAllPendingReservations_Handler,
		},
		{
			MethodName: "GetAllRejectedReservations",
			Handler:    _ReservationService_GetAllRejectedReservations_Handler,
		},
		{
			MethodName: "RejectReservation",
			Handler:    _ReservationService_RejectReservation_Handler,
		},
		{
			MethodName: "AcceptReservation",
			Handler:    _ReservationService_AcceptReservation_Handler,
		},
		{
			MethodName: "CancelReservation",
			Handler:    _ReservationService_CancelReservation_Handler,
		},
		{
			MethodName: "GetAllReservationsForGuest",
			Handler:    _ReservationService_GetAllReservationsForGuest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reservation_service.proto",
}
