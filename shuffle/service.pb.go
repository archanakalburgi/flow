// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: shuffle/service.proto

package shuffle

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	protocol "go.gazette.dev/core/broker/protocol"
	go_gazette_dev_core_message "go.gazette.dev/core/message"
	grpc "google.golang.org/grpc"
	io "io"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Config struct {
	Processors []Config_Processor `protobuf:"bytes,1,rep,name=processors,proto3" json:"processors"`
	// JSON-Pointer to the message UUID.
	UuidJsonPtr string `protobuf:"bytes,2,opt,name=uuid_json_ptr,json=uuidJsonPtr,proto3" json:"uuid_json_ptr,omitempty"`
	// Composite key over which shuffling occurs, specified as one or more
	// JSON-Pointers indicating a message location to extract.
	ShuffleKeyPtr []string `protobuf:"bytes,3,rep,name=shuffle_key_ptr,json=shuffleKeyPtr,proto3" json:"shuffle_key_ptr,omitempty"`
	// Number of top-ranked processors from which a single processor will be
	// randomly selected, after shuffling. This final selection is deterministic,
	// being derived from the message clock time. Usually this is one. Larger
	// values can help distribute the effect of "hot keys" which would otherwise
	// overwhelm specific processors. If non-zero, |broadcast_to| cannot be set.
	ChooseFrom uint32 `protobuf:"varint,5,opt,name=choose_from,json=chooseFrom,proto3" json:"choose_from,omitempty"`
	// Number of top-ranked processors to broadcast each message to, after shuffling.
	// Usually this is zero. If non-zero, |choose_from| cannot be set.
	BroadcastTo uint32 `protobuf:"varint,4,opt,name=broadcast_to,json=broadcastTo,proto3" json:"broadcast_to,omitempty"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_6033d9fdececa05a, []int{0}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

type Config_Processor struct {
	MinMsgClock go_gazette_dev_core_message.Clock `protobuf:"varint,1,opt,name=min_msg_clock,json=minMsgClock,proto3,casttype=go.gazette.dev/core/message.Clock" json:"min_msg_clock,omitempty"`
	MaxMsgClock go_gazette_dev_core_message.Clock `protobuf:"varint,2,opt,name=max_msg_clock,json=maxMsgClock,proto3,casttype=go.gazette.dev/core/message.Clock" json:"max_msg_clock,omitempty"`
}

func (m *Config_Processor) Reset()         { *m = Config_Processor{} }
func (m *Config_Processor) String() string { return proto.CompactTextString(m) }
func (*Config_Processor) ProtoMessage()    {}
func (*Config_Processor) Descriptor() ([]byte, []int) {
	return fileDescriptor_6033d9fdececa05a, []int{0, 0}
}
func (m *Config_Processor) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Config_Processor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Config_Processor.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Config_Processor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config_Processor.Merge(m, src)
}
func (m *Config_Processor) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Config_Processor) XXX_DiscardUnknown() {
	xxx_messageInfo_Config_Processor.DiscardUnknown(m)
}

var xxx_messageInfo_Config_Processor proto.InternalMessageInfo

type Request struct {
	// Configuration under which shuffling is to occur.
	Config Config `protobuf:"bytes,1,opt,name=config,proto3" json:"config"`
	// Index of this reader within the shuffle topology.
	ShuffleIndex int32 `protobuf:"varint,2,opt,name=shuffle_index,json=shuffleIndex,proto3" json:"shuffle_index,omitempty"`
	// Nominal ReadRequest whose response is to be shuffled.
	protocol.ReadRequest `protobuf:"bytes,3,opt,name=read_request,json=readRequest,proto3,embedded=read_request" json:"read_request"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_6033d9fdececa05a, []int{1}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return m.ProtoSize()
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Config)(nil), "shuffle.Config")
	proto.RegisterType((*Config_Processor)(nil), "shuffle.Config.Processor")
	proto.RegisterType((*Request)(nil), "shuffle.Request")
}

func init() { proto.RegisterFile("shuffle/service.proto", fileDescriptor_6033d9fdececa05a) }

