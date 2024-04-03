package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kulinsky/gotestexample/internal/app/command"
	"github.com/kulinsky/gotestexample/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) Save(ctx context.Context, id, fullURL string) error {
	args := m.Called(id, fullURL)

	return args.Error(0)
}

type idGeneratorStub struct {
	ID string
}

func (d idGeneratorStub) Generate() string {
	return d.ID
}

func TestCreateShortURLCommandTable(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	assert := assert.New(t)

	type in struct {
		fullURL string
	}

	type out struct {
		err error
		id  string
	}

	tests := []struct {
		setup  func(*repositoryMock, *idGeneratorStub, *in)
		assert func(*repositoryMock, *out)
		name   string
		in     in
	}{
		{
			name: "successfully save",
			in:   in{fullURL: "https://google.com"},
			setup: func(repo *repositoryMock, idGen *idGeneratorStub, in *in) {
				idGen.ID = "1"
				repo.On("Save", "1", in.fullURL).Return(nil)
			},
			assert: func(repo *repositoryMock, out *out) {
				assert.NoError(out.err)
				assert.Equal("1", out.id)
				repo.AssertExpectations(t)
			},
		},
		{
			name: "invalid url",
			in:   in{fullURL: "this is invalid url"},
			setup: func(_ *repositoryMock, idGen *idGeneratorStub, _ *in) {
				idGen.ID = "1"
			},
			assert: func(_ *repositoryMock, out *out) {
				assert.ErrorIs(out.err, common.ErrValidation)
				assert.Empty(out.id)
			},
		},
		{
			name: "repo error on save",
			in:   in{fullURL: "https://google.com"},
			setup: func(repo *repositoryMock, idGen *idGeneratorStub, in *in) {
				idGen.ID = "1"
				repo.On("Save", "1", in.fullURL).Return(errors.New("unexpected repository error"))
			},
			assert: func(repo *repositoryMock, out *out) {
				assert.ErrorIs(out.err, common.ErrTechnical)
				assert.Empty(out.id)
				repo.AssertExpectations(t)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			idp := &idGeneratorStub{}
			repo := &repositoryMock{}
			tt.setup(repo, idp, &tt.in)

			sut := command.NewCreateShortURLCommand(idp, repo)

			// When
			res, err := sut.Execute(ctx, tt.in.fullURL)

			// Then
			tt.assert(repo, &out{id: res, err: err})
		})
	}
}
