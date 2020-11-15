// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: oneof/oneof.proto

package oneof_pb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthorId int32 `protobuf:"varint,1,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	// Types that are assignable to Content:
	//	*Entry_Text
	//	*Entry_Image
	//	*Entry_Link
	Content isEntry_Content `protobuf_oneof:"content"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oneof_oneof_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_oneof_oneof_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_oneof_oneof_proto_rawDescGZIP(), []int{0}
}

func (x *Entry) GetAuthorId() int32 {
	if x != nil {
		return x.AuthorId
	}
	return 0
}

func (m *Entry) GetContent() isEntry_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *Entry) GetText() *Text {
	if x, ok := x.GetContent().(*Entry_Text); ok {
		return x.Text
	}
	return nil
}

func (x *Entry) GetImage() *Image {
	if x, ok := x.GetContent().(*Entry_Image); ok {
		return x.Image
	}
	return nil
}

func (x *Entry) GetLink() *Link {
	if x, ok := x.GetContent().(*Entry_Link); ok {
		return x.Link
	}
	return nil
}

type isEntry_Content interface {
	isEntry_Content()
}

type Entry_Text struct {
	Text *Text `protobuf:"bytes,11,opt,name=text,proto3,oneof"`
}

type Entry_Image struct {
	Image *Image `protobuf:"bytes,12,opt,name=image,proto3,oneof"`
}

type Entry_Link struct {
	Link *Link `protobuf:"bytes,13,opt,name=link,proto3,oneof"`
}

func (*Entry_Text) isEntry_Content() {}

func (*Entry_Image) isEntry_Content() {}

func (*Entry_Link) isEntry_Content() {}

type Text struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Body string `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Text) Reset() {
	*x = Text{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oneof_oneof_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Text) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Text) ProtoMessage() {}

func (x *Text) ProtoReflect() protoreflect.Message {
	mi := &file_oneof_oneof_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Text.ProtoReflect.Descriptor instead.
func (*Text) Descriptor() ([]byte, []int) {
	return file_oneof_oneof_proto_rawDescGZIP(), []int{1}
}

func (x *Text) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Text) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type Image struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url    string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Width  uint32 `protobuf:"varint,3,opt,name=width,proto3" json:"width,omitempty"`
	Height uint32 `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *Image) Reset() {
	*x = Image{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oneof_oneof_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Image) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Image) ProtoMessage() {}

func (x *Image) ProtoReflect() protoreflect.Message {
	mi := &file_oneof_oneof_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Image.ProtoReflect.Descriptor instead.
func (*Image) Descriptor() ([]byte, []int) {
	return file_oneof_oneof_proto_rawDescGZIP(), []int{2}
}

func (x *Image) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Image) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Image) GetWidth() uint32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *Image) GetHeight() uint32 {
	if x != nil {
		return x.Height
	}
	return 0
}

type Link struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Link) Reset() {
	*x = Link{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oneof_oneof_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Link) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Link) ProtoMessage() {}

func (x *Link) ProtoReflect() protoreflect.Message {
	mi := &file_oneof_oneof_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Link.ProtoReflect.Descriptor instead.
func (*Link) Descriptor() ([]byte, []int) {
	return file_oneof_oneof_proto_rawDescGZIP(), []int{3}
}

func (x *Link) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Link) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_oneof_oneof_proto protoreflect.FileDescriptor

var file_oneof_oneof_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x2f, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x6f, 0x6e, 0x65,
	0x6f, 0x66, 0x22, 0xb3, 0x01, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70,
	0x69, 0x2e, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x2e, 0x54, 0x65, 0x78, 0x74, 0x48, 0x00, 0x52, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x12, 0x2c, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x6f, 0x6e,
	0x65, 0x6f, 0x66, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x05, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x6f, 0x6e, 0x65, 0x6f, 0x66,
	0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x6b, 0x42, 0x09, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x2a, 0x0a, 0x04, 0x54, 0x65, 0x78, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x22, 0x57, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x77, 0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x28, 0x0a,
	0x04, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x42, 0x18, 0x5a, 0x16, 0x61, 0x70, 0x69, 0x73, 0x2f,
	0x67, 0x6f, 0x2f, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x3b, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x5f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_oneof_oneof_proto_rawDescOnce sync.Once
	file_oneof_oneof_proto_rawDescData = file_oneof_oneof_proto_rawDesc
)

func file_oneof_oneof_proto_rawDescGZIP() []byte {
	file_oneof_oneof_proto_rawDescOnce.Do(func() {
		file_oneof_oneof_proto_rawDescData = protoimpl.X.CompressGZIP(file_oneof_oneof_proto_rawDescData)
	})
	return file_oneof_oneof_proto_rawDescData
}

var file_oneof_oneof_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_oneof_oneof_proto_goTypes = []interface{}{
	(*Entry)(nil), // 0: testapi.oneof.Entry
	(*Text)(nil),  // 1: testapi.oneof.Text
	(*Image)(nil), // 2: testapi.oneof.Image
	(*Link)(nil),  // 3: testapi.oneof.Link
}
var file_oneof_oneof_proto_depIdxs = []int32{
	1, // 0: testapi.oneof.Entry.text:type_name -> testapi.oneof.Text
	2, // 1: testapi.oneof.Entry.image:type_name -> testapi.oneof.Image
	3, // 2: testapi.oneof.Entry.link:type_name -> testapi.oneof.Link
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_oneof_oneof_proto_init() }
func file_oneof_oneof_proto_init() {
	if File_oneof_oneof_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_oneof_oneof_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_oneof_oneof_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Text); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_oneof_oneof_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Image); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_oneof_oneof_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Link); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_oneof_oneof_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Entry_Text)(nil),
		(*Entry_Image)(nil),
		(*Entry_Link)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_oneof_oneof_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_oneof_oneof_proto_goTypes,
		DependencyIndexes: file_oneof_oneof_proto_depIdxs,
		MessageInfos:      file_oneof_oneof_proto_msgTypes,
	}.Build()
	File_oneof_oneof_proto = out.File
	file_oneof_oneof_proto_rawDesc = nil
	file_oneof_oneof_proto_goTypes = nil
	file_oneof_oneof_proto_depIdxs = nil
}