var fileDescriptor_6033d9fdececa05a = []byte{
	// 485 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xc7, 0x7d, 0xb5, 0xdb, 0x90, 0x73, 0xa2, 0xa2, 0x13, 0x45, 0x26, 0x83, 0xed, 0x06, 0x15,
	0x79, 0xc1, 0x46, 0x61, 0x43, 0x42, 0x48, 0xae, 0x84, 0x54, 0x10, 0x52, 0x65, 0x98, 0x58, 0x2c,
	0xc7, 0x7e, 0x71, 0x4d, 0x62, 0xbf, 0x70, 0xe7, 0x54, 0x29, 0xdf, 0x80, 0x8d, 0x8f, 0xd0, 0xad,
	0x5f, 0x25, 0x63, 0x46, 0xa6, 0x48, 0x34, 0x0b, 0x9f, 0x81, 0x09, 0xf9, 0xec, 0x98, 0xaa, 0x1b,
	0x8b, 0x75, 0xf7, 0x7f, 0xff, 0xf7, 0xf3, 0x5f, 0xef, 0x1e, 0x3d, 0x12, 0x17, 0x8b, 0xc9, 0x64,
	0x06, 0x9e, 0x00, 0x7e, 0x99, 0xc5, 0xe0, 0xce, 0x39, 0x96, 0xc8, 0x3a, 0x8d, 0x3c, 0x78, 0x94,
	0x62, 0x8a, 0x52, 0xf3, 0xaa, 0x53, 0x5d, 0x1e, 0x98, 0x63, 0x8e, 0x53, 0xe0, 0x9e, 0xbc, 0xc5,
	0x38, 0x6b, 0x0f, 0x75, 0x7d, 0xf8, 0x5d, 0xa5, 0x07, 0xa7, 0x58, 0x4c, 0xb2, 0x94, 0xbd, 0xa1,
	0x74, 0xce, 0x31, 0x06, 0x21, 0x90, 0x0b, 0x83, 0xd8, 0xaa, 0xa3, 0x8f, 0x9e, 0xb8, 0x0d, 0xde,
	0xad, 0x4d, 0xee, 0xf9, 0xce, 0xe1, 0x6b, 0xab, 0x8d, 0xa5, 0x04, 0x77, 0x5a, 0xd8, 0x90, 0xf6,
	0x17, 0x8b, 0x2c, 0x09, 0xbf, 0x08, 0x2c, 0xc2, 0x79, 0xc9, 0x8d, 0x3d, 0x9b, 0x38, 0xdd, 0x40,
	0xaf, 0xc4, 0x77, 0x02, 0x8b, 0xf3, 0x92, 0xb3, 0x67, 0xf4, 0xb0, 0x21, 0x86, 0x53, 0xb8, 0x92,
	0x2e, 0xd5, 0x56, 0x9d, 0x6e, 0xd0, 0x6f, 0xe4, 0xf7, 0x70, 0x55, 0xf9, 0x2c, 0xaa, 0xc7, 0x17,
	0x88, 0x02, 0xc2, 0x09, 0xc7, 0xdc, 0xd8, 0xb7, 0x89, 0xd3, 0x0f, 0x68, 0x2d, 0xbd, 0xe5, 0x98,
	0xb3, 0x63, 0xda, 0x1b, 0x73, 0x8c, 0x92, 0x38, 0x12, 0x65, 0x58, 0xa2, 0xa1, 0x49, 0x87, 0xde,
	0x6a, 0x9f, 0x70, 0x70, 0x43, 0x68, 0xb7, 0xcd, 0xcb, 0xce, 0x68, 0x3f, 0xcf, 0x8a, 0x30, 0x17,
	0x69, 0x18, 0xcf, 0x30, 0x9e, 0x1a, 0xc4, 0x26, 0x8e, 0xe6, 0x9f, 0xfc, 0xd9, 0x58, 0xc7, 0x29,
	0xba, 0x69, 0xf4, 0x0d, 0xca, 0x12, 0xdc, 0x04, 0x2e, 0xbd, 0x18, 0x39, 0x78, 0x39, 0x08, 0x11,
	0xa5, 0xe0, 0x9e, 0x56, 0xe6, 0x40, 0xcf, 0xb3, 0xe2, 0x83, 0x48, 0xe5, 0x45, 0xa2, 0xa2, 0xe5,
	0x1d, 0xd4, 0xde, 0xff, 0xa1, 0xa2, 0xe5, 0x0e, 0xf5, 0x4a, 0xfb, 0x7d, 0x6d, 0x91, 0xfa, 0x3b,
	0xbc, 0x21, 0xb4, 0x13, 0xc0, 0xd7, 0x05, 0x88, 0x92, 0x3d, 0xa7, 0x07, 0xb1, 0x9c, 0xb8, 0x8c,
	0xa9, 0x8f, 0x0e, 0xef, 0x3d, 0x44, 0x33, 0xfe, 0xc6, 0xc4, 0x9e, 0xd2, 0xdd, 0xfc, 0xc2, 0xac,
	0x48, 0x60, 0x29, 0x13, 0xed, 0x07, 0xbd, 0x46, 0x3c, 0xab, 0x34, 0xe6, 0xd3, 0x1e, 0x87, 0x28,
	0x09, 0x79, 0xfd, 0x0f, 0x43, 0x95, 0xe4, 0x23, 0xb7, 0x5d, 0x89, 0x00, 0xa2, 0xa4, 0x09, 0xe0,
	0x3f, 0xa8, 0xf8, 0xeb, 0x8d, 0x45, 0x02, 0x9d, 0xff, 0x93, 0xeb, 0xa4, 0xa3, 0xd7, 0xb4, 0xf3,
	0xb1, 0x26, 0xb3, 0x11, 0xd5, 0xaa, 0x36, 0xf6, 0xb0, 0x0d, 0xd8, 0x58, 0x07, 0x8f, 0xef, 0x83,
	0xc5, 0x1c, 0x0b, 0x01, 0x2f, 0x88, 0x7f, 0xb2, 0xfa, 0x65, 0x2a, 0xab, 0x5b, 0x93, 0xac, 0x6f,
	0x4d, 0xf2, 0x63, 0x6b, 0x2a, 0xd7, 0x5b, 0x93, 0xac, 0xb7, 0xa6, 0xf2, 0x73, 0x6b, 0x2a, 0x9f,
	0x77, 0x1b, 0x3d, 0x3e, 0x90, 0xfd, 0x2f, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x99, 0x88, 0xee,
	0x8e, 0xfa, 0x02, 0x00, 0x00,
}

