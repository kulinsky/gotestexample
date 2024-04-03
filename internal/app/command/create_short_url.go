//go:generate mockgen -destination=../../../mocks/create_short_url_cmd.go -package=mocks -mock_names=Repository=MockCmdRepo,IDGenerator=MockIDGenerator integrationtest/internal/app/command Repository,IDGenerator

package command

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kulinsky/gotestexample/internal/common"
)

var (
	ErrInvalidURL = fmt.Errorf("invalid url: %w", common.ErrValidation)
	ErrRepository = fmt.Errorf("command repo error: %w", common.ErrTechnical)
)

type IDGenerator interface {
	Generate() string
}

type Repository interface {
	Save(ctx context.Context, id string, full_url string) error
}

type CreateShortUrlCommand struct {
	idGenerator IDGenerator
	repo        Repository
}

func NewCreateShortURLCommand(idGen IDGenerator, repo Repository) *CreateShortUrlCommand {
	return &CreateShortUrlCommand{
		idGenerator: idGen,
		repo:        repo,
	}
}

func (cmd *CreateShortUrlCommand) Execute(ctx context.Context, rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", ErrInvalidURL
	}

	id := cmd.idGenerator.Generate()

	if err := cmd.repo.Save(ctx, id, parsedURL.String()); err != nil {
		return "", fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return id, nil
}
