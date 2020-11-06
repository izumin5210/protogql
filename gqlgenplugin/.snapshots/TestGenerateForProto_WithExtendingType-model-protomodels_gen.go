// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	todo_pb "apis/go/todo"
	user_pb "apis/go/user"
	"fmt"
	"io"
	"strconv"
)

type Task struct {
	ID uint64

	Title string

	Status *TaskStatus

	AssigneeIds []uint64

	AuthorID uint64

	Assignees []*User

	Author *User
}

func TaskListFromRepeatedProto(in []*todo_pb.Task) []*Task {
	out := make([]*Task, len(in))
	for i, m := range in {
		out[i] = TaskFromProto(m)
	}
	return out
}

func TaskFromProto(in *todo_pb.Task) *Task {
	return &Task{

		ID: in.Id,

		Title: in.Title,

		Status: TaskStatusFromProto(in.Status),

		AssigneeIds: in.AssigneeIds,

		AuthorID: in.AuthorId,
	}
}

func TaskListToRepeatedProto(in []*Task) []*todo_pb.Task {
	out := make([]*todo_pb.Task, len(in))
	for i, m := range in {
		out[i] = TaskToProto(m)
	}
	return out
}

func TaskToProto(in *Task) *todo_pb.Task {
	return &todo_pb.Task{

		Id: in.ID,

		Title: in.Title,

		Status: TaskStatusToProto(in.Status),

		AssigneeIds: in.AssigneeIds,

		AuthorId: in.AuthorID,
	}
}

type TaskInput struct {
	ID uint64

	Title string

	Status *TaskStatus

	AssigneeIds []uint64

	AuthorID uint64
}

func TaskInputListFromRepeatedProto(in []*todo_pb.Task) []*TaskInput {
	out := make([]*TaskInput, len(in))
	for i, m := range in {
		out[i] = TaskInputFromProto(m)
	}
	return out
}

func TaskInputFromProto(in *todo_pb.Task) *TaskInput {
	return &TaskInput{

		ID: in.Id,

		Title: in.Title,

		Status: TaskStatusFromProto(in.Status),

		AssigneeIds: in.AssigneeIds,

		AuthorID: in.AuthorId,
	}
}

func TaskInputListToRepeatedProto(in []*TaskInput) []*todo_pb.Task {
	out := make([]*todo_pb.Task, len(in))
	for i, m := range in {
		out[i] = TaskInputToProto(m)
	}
	return out
}

func TaskInputToProto(in *TaskInput) *todo_pb.Task {
	return &todo_pb.Task{

		Id: in.ID,

		Title: in.Title,

		Status: TaskStatusToProto(in.Status),

		AssigneeIds: in.AssigneeIds,

		AuthorId: in.AuthorID,
	}
}

type User struct {
	ID uint64

	FullName string

	Role *UserRole

	AssignedTasks []*Task
}

func UserListFromRepeatedProto(in []*user_pb.User) []*User {
	out := make([]*User, len(in))
	for i, m := range in {
		out[i] = UserFromProto(m)
	}
	return out
}

func UserFromProto(in *user_pb.User) *User {
	return &User{

		ID: in.Id,

		FullName: in.FullName,

		Role: UserRoleFromProto(in.Role),
	}
}

func UserListToRepeatedProto(in []*User) []*user_pb.User {
	out := make([]*user_pb.User, len(in))
	for i, m := range in {
		out[i] = UserToProto(m)
	}
	return out
}

func UserToProto(in *User) *user_pb.User {
	return &user_pb.User{

		Id: in.ID,

		FullName: in.FullName,

		Role: UserRoleToProto(in.Role),
	}
}

type UserInput struct {
	ID uint64

	FullName string

	Role *UserRole
}

func UserInputListFromRepeatedProto(in []*user_pb.User) []*UserInput {
	out := make([]*UserInput, len(in))
	for i, m := range in {
		out[i] = UserInputFromProto(m)
	}
	return out
}

func UserInputFromProto(in *user_pb.User) *UserInput {
	return &UserInput{

		ID: in.Id,

		FullName: in.FullName,

		Role: UserRoleFromProto(in.Role),
	}
}

func UserInputListToRepeatedProto(in []*UserInput) []*user_pb.User {
	out := make([]*user_pb.User, len(in))
	for i, m := range in {
		out[i] = UserInputToProto(m)
	}
	return out
}

func UserInputToProto(in *UserInput) *user_pb.User {
	return &user_pb.User{

		Id: in.ID,

		FullName: in.FullName,

		Role: UserRoleToProto(in.Role),
	}
}

type TaskStatus struct {
	Proto todo_pb.Task_Status
}

func TaskStatusListFromRepeatedProto(in []todo_pb.Task_Status) []*TaskStatus {
	out := make([]*TaskStatus, len(in))
	for i, m := range in {
		out[i] = TaskStatusFromProto(m)
	}
	return out
}

func TaskStatusFromProto(in todo_pb.Task_Status) *TaskStatus {
	return &TaskStatus{Proto: in}
}

func TaskStatusListToRepeatedProto(in []*TaskStatus) []todo_pb.Task_Status {
	out := make([]todo_pb.Task_Status, len(in))
	for i, m := range in {
		out[i] = TaskStatusToProto(m)
	}
	return out
}

func TaskStatusToProto(in *TaskStatus) todo_pb.Task_Status {
	return in.Proto
}

func (e TaskStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.Proto.String()))
}

func (e *TaskStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	e.Proto = todo_pb.Task_Status(todo_pb.Task_Status_value[str])
	return nil
}

type UserRole struct {
	Proto user_pb.User_Role
}

func UserRoleListFromRepeatedProto(in []user_pb.User_Role) []*UserRole {
	out := make([]*UserRole, len(in))
	for i, m := range in {
		out[i] = UserRoleFromProto(m)
	}
	return out
}

func UserRoleFromProto(in user_pb.User_Role) *UserRole {
	return &UserRole{Proto: in}
}

func UserRoleListToRepeatedProto(in []*UserRole) []user_pb.User_Role {
	out := make([]user_pb.User_Role, len(in))
	for i, m := range in {
		out[i] = UserRoleToProto(m)
	}
	return out
}

func UserRoleToProto(in *UserRole) user_pb.User_Role {
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

	e.Proto = user_pb.User_Role(user_pb.User_Role_value[str])
	return nil
}