func (this *Config) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Config)
	if !ok {
		that2, ok := that.(Config)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Processors) != len(that1.Processors) {
		return false
	}
	for i := range this.Processors {
		if !this.Processors[i].Equal(&that1.Processors[i]) {
			return false
		}
	}
	if this.UuidJsonPtr != that1.UuidJsonPtr {
		return false
	}
	if len(this.ShuffleKeyPtr) != len(that1.ShuffleKeyPtr) {
		return false
	}
	for i := range this.ShuffleKeyPtr {
		if this.ShuffleKeyPtr[i] != that1.ShuffleKeyPtr[i] {
			return false
		}
	}
	if this.ChooseFrom != that1.ChooseFrom {
		return false
	}
	if this.BroadcastTo != that1.BroadcastTo {
		return false
	}
	return true
}
func (this *Config_Processor) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Config_Processor)
	if !ok {
		that2, ok := that.(Config_Processor)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.MinMsgClock != that1.MinMsgClock {
		return false
	}
	if this.MaxMsgClock != that1.MaxMsgClock {
		return false
	}
	return true
}
func (this *Request) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Request)
	if !ok {
		that2, ok := that.(Request)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Config.Equal(&that1.Config) {
		return false
	}
	if this.ShuffleIndex != that1.ShuffleIndex {
		return false
	}
	if !this.ReadRequest.Equal(&that1.ReadRequest) {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ShuffleClient is the client API for Shuffle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ShuffleClient interface {
	// Subscribe a Shard to read a Journal.
	Read(ctx context.Context, in *Request, opts ...grpc.CallOption) (Shuffle_ReadClient, error)
}

type shuffleClient struct {
	cc *grpc.ClientConn
}

func NewShuffleClient(cc *grpc.ClientConn) ShuffleClient {
	return &shuffleClient{cc}
}

func (c *shuffleClient) Read(ctx context.Context, in *Request, opts ...grpc.CallOption) (Shuffle_ReadClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Shuffle_serviceDesc.Streams[0], "/shuffle.Shuffle/Read", opts...)
	if err != nil {
		return nil, err
	}
	x := &shuffleReadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Shuffle_ReadClient interface {
	Recv() (*protocol.ReadResponse, error)
	grpc.ClientStream
}

type shuffleReadClient struct {
	grpc.ClientStream
}

func (x *shuffleReadClient) Recv() (*protocol.ReadResponse, error) {
	m := new(protocol.ReadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ShuffleServer is the server API for Shuffle service.
type ShuffleServer interface {
	// Subscribe a Shard to read a Journal.
	Read(*Request, Shuffle_ReadServer) error
}

func RegisterShuffleServer(s *grpc.Server, srv ShuffleServer) {
	s.RegisterService(&_Shuffle_serviceDesc, srv)
}

func _Shuffle_Read_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ShuffleServer).Read(m, &shuffleReadServer{stream})
}

type Shuffle_ReadServer interface {
	Send(*protocol.ReadResponse) error
	grpc.ServerStream
}

type shuffleReadServer struct {
	grpc.ServerStream
}

func (x *shuffleReadServer) Send(m *protocol.ReadResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Shuffle_serviceDesc = grpc.ServiceDesc{
	ServiceName: "shuffle.Shuffle",
	HandlerType: (*ShuffleServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Read",
			Handler:       _Shuffle_Read_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "shuffle/service.proto",
}

func (m *Config) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Config) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Processors) > 0 {
		for _, msg := range m.Processors {
			dAtA[i] = 0xa
			i++
			i = encodeVarintService(dAtA, i, uint64(msg.ProtoSize()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.UuidJsonPtr) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintService(dAtA, i, uint64(len(m.UuidJsonPtr)))
		i += copy(dAtA[i:], m.UuidJsonPtr)
	}
	if len(m.ShuffleKeyPtr) > 0 {
		for _, s := range m.ShuffleKeyPtr {
			dAtA[i] = 0x1a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if m.BroadcastTo != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintService(dAtA, i, uint64(m.BroadcastTo))
	}
	if m.ChooseFrom != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintService(dAtA, i, uint64(m.ChooseFrom))
	}
	return i, nil
}

func (m *Config_Processor) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Config_Processor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.MinMsgClock != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintService(dAtA, i, uint64(m.MinMsgClock))
	}
	if m.MaxMsgClock != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintService(dAtA, i, uint64(m.MaxMsgClock))
	}
	return i, nil
}

