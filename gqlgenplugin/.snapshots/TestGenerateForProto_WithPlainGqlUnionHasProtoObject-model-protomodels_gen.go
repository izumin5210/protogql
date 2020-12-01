// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	hello_pb "apis/go/hello"
	"fmt"
	"io"
	"strconv"
)

type Hello struct {
	ID            uint64
	Message       string
	UserID        uint64
	ReplyHelloIds []uint64
}

func HelloListFromRepeatedProto(in []*hello_pb.Hello) []*Hello {
	out := make([]*Hello, len(in))
	for i, m := range in {
		out[i] = HelloFromProto(m)
	}
	return out
}

func HelloFromProto(in *hello_pb.Hello) *Hello {
	out := &Hello{
		ID:            in.GetId(),
		Message:       in.GetMessage(),
		UserID:        in.GetUserId(),
		ReplyHelloIds: in.GetReplyHelloIds(),
	}

	return out
}

func HelloListToRepeatedProto(in []*Hello) []*hello_pb.Hello {
	out := make([]*hello_pb.Hello, len(in))
	for i, m := range in {
		out[i] = HelloToProto(m)
	}
	return out
}

func HelloToProto(in *Hello) *hello_pb.Hello {
	out := &hello_pb.Hello{
		Id:            in.ID,
		Message:       in.Message,
		UserId:        in.UserID,
		ReplyHelloIds: in.ReplyHelloIds,
	}

	return out
}

type HelloInput struct {
	ID            uint64
	Message       string
	UserID        uint64
	ReplyHelloIds []uint64
}

func HelloInputListFromRepeatedProto(in []*hello_pb.Hello) []*HelloInput {
	out := make([]*HelloInput, len(in))
	for i, m := range in {
		out[i] = HelloInputFromProto(m)
	}
	return out
}

func HelloInputFromProto(in *hello_pb.Hello) *HelloInput {
	out := &HelloInput{
		ID:            in.GetId(),
		Message:       in.GetMessage(),
		UserID:        in.GetUserId(),
		ReplyHelloIds: in.GetReplyHelloIds(),
	}

	return out
}

func HelloInputListToRepeatedProto(in []*HelloInput) []*hello_pb.Hello {
	out := make([]*hello_pb.Hello, len(in))
	for i, m := range in {
		out[i] = HelloInputToProto(m)
	}
	return out
}

func HelloInputToProto(in *HelloInput) *hello_pb.Hello {
	out := &hello_pb.Hello{
		Id:            in.ID,
		Message:       in.Message,
		UserId:        in.UserID,
		ReplyHelloIds: in.ReplyHelloIds,
	}

	return out
}

type User struct {
	ID   uint64
	Name string
}

func UserListFromRepeatedProto(in []*hello_pb.User) []*User {
	out := make([]*User, len(in))
	for i, m := range in {
		out[i] = UserFromProto(m)
	}
	return out
}

func UserFromProto(in *hello_pb.User) *User {
	out := &User{
		ID:   in.GetId(),
		Name: in.GetName(),
	}

	return out
}

func UserListToRepeatedProto(in []*User) []*hello_pb.User {
	out := make([]*hello_pb.User, len(in))
	for i, m := range in {
		out[i] = UserToProto(m)
	}
	return out
}

func UserToProto(in *User) *hello_pb.User {
	out := &hello_pb.User{
		Id:   in.ID,
		Name: in.Name,
	}

	return out
}

type UserInput struct {
	ID   uint64
	Name string
}

func UserInputListFromRepeatedProto(in []*hello_pb.User) []*UserInput {
	out := make([]*UserInput, len(in))
	for i, m := range in {
		out[i] = UserInputFromProto(m)
	}
	return out
}

func UserInputFromProto(in *hello_pb.User) *UserInput {
	out := &UserInput{
		ID:   in.GetId(),
		Name: in.GetName(),
	}

	return out
}

func UserInputListToRepeatedProto(in []*UserInput) []*hello_pb.User {
	out := make([]*hello_pb.User, len(in))
	for i, m := range in {
		out[i] = UserInputToProto(m)
	}
	return out
}

func UserInputToProto(in *UserInput) *hello_pb.User {
	out := &hello_pb.User{
		Id:   in.ID,
		Name: in.Name,
	}

	return out
}

