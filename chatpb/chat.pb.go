// Code generated by protoc-gen-go.
// source: chat.proto
// DO NOT EDIT!

/*
Package chatpb is a generated protocol buffer package.

It is generated from these files:
	chat.proto

It has these top-level messages:
	ChatRequest
	ChatResponse
*/
package chatpb

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

type ChatRequest struct {
	// Types that are valid to be assigned to MessageRequest:
	//	*ChatRequest_SubRequest
	//	*ChatRequest_UnsubRequest
	//	*ChatRequest_PubRequest
	MessageRequest isChatRequest_MessageRequest `protobuf_oneof:"MessageRequest"`
}

func (m *ChatRequest) Reset()                    { *m = ChatRequest{} }
func (m *ChatRequest) String() string            { return proto.CompactTextString(m) }
func (*ChatRequest) ProtoMessage()               {}
func (*ChatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isChatRequest_MessageRequest interface {
	isChatRequest_MessageRequest()
}

type ChatRequest_SubRequest struct {
	SubRequest *gatewaypb.SubRequest `protobuf:"bytes,1,opt,name=subRequest,oneof"`
}
type ChatRequest_UnsubRequest struct {
	UnsubRequest *gatewaypb.UnsubRequest `protobuf:"bytes,2,opt,name=unsubRequest,oneof"`
}
type ChatRequest_PubRequest struct {
	PubRequest *gatewaypb.PubRequest `protobuf:"bytes,3,opt,name=pubRequest,oneof"`
}

func (*ChatRequest_SubRequest) isChatRequest_MessageRequest()   {}
func (*ChatRequest_UnsubRequest) isChatRequest_MessageRequest() {}
func (*ChatRequest_PubRequest) isChatRequest_MessageRequest()   {}

func (m *ChatRequest) GetMessageRequest() isChatRequest_MessageRequest {
	if m != nil {
		return m.MessageRequest
	}
	return nil
}

func (m *ChatRequest) GetSubRequest() *gatewaypb.SubRequest {
	if x, ok := m.GetMessageRequest().(*ChatRequest_SubRequest); ok {
		return x.SubRequest
	}
	return nil
}

func (m *ChatRequest) GetUnsubRequest() *gatewaypb.UnsubRequest {
	if x, ok := m.GetMessageRequest().(*ChatRequest_UnsubRequest); ok {
		return x.UnsubRequest
	}
	return nil
}

func (m *ChatRequest) GetPubRequest() *gatewaypb.PubRequest {
	if x, ok := m.GetMessageRequest().(*ChatRequest_PubRequest); ok {
		return x.PubRequest
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ChatRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ChatRequest_OneofMarshaler, _ChatRequest_OneofUnmarshaler, _ChatRequest_OneofSizer, []interface{}{
		(*ChatRequest_SubRequest)(nil),
		(*ChatRequest_UnsubRequest)(nil),
		(*ChatRequest_PubRequest)(nil),
	}
}

func _ChatRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ChatRequest)
	// MessageRequest
	switch x := m.MessageRequest.(type) {
	case *ChatRequest_SubRequest:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SubRequest); err != nil {
			return err
		}
	case *ChatRequest_UnsubRequest:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.UnsubRequest); err != nil {
			return err
		}
	case *ChatRequest_PubRequest:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PubRequest); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ChatRequest.MessageRequest has unexpected type %T", x)
	}
	return nil
}

func _ChatRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ChatRequest)
	switch tag {
	case 1: // MessageRequest.subRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(gatewaypb.SubRequest)
		err := b.DecodeMessage(msg)
		m.MessageRequest = &ChatRequest_SubRequest{msg}
		return true, err
	case 2: // MessageRequest.unsubRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(gatewaypb.UnsubRequest)
		err := b.DecodeMessage(msg)
		m.MessageRequest = &ChatRequest_UnsubRequest{msg}
		return true, err
	case 3: // MessageRequest.pubRequest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(gatewaypb.PubRequest)
		err := b.DecodeMessage(msg)
		m.MessageRequest = &ChatRequest_PubRequest{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ChatRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ChatRequest)
	// MessageRequest
	switch x := m.MessageRequest.(type) {
	case *ChatRequest_SubRequest:
		s := proto.Size(x.SubRequest)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ChatRequest_UnsubRequest:
		s := proto.Size(x.UnsubRequest)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ChatRequest_PubRequest:
		s := proto.Size(x.PubRequest)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ChatResponse struct {
	// Types that are valid to be assigned to MessageResponse:
	//	*ChatResponse_PubResponse
	MessageResponse isChatResponse_MessageResponse `protobuf_oneof:"MessageResponse"`
}

func (m *ChatResponse) Reset()                    { *m = ChatResponse{} }
func (m *ChatResponse) String() string            { return proto.CompactTextString(m) }
func (*ChatResponse) ProtoMessage()               {}
func (*ChatResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isChatResponse_MessageResponse interface {
	isChatResponse_MessageResponse()
}

type ChatResponse_PubResponse struct {
	PubResponse *gatewaypb.PubResponse `protobuf:"bytes,1,opt,name=pubResponse,oneof"`
}

func (*ChatResponse_PubResponse) isChatResponse_MessageResponse() {}

func (m *ChatResponse) GetMessageResponse() isChatResponse_MessageResponse {
	if m != nil {
		return m.MessageResponse
	}
	return nil
}

func (m *ChatResponse) GetPubResponse() *gatewaypb.PubResponse {
	if x, ok := m.GetMessageResponse().(*ChatResponse_PubResponse); ok {
		return x.PubResponse
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ChatResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ChatResponse_OneofMarshaler, _ChatResponse_OneofUnmarshaler, _ChatResponse_OneofSizer, []interface{}{
		(*ChatResponse_PubResponse)(nil),
	}
}

func _ChatResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ChatResponse)
	// MessageResponse
	switch x := m.MessageResponse.(type) {
	case *ChatResponse_PubResponse:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PubResponse); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ChatResponse.MessageResponse has unexpected type %T", x)
	}
	return nil
}

func _ChatResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ChatResponse)
	switch tag {
	case 1: // MessageResponse.pubResponse
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(gatewaypb.PubResponse)
		err := b.DecodeMessage(msg)
		m.MessageResponse = &ChatResponse_PubResponse{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ChatResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ChatResponse)
	// MessageResponse
	switch x := m.MessageResponse.(type) {
	case *ChatResponse_PubResponse:
		s := proto.Size(x.PubResponse)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ChatRequest)(nil), "chatpb.ChatRequest")
	proto.RegisterType((*ChatResponse)(nil), "chatpb.ChatResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ChatService service

type ChatServiceClient interface {
	Serve(ctx context.Context, opts ...grpc.CallOption) (ChatService_ServeClient, error)
}

type chatServiceClient struct {
	cc *grpc.ClientConn
}

func NewChatServiceClient(cc *grpc.ClientConn) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) Serve(ctx context.Context, opts ...grpc.CallOption) (ChatService_ServeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_ChatService_serviceDesc.Streams[0], c.cc, "/chatpb.ChatService/Serve", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceServeClient{stream}
	return x, nil
}