func (m *Request) Marshal() (dAtA []byte, err error) {
	size := m.ProtoSize()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Request) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintService(dAtA, i, uint64(m.Config.ProtoSize()))
	n1, err := m.Config.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if m.ShuffleIndex != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintService(dAtA, i, uint64(m.ShuffleIndex))
	}
	dAtA[i] = 0x1a
	i++
	i = encodeVarintService(dAtA, i, uint64(m.ReadRequest.ProtoSize()))
	n2, err := m.ReadRequest.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	return i, nil
}

func encodeVarintService(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Config) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Processors) > 0 {
		for _, e := range m.Processors {
			l = e.ProtoSize()
			n += 1 + l + sovService(uint64(l))
		}
	}
	l = len(m.UuidJsonPtr)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	if len(m.ShuffleKeyPtr) > 0 {
		for _, s := range m.ShuffleKeyPtr {
			l = len(s)
			n += 1 + l + sovService(uint64(l))
		}
	}
	if m.BroadcastTo != 0 {
		n += 1 + sovService(uint64(m.BroadcastTo))
	}
	if m.ChooseFrom != 0 {
		n += 1 + sovService(uint64(m.ChooseFrom))
	}
	return n
}

func (m *Config_Processor) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MinMsgClock != 0 {
		n += 1 + sovService(uint64(m.MinMsgClock))
	}
	if m.MaxMsgClock != 0 {
		n += 1 + sovService(uint64(m.MaxMsgClock))
	}
	return n
}

func (m *Request) ProtoSize() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Config.ProtoSize()
	n += 1 + l + sovService(uint64(l))
	if m.ShuffleIndex != 0 {
		n += 1 + sovService(uint64(m.ShuffleIndex))
	}
	l = m.ReadRequest.ProtoSize()
	n += 1 + l + sovService(uint64(l))
	return n
}

func sovService(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozService(x uint64) (n int) {
	return sovService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Processors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Processors = append(m.Processors, Config_Processor{})
			if err := m.Processors[len(m.Processors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UuidJsonPtr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UuidJsonPtr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShuffleKeyPtr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ShuffleKeyPtr = append(m.ShuffleKeyPtr, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BroadcastTo", wireType)
			}
			m.BroadcastTo = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BroadcastTo |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChooseFrom", wireType)
			}
			m.ChooseFrom = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChooseFrom |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Config_Processor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Processor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Processor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinMsgClock", wireType)
			}
			m.MinMsgClock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinMsgClock |= go_gazette_dev_core_message.Clock(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxMsgClock", wireType)
			}
			m.MaxMsgClock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxMsgClock |= go_gazette_dev_core_message.Clock(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShuffleIndex", wireType)
			}
			m.ShuffleIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ShuffleIndex |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReadRequest", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ReadRequest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthService
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthService
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowService
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipService(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthService
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthService = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowService   = fmt.Errorf("proto: integer overflow")
)