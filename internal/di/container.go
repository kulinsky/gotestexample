package di

import (
	"github.com/kulinsky/gotestexample/internal/app/command"
	"github.com/kulinsky/gotestexample/internal/app/query"
)

type Container struct {
	CreateShortURLCmd *command.CreateShortURLCmd
	GetLongURLQuery   *query.GetLongURLQuery
}

func New(cmd *command.CreateShortURLCmd, q *query.GetLongURLQuery) *Container {
	return &Container{
		CreateShortURLCmd: cmd,
		GetLongURLQuery:   q,
	}
}