type ChatService_ServeClient interface {
	Send(*ChatRequest) error
	Recv() (*ChatResponse, error)
	grpc.ClientStream
}

type chatServiceServeClient struct {
	grpc.ClientStream
}

func (x *chatServiceServeClient) Send(m *ChatRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceServeClient) Recv() (*ChatResponse, error) {
	m := new(ChatResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for ChatService service

type ChatServiceServer interface {
	Serve(ChatService_ServeServer) error
}

func RegisterChatServiceServer(s *grpc.Server, srv ChatServiceServer) {
	s.RegisterService(&_ChatService_serviceDesc, srv)
}

func _ChatService_Serve_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Serve(&chatServiceServeServer{stream})
}

type ChatService_ServeServer interface {
	Send(*ChatResponse) error
	Recv() (*ChatRequest, error)
	grpc.ServerStream
}

type chatServiceServeServer struct {
	grpc.ServerStream
}

func (x *chatServiceServeServer) Send(m *ChatResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceServeServer) Recv() (*ChatRequest, error) {
	m := new(ChatRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _ChatService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chatpb.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Serve",
			Handler:       _ChatService_Serve_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chat.proto",
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xce, 0x48, 0x2c,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0xb1, 0x0b, 0x92, 0xa4, 0xe4, 0x92, 0x4b,
	0x2a, 0x8a, 0xcb, 0x33, 0xf4, 0xd3, 0x13, 0x4b, 0x52, 0xcb, 0x13, 0x2b, 0x0b, 0x92, 0x60, 0x2c,
	0x88, 0x3a, 0xa5, 0x0b, 0x8c, 0x5c, 0xdc, 0xce, 0x19, 0x89, 0x25, 0x41, 0xa9, 0x85, 0xa5, 0xa9,
	0xc5, 0x25, 0x42, 0xe6, 0x5c, 0x5c, 0xc5, 0xa5, 0x49, 0x50, 0x9e, 0x04, 0xa3, 0x02, 0xa3, 0x06,
	0xb7, 0x91, 0xa8, 0x1e, 0x5c, 0xb7, 0x5e, 0x30, 0x5c, 0xd2, 0x83, 0x21, 0x08, 0x49, 0xa9, 0x90,
	0x2d, 0x17, 0x4f, 0x69, 0x1e, 0x92, 0x56, 0x26, 0xb0, 0x56, 0x71, 0x24, 0xad, 0xa1, 0x48, 0xd2,
	0x1e, 0x0c, 0x41, 0x28, 0xca, 0x41, 0xf6, 0x16, 0x20, 0x34, 0x33, 0x63, 0xd8, 0x1b, 0x80, 0x62,
	0x2f, 0x42, 0xa9, 0x93, 0x00, 0x17, 0x9f, 0x6f, 0x6a, 0x71, 0x71, 0x62, 0x7a, 0x2a, 0x54, 0x44,
	0x29, 0x96, 0x8b, 0x07, 0xe2, 0xa3, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x2b, 0x2e, 0x6e,
	0xb0, 0x7a, 0x08, 0x17, 0xea, 0x27, 0x31, 0x74, 0xb3, 0x21, 0xb2, 0x1e, 0x0c, 0x41, 0xc8, 0x8a,
	0x9d, 0x04, 0xb9, 0xf8, 0xe1, 0xa6, 0x43, 0x84, 0x8c, 0xdc, 0x21, 0x01, 0x16, 0x9c, 0x5a, 0x54,
	0x96, 0x99, 0x9c, 0x2a, 0x64, 0xc1, 0xc5, 0x0a, 0x62, 0xa6, 0x0a, 0x09, 0xeb, 0x41, 0x82, 0x5c,
	0x0f, 0x29, 0x38, 0xa5, 0x44, 0x50, 0x05, 0x21, 0x46, 0x28, 0x31, 0x68, 0x30, 0x1a, 0x30, 0x26,
	0xb1, 0x81, 0x63, 0xc0, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x6b, 0x28, 0x92, 0x97, 0xb7, 0x01,
	0x00, 0x00,
}