package query

import (
	"context"
	"errors"
	"fmt"

	"github.com/kulinsky/gotestexample/internal/common"
)

var ErrRepository = fmt.Errorf("query repo error: %w", common.ErrTechnical)

//go:generate mockgen -destination=../../../mocks/get_full_url_repo.go -mock_names=Repository=MockQueryRepo -package=mocks integrationtest/internal/app/query Repository
type Repository interface {
	Get(ctx context.Context, id string) (string, error)
}

type GetFullUrlQuery struct {
	repo Repository
}

func NewGetFullURLQuery(repo Repository) *GetFullUrlQuery {
	return &GetFullUrlQuery{
		repo: repo,
	}
}

func (q *GetFullUrlQuery) Execute(ctx context.Context, id string) (string, error) {
	res, err := q.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return "", fmt.Errorf("%w: url not found: %s", err, id)
		}

		return "", fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return res, nil
}
