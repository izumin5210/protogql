package loader

import (
	user_pb "apis/go/user"
	"context"
	"time"
)

func (l *Loaders) UserByID(ctx context.Context) *UserLoader {
	return NewUserLoader(UserLoaderConfig{
		Fetch: func(keys []uint64) ([]*user_pb.User, []error) {
			req := &user_pb.BatchGetUsersRequest{UserIds: keys}
			resp, err := l.UserClient.BatchGetUsers(ctx, req)
			if err != nil {
				return nil, []error{err}
			}

			return resp.GetUsers(), nil
		},
		Wait:     10 * time.Millisecond,
		MaxBatch: 10,
	})
}
