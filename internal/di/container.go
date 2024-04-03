package di

import (
	"github.com/kulinsky/gotestexample/internal/app/command"
	"github.com/kulinsky/gotestexample/internal/app/query"
)

type Container struct {
	CreateShortUrlCmd *command.CreateShortUrlCommand
	GetFullUrlQuery   *query.GetFullUrlQuery
}

func New(cmdCreateShortURL *command.CreateShortUrlCommand, qGetFullURL *query.GetFullUrlQuery) *Container {
	return &Container{
		CreateShortUrlCmd: cmdCreateShortURL,
		GetFullUrlQuery:   qGetFullURL,
	}
}
