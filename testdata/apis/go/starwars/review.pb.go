// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: starwars/review.proto

package starwars_pb

import (
	proto "github.com/golang/protobuf/proto"
	_ "github.com/izumin5210/protogql/options"
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

// Represents a review for a movie
type Review struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The number of stars this review gave, 1-5
	Stars int32 `protobuf:"varint,1,opt,name=stars,proto3" json:"stars,omitempty"`
	// Comment about the movie
	Commentary string `protobuf:"bytes,2,opt,name=commentary,proto3" json:"commentary,omitempty"`
}

func (x *Review) Reset() {
	*x = Review{}
	if protoimpl.UnsafeEnabled {
		mi := &file_starwars_review_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Review) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Review) ProtoMessage() {}

func (x *Review) ProtoReflect() protoreflect.Message {
	mi := &file_starwars_review_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Review.ProtoReflect.Descriptor instead.
func (*Review) Descriptor() ([]byte, []int) {
	return file_starwars_review_proto_rawDescGZIP(), []int{0}
}

func (x *Review) GetStars() int32 {
	if x != nil {
		return x.Stars
	}
	return 0
}

func (x *Review) GetCommentary() string {
	if x != nil {
		return x.Commentary
	}
	return ""
}

type ListReviewsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Episode Episode `protobuf:"varint,1,opt,name=episode,proto3,enum=testapi.starwars.Episode" json:"episode,omitempty"`
}

func (x *ListReviewsRequest) Reset() {
	*x = ListReviewsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_starwars_review_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListReviewsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReviewsRequest) ProtoMessage() {}

func (x *ListReviewsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_starwars_review_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReviewsRequest.ProtoReflect.Descriptor instead.
func (*ListReviewsRequest) Descriptor() ([]byte, []int) {
	return file_starwars_review_proto_rawDescGZIP(), []int{1}
}

func (x *ListReviewsRequest) GetEpisode() Episode {
	if x != nil {
		return x.Episode
	}
	return Episode_EPISODE_UNSPECIFIED
}

type ListReviewsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reviews []*Review `protobuf:"bytes,1,rep,name=reviews,proto3" json:"reviews,omitempty"`
}

func (x *ListReviewsResponse) Reset() {
	*x = ListReviewsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_starwars_review_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListReviewsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReviewsResponse) ProtoMessage() {}

func (x *ListReviewsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_starwars_review_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReviewsResponse.ProtoReflect.Descriptor instead.
func (*ListReviewsResponse) Descriptor() ([]byte, []int) {
	return file_starwars_review_proto_rawDescGZIP(), []int{2}
}

func (x *ListReviewsResponse) GetReviews() []*Review {
	if x != nil {
		return x.Reviews
	}
	return nil
}

type CreateReviewRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Episode Episode `protobuf:"varint,1,opt,name=episode,proto3,enum=testapi.starwars.Episode" json:"episode,omitempty"`
	// The input object sent when someone is creating a new review
	Review *Review `protobuf:"bytes,2,opt,name=review,proto3" json:"review,omitempty"`
	// Favorite color, optional
	FavoriteColor *Color `protobuf:"bytes,3,opt,name=favorite_color,json=favoriteColor,proto3" json:"favorite_color,omitempty"`
}

func (x *CreateReviewRequest) Reset() {
	*x = CreateReviewRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_starwars_review_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateReviewRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateReviewRequest) ProtoMessage() {}

func (x *CreateReviewRequest) ProtoReflect() protoreflect.Message {
	mi := &file_starwars_review_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateReviewRequest.ProtoReflect.Descriptor instead.
func (*CreateReviewRequest) Descriptor() ([]byte, []int) {
	return file_starwars_review_proto_rawDescGZIP(), []int{3}
}

func (x *CreateReviewRequest) GetEpisode() Episode {
	if x != nil {
		return x.Episode
	}
	return Episode_EPISODE_UNSPECIFIED
}

func (x *CreateReviewRequest) GetReview() *Review {
	if x != nil {
		return x.Review
	}
	return nil
}

func (x *CreateReviewRequest) GetFavoriteColor() *Color {
	if x != nil {
		return x.FavoriteColor
	}
	return nil
}

// The input object sent when passing in a color
type Color struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Red   int32 `protobuf:"varint,1,opt,name=red,proto3" json:"red,omitempty"`
	Green int32 `protobuf:"varint,2,opt,name=green,proto3" json:"green,omitempty"`
	Blue  int32 `protobuf:"varint,3,opt,name=blue,proto3" json:"blue,omitempty"`
}

func (x *Color) Reset() {
	*x = Color{}
	if protoimpl.UnsafeEnabled {
		mi := &file_starwars_review_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Color) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Color) ProtoMessage() {}

