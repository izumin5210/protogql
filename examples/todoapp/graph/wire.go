//+build wireinject

package graph

import (
	"context"

	"github.com/google/wire"
)

func NewApp(context.Context) (*App, func(), error) {
	wire.Build(
		appSet,
	)
	return nil, nil, nil
}
