// Code generated by protoc-gen-go.
// source: login.proto
// DO NOT EDIT!

/*
Package loginpb is a generated protocol buffer package.

It is generated from these files:
	login.proto

It has these top-level messages:
	LoginRequest
	LoginResponse
	IpCount
	Empty
*/
package loginpb

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"

	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LoginRequest struct {
	// Types that are valid to be assigned to MessageRequest:
	//	*LoginRequest_IpCount
	MessageRequest isLoginRequest_MessageRequest `protobuf_oneof:"MessageRequest"`
}

func (m *LoginRequest) Reset()                    { *m = LoginRequest{} }
func (m *LoginRequest) String() string            { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()               {}
func (*LoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isLoginRequest_MessageRequest interface {
	isLoginRequest_MessageRequest()
}

type LoginRequest_IpCount struct {
	IpCount *IpCount `protobuf:"bytes,1,opt,name=ipCount,oneof"`
}

func (*LoginRequest_IpCount) isLoginRequest_MessageRequest() {}

func (m *LoginRequest) GetMessageRequest() isLoginRequest_MessageRequest {
	if m != nil {
		return m.MessageRequest
	}
	return nil
}

func (m *LoginRequest) GetIpCount() *IpCount {
	if x, ok := m.GetMessageRequest().(*LoginRequest_IpCount); ok {
		return x.IpCount
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LoginRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoginRequest_OneofMarshaler, _LoginRequest_OneofUnmarshaler, _LoginRequest_OneofSizer, []interface{}{
		(*LoginRequest_IpCount)(nil),
	}
}

func _LoginRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoginRequest)
	// MessageRequest
	switch x := m.MessageRequest.(type) {
	case *LoginRequest_IpCount:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.IpCount); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoginRequest.MessageRequest has unexpected type %T", x)
	}
	return nil
}

