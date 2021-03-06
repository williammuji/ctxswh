// Code generated by protoc-gen-go.
// source: auth.proto
// DO NOT EDIT!

/*
Package authpb is a generated protocol buffer package.

It is generated from these files:
	auth.proto

It has these top-level messages:
	AuthRequest
	AuthResponse
	Account
*/
package authpb

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"

	math "math"

	context "golang.org/x/net/context"

	grpc "google.golang.org/grpc"
)

import gatewaypb "github.com/williammuji/ctxswh/gatewaypb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AuthRequest struct {
	// Types that are valid to be assigned to MessageRequest:
	//	*AuthRequest_LoginRequest
	MessageRequest isAuthRequest_MessageRequest `protobuf_oneof:"MessageRequest"`
}

func (m *AuthRequest) Reset()                    { *m = AuthRequest{} }
func (m *AuthRequest) String() string            { return proto.CompactTextString(m) }
func (*AuthRequest) ProtoMessage()               {}
func (*AuthRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isAuthRequest_MessageRequest interface {
	isAuthRequest_MessageRequest()
}

type AuthRequest_LoginRequest struct {
	LoginRequest *gatewaypb.LoginRequest `protobuf:"bytes,1,opt,name=loginRequest,oneof"`
}

func (*AuthRequest_LoginRequest) isAuthRequest_MessageRequest() {}

func (m *AuthRequest) GetMessageRequest() isAuthRequest_MessageRequest {
	if m != nil {
		return m.MessageRequest
	}
	return nil
}

func (m *AuthRequest) GetLoginRequest() *gatewaypb.LoginRequest {
	if x, ok := m.GetMessageRequest().(*AuthRequest_LoginRequest); ok {
		return x.LoginRequest
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*AuthRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _AuthRequest_OneofMarshaler, _AuthRequest_OneofUnmarshaler, _AuthRequest_OneofSizer, []interface{}{
		(*AuthRequest_LoginRequest)(nil),
	}
}

func _AuthRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*AuthRequest)
	// MessageRequest
	switch x := m.MessageRequest.(type) {
	case *AuthRequest_LoginRequest:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LoginRequest); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("AuthRequest.MessageRequest has unexpected type %T", x)
	}
	return nil
}

func _AuthRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*AuthRequest)
	switch tag {
	case 1: // MessageRequest.loginRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(gatewaypb.LoginRequest)
		err := b.DecodeMessage(msg)
		m.MessageRequest = &AuthRequest_LoginRequest{msg}
		return true, err
	default:
		return false, nil
	}
}

func _AuthRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*AuthRequest)
	// MessageRequest
	switch x := m.MessageRequest.(type) {
	case *AuthRequest_LoginRequest:
		s := proto.Size(x.LoginRequest)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type AuthResponse struct {
	// Types that are valid to be assigned to MessageResponse:
	//	*AuthResponse_LoginResponse
	MessageResponse isAuthResponse_MessageResponse `protobuf_oneof:"MessageResponse"`
}

func (m *AuthResponse) Reset()                    { *m = AuthResponse{} }
func (m *AuthResponse) String() string            { return proto.CompactTextString(m) }
func (*AuthResponse) ProtoMessage()               {}
func (*AuthResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isAuthResponse_MessageResponse interface {
	isAuthResponse_MessageResponse()
}

type AuthResponse_LoginResponse struct {
	LoginResponse *gatewaypb.LoginResponse `protobuf:"bytes,1,opt,name=loginResponse,oneof"`
}

func (*AuthResponse_LoginResponse) isAuthResponse_MessageResponse() {}

func (m *AuthResponse) GetMessageResponse() isAuthResponse_MessageResponse {
	if m != nil {
		return m.MessageResponse
	}
	return nil
}

func (m *AuthResponse) GetLoginResponse() *gatewaypb.LoginResponse {
	if x, ok := m.GetMessageResponse().(*AuthResponse_LoginResponse); ok {
		return x.LoginResponse
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*AuthResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _AuthResponse_OneofMarshaler, _AuthResponse_OneofUnmarshaler, _AuthResponse_OneofSizer, []interface{}{
		(*AuthResponse_LoginResponse)(nil),
	}
}

func _AuthResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*AuthResponse)
	// MessageResponse
	switch x := m.MessageResponse.(type) {
	case *AuthResponse_LoginResponse:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LoginResponse); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("AuthResponse.MessageResponse has unexpected type %T", x)
	}
	return nil
}

func _AuthResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*AuthResponse)
	switch tag {
	case 1: // MessageResponse.loginResponse
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(gatewaypb.LoginResponse)
		err := b.DecodeMessage(msg)
		m.MessageResponse = &AuthResponse_LoginResponse{msg}
		return true, err
	default:
		return false, nil
	}
}

func _AuthResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*AuthResponse)
	// MessageResponse
	switch x := m.MessageResponse.(type) {
	case *AuthResponse_LoginResponse:
		s := proto.Size(x.LoginResponse)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Account struct {
	Username   string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	PasswdHash string `protobuf:"bytes,2,opt,name=passwdHash" json:"passwdHash,omitempty"`
	Email      string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
	Gm         bool   `protobuf:"varint,4,opt,name=gm" json:"gm,omitempty"`
}

func (m *Account) Reset()                    { *m = Account{} }
func (m *Account) String() string            { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()               {}
func (*Account) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Account) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Account) GetPasswdHash() string {
	if m != nil {
		return m.PasswdHash
	}
	return ""
}

func (m *Account) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Account) GetGm() bool {
	if m != nil {
		return m.Gm
	}
	return false
}

