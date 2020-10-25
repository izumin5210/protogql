package model

import (
	"fmt"
	"io"
	"strconv"

	"testdata/task/api"
)

type User struct {
	Id          uint64
	FullName    string
	Role        *UserRole
	AssigneeIds []uint64
}

func UserListFromProto(in []*api.User) []*User {
	out := make([]*User, len(in))
	for i, m := range in {
		out[i] = UserFromProto(m)
	}
	return out
}

func UserFromProto(in *api.User) *User {
	out := new(User)
	out.Id = in.Id
	out.FullName = in.FullName
	out.Role = UserRoleFromProto(in.Role)
	return out
}

func (u *User) Proto() *api.User {
	out := new(api.User)
	out.Id = u.Id
	out.FullName = u.FullName
	out.Role = u.Role.Proto()
	return out
}

type UserRole struct {
	value api.User_Role
}

func UserRoleFromProto(in api.User_Role) *UserRole {
	out := new(UserRole)
	out.value = in
	return out
}

func (u UserRole) Proto() api.User_Role {
	return u.value
}

func (u UserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(u.value.String()))
}

func (u *UserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	u.value = api.User_Role(api.User_Role_value[str])
	return nil
}