func _LoginRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoginRequest)
	switch tag {
	case 1: // MessageRequest.ipCount
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(IpCount)
		err := b.DecodeMessage(msg)
		m.MessageRequest = &LoginRequest_IpCount{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoginRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoginRequest)
	// MessageRequest
	switch x := m.MessageRequest.(type) {
	case *LoginRequest_IpCount:
		s := proto.Size(x.IpCount)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type LoginResponse struct {
	// Types that are valid to be assigned to MessageResponse:
	//	*LoginResponse_Empty
	MessageResponse isLoginResponse_MessageResponse `protobuf_oneof:"MessageResponse"`
}

func (m *LoginResponse) Reset()                    { *m = LoginResponse{} }
func (m *LoginResponse) String() string            { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()               {}
func (*LoginResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isLoginResponse_MessageResponse interface {
	isLoginResponse_MessageResponse()
}

type LoginResponse_Empty struct {
	Empty *Empty `protobuf:"bytes,1,opt,name=empty,oneof"`
}

func (*LoginResponse_Empty) isLoginResponse_MessageResponse() {}

func (m *LoginResponse) GetMessageResponse() isLoginResponse_MessageResponse {
	if m != nil {
		return m.MessageResponse
	}
	return nil
}

func (m *LoginResponse) GetEmpty() *Empty {
	if x, ok := m.GetMessageResponse().(*LoginResponse_Empty); ok {
		return x.Empty
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LoginResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoginResponse_OneofMarshaler, _LoginResponse_OneofUnmarshaler, _LoginResponse_OneofSizer, []interface{}{
		(*LoginResponse_Empty)(nil),
	}
}

func _LoginResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoginResponse)
	// MessageResponse
	switch x := m.MessageResponse.(type) {
	case *LoginResponse_Empty:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Empty); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoginResponse.MessageResponse has unexpected type %T", x)
	}
	return nil
}

func _LoginResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoginResponse)
	switch tag {
	case 1: // MessageResponse.empty
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Empty)
		err := b.DecodeMessage(msg)
		m.MessageResponse = &LoginResponse_Empty{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoginResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoginResponse)
	// MessageResponse
	switch x := m.MessageResponse.(type) {
	case *LoginResponse_Empty:
		s := proto.Size(x.Empty)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type IpCount struct {
	Ip    string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Count uint64 `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
}

func (m *IpCount) Reset()                    { *m = IpCount{} }
func (m *IpCount) String() string            { return proto.CompactTextString(m) }
func (*IpCount) ProtoMessage()               {}
func (*IpCount) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *IpCount) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *IpCount) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*LoginRequest)(nil), "loginpb.LoginRequest")
	proto.RegisterType((*LoginResponse)(nil), "loginpb.LoginResponse")
	proto.RegisterType((*IpCount)(nil), "loginpb.IpCount")
	proto.RegisterType((*Empty)(nil), "loginpb.Empty")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for LoginService service

type LoginServiceClient interface {
	Serve(ctx context.Context, opts ...grpc.CallOption) (LoginService_ServeClient, error)
}

type loginServiceClient struct {
	cc *grpc.ClientConn
}

func NewLoginServiceClient(cc *grpc.ClientConn) LoginServiceClient {
	return &loginServiceClient{cc}
}

func (c *loginServiceClient) Serve(ctx context.Context, opts ...grpc.CallOption) (LoginService_ServeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_LoginService_serviceDesc.Streams[0], c.cc, "/loginpb.LoginService/Serve", opts...)
	if err != nil {
		return nil, err
	}
	x := &loginServiceServeClient{stream}
	return x, nil
}

type LoginService_ServeClient interface {
	Send(*LoginRequest) error
	Recv() (*LoginResponse, error)
	grpc.ClientStream
}

type loginServiceServeClient struct {
	grpc.ClientStream
}

func (x *loginServiceServeClient) Send(m *LoginRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *loginServiceServeClient) Recv() (*LoginResponse, error) {
	m := new(LoginResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for LoginService service

type LoginServiceServer interface {
	Serve(LoginService_ServeServer) error
}

func RegisterLoginServiceServer(s *grpc.Server, srv LoginServiceServer) {
	s.RegisterService(&_LoginService_serviceDesc, srv)
}

func _LoginService_Serve_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LoginServiceServer).Serve(&loginServiceServeServer{stream})
}

type LoginService_ServeServer interface {
	Send(*LoginResponse) error
	Recv() (*LoginRequest, error)
	grpc.ServerStream
}

type loginServiceServeServer struct {
	grpc.ServerStream
}

func (x *loginServiceServeServer) Send(m *LoginResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *loginServiceServeServer) Recv() (*LoginRequest, error) {
	m := new(LoginRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _LoginService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "loginpb.LoginService",
	HandlerType: (*LoginServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Serve",
			Handler:       _LoginService_Serve_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "login.proto",
}

func init() { proto.RegisterFile("login.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xc9, 0x4f, 0xcf,
	0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x07, 0x73, 0x0a, 0x92, 0x94, 0xfc, 0xb8,
	0x78, 0x7c, 0x40, 0xcc, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x1d, 0x2e, 0xf6, 0xcc,
	0x02, 0xe7, 0xfc, 0xd2, 0xbc, 0x12, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x6e, 0x23, 0x01, 0x3d, 0xa8,
	0x52, 0x3d, 0x4f, 0x88, 0xb8, 0x07, 0x43, 0x10, 0x4c, 0x89, 0x93, 0x00, 0x17, 0x9f, 0x6f, 0x6a,
	0x71, 0x71, 0x62, 0x7a, 0x2a, 0x54, 0xbf, 0x92, 0x17, 0x17, 0x2f, 0xd4, 0xbc, 0xe2, 0x82, 0xfc,
	0xbc, 0xe2, 0x54, 0x21, 0x35, 0x2e, 0xd6, 0xd4, 0xdc, 0x82, 0x92, 0x4a, 0xa8, 0x71, 0x7c, 0x70,
	0xe3, 0x5c, 0x41, 0xa2, 0x1e, 0x0c, 0x41, 0x10, 0x69, 0x27, 0x41, 0x2e, 0x7e, 0xb8, 0x51, 0x10,
	0xad, 0x4a, 0xfa, 0x5c, 0xec, 0x50, 0x3b, 0x85, 0xf8, 0xb8, 0x98, 0x32, 0x0b, 0xc0, 0x46, 0x70,
	0x06, 0x31, 0x65, 0x16, 0x08, 0x89, 0x70, 0xb1, 0x26, 0x83, 0x1d, 0xc9, 0xa4, 0xc0, 0xa8, 0xc1,
	0x12, 0x04, 0xe1, 0x28, 0xb1, 0x73, 0xb1, 0x82, 0x4d, 0x35, 0xf2, 0x81, 0xfa, 0x2a, 0x38, 0xb5,
	0xa8, 0x2c, 0x33, 0x39, 0x55, 0xc8, 0x86, 0x8b, 0x15, 0xc4, 0x4c, 0x15, 0x12, 0x85, 0x5b, 0x8f,
	0xec, 0x6b, 0x29, 0x31, 0x74, 0x61, 0xa8, 0x0b, 0x18, 0x34, 0x18, 0x0d, 0x18, 0x93, 0xd8, 0xc0,
	0x61, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x41, 0xb7, 0x99, 0x98, 0x42, 0x01, 0x00, 0x00,
}