func init() {
	proto.RegisterType((*AuthRequest)(nil), "authpb.AuthRequest")
	proto.RegisterType((*AuthResponse)(nil), "authpb.AuthResponse")
	proto.RegisterType((*Account)(nil), "authpb.Account")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AuthService service

type AuthServiceClient interface {
	Serve(ctx context.Context, opts ...grpc.CallOption) (AuthService_ServeClient, error)
}

type authServiceClient struct {
	cc *grpc.ClientConn
}

func NewAuthServiceClient(cc *grpc.ClientConn) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Serve(ctx context.Context, opts ...grpc.CallOption) (AuthService_ServeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_AuthService_serviceDesc.Streams[0], c.cc, "/authpb.AuthService/Serve", opts...)
	if err != nil {
		return nil, err
	}
	x := &authServiceServeClient{stream}
	return x, nil
}

type AuthService_ServeClient interface {
	Send(*AuthRequest) error
	Recv() (*AuthResponse, error)
	grpc.ClientStream
}

type authServiceServeClient struct {
	grpc.ClientStream
}

func (x *authServiceServeClient) Send(m *AuthRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *authServiceServeClient) Recv() (*AuthResponse, error) {
	m := new(AuthResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for AuthService service

type AuthServiceServer interface {
	Serve(AuthService_ServeServer) error
}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
}

func _AuthService_Serve_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AuthServiceServer).Serve(&authServiceServeServer{stream})
}

type AuthService_ServeServer interface {
	Send(*AuthResponse) error
	Recv() (*AuthRequest, error)
	grpc.ServerStream
}

type authServiceServeServer struct {
	grpc.ServerStream
}

func (x *authServiceServeServer) Send(m *AuthResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *authServiceServeServer) Recv() (*AuthRequest, error) {
	m := new(AuthRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "authpb.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Serve",
			Handler:       _AuthService_Serve_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "auth.proto",
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 274 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x51, 0xc1, 0x4e, 0x83, 0x40,
	0x10, 0x2d, 0x68, 0x6b, 0x9d, 0xd6, 0xaa, 0x6b, 0x13, 0x09, 0x87, 0xa6, 0xe1, 0xc4, 0x09, 0x4d,
	0xbd, 0x78, 0x31, 0xb1, 0xbd, 0xc8, 0x41, 0x2f, 0x78, 0x37, 0x59, 0x70, 0x02, 0xc4, 0x02, 0x2b,
	0xb3, 0x2b, 0xfa, 0xf7, 0xa6, 0xec, 0x4a, 0x20, 0xde, 0xe6, 0xbd, 0x37, 0x6f, 0xde, 0xe6, 0x2d,
	0x00, 0x57, 0x32, 0x0b, 0x44, 0x5d, 0xc9, 0x8a, 0x4d, 0x0e, 0xb3, 0x88, 0xdd, 0x55, 0x22, 0xbf,
	0xa9, 0xc9, 0x6e, 0x52, 0x2e, 0xb1, 0xe1, 0x3f, 0x22, 0xfe, 0x9b, 0xf4, 0x9e, 0xf7, 0x06, 0xb3,
	0xad, 0x92, 0x59, 0x84, 0x9f, 0x0a, 0x49, 0xb2, 0x07, 0x98, 0xef, 0xab, 0x34, 0x2f, 0x0d, 0x76,
	0xac, 0xb5, 0xe5, 0xcf, 0x36, 0xd7, 0x41, 0x67, 0x0f, 0x9e, 0x7b, 0x72, 0x38, 0x8a, 0x06, 0xeb,
	0xbb, 0x0b, 0x58, 0xbc, 0x20, 0x11, 0x4f, 0xd1, 0x30, 0x5e, 0x02, 0x73, 0x7d, 0x9f, 0x44, 0x55,
	0x12, 0xb2, 0x47, 0x38, 0x33, 0x0e, 0x4d, 0x98, 0x04, 0xe7, 0x7f, 0x82, 0xd6, 0xc3, 0x51, 0x34,
	0x34, 0xec, 0x2e, 0xe1, 0xbc, 0xcb, 0xd0, 0x94, 0xf7, 0x01, 0x27, 0xdb, 0x24, 0xa9, 0x54, 0x29,
	0x99, 0x0b, 0x53, 0x45, 0x58, 0x97, 0xbc, 0xd0, 0xa7, 0x4f, 0xa3, 0x0e, 0xb3, 0x15, 0x80, 0xe0,
	0x44, 0xcd, 0x7b, 0xc8, 0x29, 0x73, 0xec, 0x56, 0xed, 0x31, 0x6c, 0x09, 0x63, 0x2c, 0x78, 0xbe,
	0x77, 0x8e, 0x5a, 0x49, 0x03, 0xb6, 0x00, 0x3b, 0x2d, 0x9c, 0xe3, 0xb5, 0xe5, 0x4f, 0x23, 0x3b,
	0x2d, 0x36, 0x4f, 0xba, 0xb1, 0x57, 0xac, 0xbf, 0xf2, 0x04, 0xd9, 0x3d, 0x8c, 0x0f, 0x23, 0xb2,
	0xab, 0x40, 0x57, 0x1e, 0xf4, 0xfa, 0x74, 0x97, 0x43, 0xd2, 0xbc, 0x77, 0xe4, 0x5b, 0xb7, 0x56,
	0x3c, 0x69, 0x7f, 0xe0, 0xee, 0x37, 0x00, 0x00, 0xff, 0xff, 0xed, 0x20, 0x45, 0x0f, 0xb7, 0x01,
	0x00, 0x00,
}
