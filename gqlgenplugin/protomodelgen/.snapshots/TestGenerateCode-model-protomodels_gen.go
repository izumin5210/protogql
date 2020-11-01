// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"apis/go/task"
	"fmt"
	"io"
	"strconv"
)

type Task struct {
	ID uint64

	Title string

	Status *TaskStatus

	AssigneeIds []uint64
}

func TaskListFromRepeatedProto(in []*task.Task) []*Task {
	out := make([]*Task, len(in))
	for i, m := range in {
		out[i] = TaskFromProto(m)
	}
	return out
}

func TaskFromProto(in *task.Task) *Task {
	return &Task{

		ID: in.Id,

		Title: in.Title,

		Status: TaskStatusFromProto(in.Status),

		AssigneeIds: in.AssigneeIds,
	}
}

func TaskListToRepeatedProto(in []*Task) []*task.Task {
	out := make([]*task.Task, len(in))
	for i, m := range in {
		out[i] = TaskToProto(m)
	}
	return out
}

func TaskToProto(in *Task) *task.Task {
	return &task.Task{

		Id: in.ID,

		Title: in.Title,

		Status: TaskStatusToProto(in.Status),

		AssigneeIds: in.AssigneeIds,
	}
}

type TaskInput struct {
	ID uint64

	Title string

	Status *TaskStatus

	AssigneeIds []uint64
}

func TaskInputListFromRepeatedProto(in []*task.Task) []*TaskInput {
	out := make([]*TaskInput, len(in))
	for i, m := range in {
		out[i] = TaskInputFromProto(m)
	}
	return out
}

func TaskInputFromProto(in *task.Task) *TaskInput {
	return &TaskInput{

		ID: in.Id,

		Title: in.Title,

		Status: TaskStatusFromProto(in.Status),

		AssigneeIds: in.AssigneeIds,
	}
}

func TaskInputListToRepeatedProto(in []*TaskInput) []*task.Task {
	out := make([]*task.Task, len(in))
	for i, m := range in {
		out[i] = TaskInputToProto(m)
	}
	return out
}

func TaskInputToProto(in *TaskInput) *task.Task {
	return &task.Task{

		Id: in.ID,

		Title: in.Title,

		Status: TaskStatusToProto(in.Status),

		AssigneeIds: in.AssigneeIds,
	}
}

type TaskStatus struct {
	Proto task.Task_Status
}

func TaskStatusListFromRepeatedProto(in []task.Task_Status) []*TaskStatus {
	out := make([]*TaskStatus, len(in))
	for i, m := range in {
		out[i] = TaskStatusFromProto(m)
	}
	return out
}

func TaskStatusFromProto(in task.Task_Status) *TaskStatus {
	return &TaskStatus{Proto: in}
}

func TaskStatusListToRepeatedProto(in []*TaskStatus) []task.Task_Status {
	out := make([]task.Task_Status, len(in))
	for i, m := range in {
		out[i] = TaskStatusToProto(m)
	}
	return out
}

func TaskStatusToProto(in *TaskStatus) task.Task_Status {
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

	e.Proto = task.Task_Status(task.Task_Status_value[str])
	return nil
}

