package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"

	"apis/go/todo"
	"apis/go/user"
	"todoapp/graph"
	"todoapp/graph/generated"
	"todoapp/todo"
	"todoapp/user"
)

func TestTodoApp(t *testing.T) {
	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskServerAddr := startGrpcServer(ctx, t, &wg, func(s *grpc.Server) {
		todo_pb.RegisterTaskServiceServer(s, todo.NewTaskServiceServer())
	})
	os.Setenv("TASK_PORT", fmt.Sprint(taskServerAddr.Port))

	userServerAddr := startGrpcServer(ctx, t, &wg, func(s *grpc.Server) {
		user_pb.RegisterUserServiceServer(s, user.NewUserServiceServer())
	})
	os.Setenv("USER_PORT", fmt.Sprint(userServerAddr.Port))

	gqlServerAddr := startGqlServer(ctx, t, &wg)
	user.Users[1] = &user_pb.User{Id: 1, FullName: "Test User1"}
	user.Users[2] = &user_pb.User{Id: 2, FullName: "Test User2"}
	todo.Tasks[1] = &todo_pb.Task{Id: 1, AuthorId: 2, AssigneeIds: []uint64{2}, Title: "Test task 1", Status: todo_pb.Task_DONE}
	todo.Tasks[2] = &todo_pb.Task{Id: 2, AuthorId: 1, AssigneeIds: []uint64{1, 2}, Title: "Test task 2", Status: todo_pb.Task_IN_PROGRESS}

	if diff := cmp.Diff(
		map[string]interface{}{
			"data": map[string]interface{}{
				"tasks": []interface{}{
					map[string]interface{}{
						"id":    float64(1),
						"title": "Test task 1",
						"author": map[string]interface{}{
							"id":       float64(2),
							"fullName": "Test User2",
						},
						"assignees": []interface{}{
							map[string]interface{}{
								"id":       float64(2),
								"fullName": "Test User2",
							},
						},
					},
					map[string]interface{}{
						"id":    float64(2),
						"title": "Test task 2",
						"author": map[string]interface{}{
							"id":       float64(1),
							"fullName": "Test User1",
						},
						"assignees": []interface{}{
							map[string]interface{}{
								"id":       float64(1),
								"fullName": "Test User1",
							},
							map[string]interface{}{
								"id":       float64(2),
								"fullName": "Test User2",
							},
						},
					},
				},
			},
		},
		execQuery(ctx, t, gqlServerAddr.Port, query{Query: `
query {
  tasks {
    id
    title
    author {
      ...User
    }
    assignees {
      ...User
    }
  }
}
fragment User on User {
  id
  fullName
}
`}),
	); diff != "" {
		t.Errorf("query returned(-want +got):\n%s", diff)
	}
}

type query struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func execQuery(ctx context.Context, t *testing.T, port int, query query) map[string]interface{} {
	cli := http.Client{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		t.Errorf("failed to encode request: %v", err)
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://localhost:%d/query", port), &buf)
	if err != nil {
		t.Errorf("failed to create request: %v", err)
		return nil
	}
	req.Header.Add("content-type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Errorf("failed to send request: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("got status code %d, want %d", got, want)
	}

	result := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
		return nil
	}

	return result
}

func startGqlServer(ctx context.Context, t *testing.T, wg *sync.WaitGroup) *net.TCPAddr {
	app, cleanup, err := graph.NewApp(ctx)
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	cfg := generated.Config{
		Resolvers: app.Resolver,
	}
	mux := http.NewServeMux()
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	mux.Handle("/query", app.Loaders.Middleware(gqlHandler))

	lis, err := net.Listen("tcp", ":")
	if err != nil {
		t.Errorf("failed to listen: %v", err)
		return nil
	}

	s := &http.Server{Handler: mux}
	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cleanup()
		if err := s.Serve(lis); err != nil && err.Error() != "http: Server closed" {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	return lis.Addr().(*net.TCPAddr)
}

func startGrpcServer(ctx context.Context, t *testing.T, wg *sync.WaitGroup, register func(s *grpc.Server)) *net.TCPAddr {
	lis, err := net.Listen("tcp", ":")
	if err != nil {
		t.Errorf("failed to listen: %v", err)
		return nil
	}
	s := grpc.NewServer()
	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()
	register(s)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	return lis.Addr().(*net.TCPAddr)
}
