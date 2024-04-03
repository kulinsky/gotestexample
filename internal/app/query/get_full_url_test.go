package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kulinsky/gotestexample/internal/app/query"
	"github.com/kulinsky/gotestexample/internal/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryStub struct {
	mock.Mock
}

func (s *repositoryStub) Get(ctx context.Context, id string) (string, error) {
	args := s.Called(id)

	return args.Get(0).(string), args.Error(1)
}

func TestGetFullUrlQueryExecute(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	assert := assert.New(t)

	type in struct {
		id string
	}

	type out struct {
		err error
		id  string
	}

	tests := []struct {
		setup  func(*repositoryStub, *in)
		assert func(*out)
		name   string
		in     in
	}{
		{
			name: "happy path",
			in:   in{id: "existing-key"},
			setup: func(repo *repositoryStub, in *in) {
				repo.On("Get", in.id).Return("https://google.com", nil)
			},
			assert: func(out *out) {
				assert.NoError(out.err)
				assert.Equal(out.id, "https://google.com")
			},
		},
		{
			name: "not found",
			in:   in{id: "unknown-id"},
			setup: func(repo *repositoryStub, in *in) {
				repo.On("Get", in.id).Return("", common.ErrNotFound)
			},
			assert: func(out *out) {
				assert.ErrorIs(out.err, common.ErrNotFound)
				assert.Empty(out.id)
			},
		},
		{
			name: "error while get from repo",
			in:   in{id: "existing-key"},
			setup: func(repo *repositoryStub, in *in) {
				repo.On("Get", in.id).Return("", errors.New("repository down"))
			},
			assert: func(out *out) {
				assert.ErrorIs(out.err, common.ErrTechnical)
				assert.Empty(out.id)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &repositoryStub{}
			sut := query.NewGetFullURLQuery(repo)

			tt.setup(repo, &tt.in)

			res, err := sut.Execute(ctx, tt.in.id)

			tt.assert(&out{id: res, err: err})
		})
	}
}
