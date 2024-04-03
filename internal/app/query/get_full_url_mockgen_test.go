package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kulinsky/gotestexample/internal/app/query"
	"github.com/kulinsky/gotestexample/internal/common"
	"github.com/kulinsky/gotestexample/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetFullUrlQueryExecuteMockgen(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	type in struct {
		err     error
		id      string
		longURL string
	}

	type out struct {
		err     error
		longURL string
	}

	setup := func(ctx context.Context, in *in) *query.GetFullUrlQuery {
		ctrl := gomock.NewController(t)
		repo := mocks.NewMockQueryRepo(ctrl)
		repo.EXPECT().Get(ctx, in.id).Return(in.longURL, in.err).Times(1)

		return query.NewGetFullURLQuery(repo)
	}

	tests := []struct {
		setup  func(context.Context, *in) *query.GetFullUrlQuery
		assert func(*out)
		name   string
		in     in
	}{
		{
			name:  "happy path",
			in:    in{id: "existing-key", longURL: "https://google.com"},
			setup: setup,
			assert: func(out *out) {
				assert.NoError(out.err)
				assert.Equal(out.longURL, "https://google.com")
			},
		},
		{
			name:  "not found",
			in:    in{id: "unknown-id", err: common.ErrNotFound},
			setup: setup,
			assert: func(out *out) {
				assert.ErrorIs(out.err, common.ErrNotFound)
				assert.Empty(out.longURL)
			},
		},
		{
			name:  "error while get from repo",
			in:    in{id: "existing-key", err: errors.New("repository down")},
			setup: setup,
			assert: func(out *out) {
				assert.ErrorIs(out.err, common.ErrTechnical)
				assert.Empty(out.longURL)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx := context.TODO()
			sut := tt.setup(ctx, &tt.in)

			// When
			res, err := sut.Execute(ctx, tt.in.id)

			// Then
			tt.assert(&out{longURL: res, err: err})
		})
	}
}