func (x *Color) ProtoReflect() protoreflect.Message {
	mi := &file_starwars_review_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Color.ProtoReflect.Descriptor instead.
func (*Color) Descriptor() ([]byte, []int) {
	return file_starwars_review_proto_rawDescGZIP(), []int{4}
}

func (x *Color) GetRed() int32 {
	if x != nil {
		return x.Red
	}
	return 0
}

func (x *Color) GetGreen() int32 {
	if x != nil {
		return x.Green
	}
	return 0
}

func (x *Color) GetBlue() int32 {
	if x != nil {
		return x.Blue
	}
	return 0
}

var File_starwars_review_proto protoreflect.FileDescriptor

var file_starwars_review_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2f, 0x72, 0x65, 0x76, 0x69, 0x65,
	0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69,
	0x2e, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x1a, 0x15, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x13, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x06, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x61, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x61, 0x72, 0x79, 0x22, 0x49, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x76,
	0x69, 0x65, 0x77, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x65,
	0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e,
	0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52, 0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65,
	0x22, 0x49, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x07, 0x72, 0x65, 0x76, 0x69, 0x65,
	0x77, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x70, 0x69, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e, 0x52, 0x65, 0x76, 0x69,
	0x65, 0x77, 0x52, 0x07, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x73, 0x22, 0xbc, 0x01, 0x0a, 0x13,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e, 0x45, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x52,
	0x07, 0x65, 0x70, 0x69, 0x73, 0x6f, 0x64, 0x65, 0x12, 0x30, 0x0a, 0x06, 0x72, 0x65, 0x76, 0x69,
	0x65, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61,
	0x70, 0x69, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e, 0x52, 0x65, 0x76, 0x69,
	0x65, 0x77, 0x52, 0x06, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x3e, 0x0a, 0x0e, 0x66, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x61,
	0x72, 0x77, 0x61, 0x72, 0x73, 0x2e, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x52, 0x0d, 0x66, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x22, 0x43, 0x0a, 0x05, 0x43, 0x6f,
	0x6c, 0x6f, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x72, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x65, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x67, 0x72, 0x65, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x62, 0x6c, 0x75, 0x65, 0x32,
	0xe8, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x72, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x73,
	0x12, 0x24, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x77,
	0x61, 0x72, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69,
	0x2e, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x76, 0x69, 0x65, 0x77, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0xc2,
	0x80, 0x01, 0x12, 0x0a, 0x07, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x73, 0x1a, 0x07, 0x72, 0x65,
	0x76, 0x69, 0x65, 0x77, 0x73, 0x12, 0x63, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x25, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e,
	0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x2e,
	0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x22, 0x12, 0xca, 0x80, 0x01, 0x0e, 0x0a, 0x0c, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x42, 0x1e, 0x5a, 0x1c, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x73, 0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x3b, 0x73,
	0x74, 0x61, 0x72, 0x77, 0x61, 0x72, 0x73, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_starwars_review_proto_rawDescOnce sync.Once
	file_starwars_review_proto_rawDescData = file_starwars_review_proto_rawDesc
)

func file_starwars_review_proto_rawDescGZIP() []byte {
	file_starwars_review_proto_rawDescOnce.Do(func() {
		file_starwars_review_proto_rawDescData = protoimpl.X.CompressGZIP(file_starwars_review_proto_rawDescData)
	})
	return file_starwars_review_proto_rawDescData
}

var file_starwars_review_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_starwars_review_proto_goTypes = []interface{}{
	(*Review)(nil),              // 0: testapi.starwars.Review
	(*ListReviewsRequest)(nil),  // 1: testapi.starwars.ListReviewsRequest
	(*ListReviewsResponse)(nil), // 2: testapi.starwars.ListReviewsResponse
	(*CreateReviewRequest)(nil), // 3: testapi.starwars.CreateReviewRequest
	(*Color)(nil),               // 4: testapi.starwars.Color
	(Episode)(0),                // 5: testapi.starwars.Episode
}
var file_starwars_review_proto_depIdxs = []int32{
	5, // 0: testapi.starwars.ListReviewsRequest.episode:type_name -> testapi.starwars.Episode
	0, // 1: testapi.starwars.ListReviewsResponse.reviews:type_name -> testapi.starwars.Review
	5, // 2: testapi.starwars.CreateReviewRequest.episode:type_name -> testapi.starwars.Episode
	0, // 3: testapi.starwars.CreateReviewRequest.review:type_name -> testapi.starwars.Review
	4, // 4: testapi.starwars.CreateReviewRequest.favorite_color:type_name -> testapi.starwars.Color
	1, // 5: testapi.starwars.ReviewService.ListReviews:input_type -> testapi.starwars.ListReviewsRequest
	3, // 6: testapi.starwars.ReviewService.CreateReview:input_type -> testapi.starwars.CreateReviewRequest
	2, // 7: testapi.starwars.ReviewService.ListReviews:output_type -> testapi.starwars.ListReviewsResponse
	0, // 8: testapi.starwars.ReviewService.CreateReview:output_type -> testapi.starwars.Review
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_starwars_review_proto_init() }
func file_starwars_review_proto_init() {
	if File_starwars_review_proto != nil {
		return
	}
	file_starwars_type_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_starwars_review_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Review); i {
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
		file_starwars_review_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListReviewsRequest); i {
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
		file_starwars_review_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListReviewsResponse); i {
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
		file_starwars_review_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateReviewRequest); i {
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
		file_starwars_review_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Color); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_starwars_review_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_starwars_review_proto_goTypes,
		DependencyIndexes: file_starwars_review_proto_depIdxs,
		MessageInfos:      file_starwars_review_proto_msgTypes,
	}.Build()
	File_starwars_review_proto = out.File
	file_starwars_review_proto_rawDesc = nil
	file_starwars_review_proto_goTypes = nil
	file_starwars_review_proto_depIdxs = nil
}