type CreateHelloSuccess_Proto struct {
	Hello *hello_pb.Hello
}

func CreateHelloSuccessListFromRepeatedProto(in []*CreateHelloSuccess_Proto) []*CreateHelloSuccess {
	out := make([]*CreateHelloSuccess, len(in))
	for i, m := range in {
		out[i] = CreateHelloSuccessFromProto(m)
	}
	return out
}

func CreateHelloSuccessFromProto(in *CreateHelloSuccess_Proto) *CreateHelloSuccess {
	return &CreateHelloSuccess{
		Hello: HelloFromProto(in.Hello),
	}
}

func CreateHelloSuccessListToRepeatedProto(in []*CreateHelloSuccess) []*CreateHelloSuccess_Proto {
	out := make([]*CreateHelloSuccess_Proto, len(in))
	for i, m := range in {
		out[i] = CreateHelloSuccessToProto(m)
	}
	return out
}

func CreateHelloSuccessToProto(in *CreateHelloSuccess) *CreateHelloSuccess_Proto {
	return &CreateHelloSuccess_Proto{
		Hello: HelloToProto(in.Hello),
	}
}

type HelloMessageInvalid_Proto struct {
	Hello   *hello_pb.Hello
	Message string
}

func HelloMessageInvalidListFromRepeatedProto(in []*HelloMessageInvalid_Proto) []*HelloMessageInvalid {
	out := make([]*HelloMessageInvalid, len(in))
	for i, m := range in {
		out[i] = HelloMessageInvalidFromProto(m)
	}
	return out
}

func HelloMessageInvalidFromProto(in *HelloMessageInvalid_Proto) *HelloMessageInvalid {
	return &HelloMessageInvalid{
		Hello:   HelloFromProto(in.Hello),
		Message: in.Message,
	}
}

func HelloMessageInvalidListToRepeatedProto(in []*HelloMessageInvalid) []*HelloMessageInvalid_Proto {
	out := make([]*HelloMessageInvalid_Proto, len(in))
	for i, m := range in {
		out[i] = HelloMessageInvalidToProto(m)
	}
	return out
}

func HelloMessageInvalidToProto(in *HelloMessageInvalid) *HelloMessageInvalid_Proto {
	return &HelloMessageInvalid_Proto{
		Hello:   HelloToProto(in.Hello),
		Message: in.Message,
	}
}

type UserRole struct {
	Proto hello_pb.User_Role
}

func UserRoleListFromRepeatedProto(in []hello_pb.User_Role) []*UserRole {
	out := make([]*UserRole, len(in))
	for i, m := range in {
		out[i] = UserRoleFromProto(m)
	}
	return out
}

func UserRoleFromProto(in hello_pb.User_Role) *UserRole {
	return &UserRole{Proto: in}
}

func UserRoleListToRepeatedProto(in []*UserRole) []hello_pb.User_Role {
	out := make([]hello_pb.User_Role, len(in))
	for i, m := range in {
		out[i] = UserRoleToProto(m)
	}
	return out
}

func UserRoleToProto(in *UserRole) hello_pb.User_Role {
	return in.Proto
}

func (e UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.Proto.String()))
}

func (e *UserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	e.Proto = hello_pb.User_Role(hello_pb.User_Role_value[str])
	return nil
}

type CreateHelloPayload_Proto struct {
	Raw CreateHelloPayload
}

func (*Hello) IsCreateHelloPayload() {}

func (u *CreateHelloPayload_Proto) GetHello() *hello_pb.Hello {
	if m, ok := u.Raw.(*Hello); ok {
		return HelloToProto(m)
	}
	return nil
}

func (u *CreateHelloPayload_Proto) GetHelloMessageInvalid() *HelloMessageInvalid_Proto {
	if m, ok := u.Raw.(*HelloMessageInvalid); ok {
		return HelloMessageInvalidToProto(m)
	}
	return nil
}

func (u *CreateHelloPayload_Proto) GetUserError() *UserError {
	if m, ok := u.Raw.(*UserError); ok {
		return m
	}
	return nil
}

func CreateHelloPayloadFromProto(in *CreateHelloPayload_Proto) CreateHelloPayload {
	return in.Raw
}

func CreateHelloPayloadToProto(in CreateHelloPayload) *CreateHelloPayload_Proto {
	return &CreateHelloPayload_Proto{Raw: in}
}

