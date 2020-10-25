package model

import (
	"fmt"
	"io"
	"strconv"

	"testdata/task/api"
)

type Task struct {
	Id          uint64
	Title       string
	Status      *TaskStatus
	AssigneeIds []uint64
}

func TaskListFromProto(in []*api.Task) []*Task {
	out := make([]*Task, len(in))
	for i, m := range in {
		out[i] = TaskFromProto(m)
	}
	return out
}

func TaskFromProto(in *api.Task) *Task {
	out := new(Task)
	out.Id = in.Id
	out.Title = in.Title
	out.Status = TaskStatusFromProto(in.Status)
	out.AssigneeIds = in.AssigneeIds
	return out
}

func (t *Task) Proto() *api.Task {
	out := new(api.Task)
	out.Id = t.Id
	out.Title = t.Title
	out.Status = t.Status.Proto()
	out.AssigneeIds = t.AssigneeIds
	return out
}

type TaskStatus struct {
	value api.Task_Status
}

func TaskStatusFromProto(in api.Task_Status) *TaskStatus {
	out := new(TaskStatus)
	out.value = in
	return out
}

func (s TaskStatus) Proto() api.Task_Status {
	return s.value
}

func (s TaskStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(s.value.String()))
}

func (s *TaskStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	s.value = api.Task_Status(api.Task_Status_value[str])
	return nil
}
