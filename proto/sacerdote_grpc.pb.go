// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: sacerdote.proto

package proto

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

// SacerdoteClient is the client API for Sacerdote service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SacerdoteClient interface {
	RegistrarTransacao(ctx context.Context, in *PedidoTransacao, opts ...grpc.CallOption) (*ResultadoTransacao, error)
	ConsultarExtrato(ctx context.Context, in *Habitante, opts ...grpc.CallOption) (*Extrato, error)
}

type sacerdoteClient struct {
	cc grpc.ClientConnInterface
}

func NewSacerdoteClient(cc grpc.ClientConnInterface) SacerdoteClient {
	return &sacerdoteClient{cc}
}

func (c *sacerdoteClient) RegistrarTransacao(ctx context.Context, in *PedidoTransacao, opts ...grpc.CallOption) (*ResultadoTransacao, error) {
	out := new(ResultadoTransacao)
	err := c.cc.Invoke(ctx, "/Sacerdote/RegistrarTransacao", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sacerdoteClient) ConsultarExtrato(ctx context.Context, in *Habitante, opts ...grpc.CallOption) (*Extrato, error) {
	out := new(Extrato)
	err := c.cc.Invoke(ctx, "/Sacerdote/ConsultarExtrato", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SacerdoteServer is the server API for Sacerdote service.
// All implementations must embed UnimplementedSacerdoteServer
// for forward compatibility
type SacerdoteServer interface {
	RegistrarTransacao(context.Context, *PedidoTransacao) (*ResultadoTransacao, error)
	ConsultarExtrato(context.Context, *Habitante) (*Extrato, error)
	mustEmbedUnimplementedSacerdoteServer()
}

// UnimplementedSacerdoteServer must be embedded to have forward compatible implementations.
type UnimplementedSacerdoteServer struct {
}

func (UnimplementedSacerdoteServer) RegistrarTransacao(context.Context, *PedidoTransacao) (*ResultadoTransacao, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegistrarTransacao not implemented")
}
func (UnimplementedSacerdoteServer) ConsultarExtrato(context.Context, *Habitante) (*Extrato, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsultarExtrato not implemented")
}
func (UnimplementedSacerdoteServer) mustEmbedUnimplementedSacerdoteServer() {}

// UnsafeSacerdoteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SacerdoteServer will
// result in compilation errors.
type UnsafeSacerdoteServer interface {
	mustEmbedUnimplementedSacerdoteServer()
}

func RegisterSacerdoteServer(s grpc.ServiceRegistrar, srv SacerdoteServer) {
	s.RegisterService(&Sacerdote_ServiceDesc, srv)
}

func _Sacerdote_RegistrarTransacao_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PedidoTransacao)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SacerdoteServer).RegistrarTransacao(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Sacerdote/RegistrarTransacao",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SacerdoteServer).RegistrarTransacao(ctx, req.(*PedidoTransacao))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sacerdote_ConsultarExtrato_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Habitante)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SacerdoteServer).ConsultarExtrato(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Sacerdote/ConsultarExtrato",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SacerdoteServer).ConsultarExtrato(ctx, req.(*Habitante))
	}
	return interceptor(ctx, in, info, handler)
}

// Sacerdote_ServiceDesc is the grpc.ServiceDesc for Sacerdote service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sacerdote_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Sacerdote",
	HandlerType: (*SacerdoteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegistrarTransacao",
			Handler:    _Sacerdote_RegistrarTransacao_Handler,
		},
		{
			MethodName: "ConsultarExtrato",
			Handler:    _Sacerdote_ConsultarExtrato_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sacerdote.proto",
}
