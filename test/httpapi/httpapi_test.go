package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/kulinsky/gotestexample/internal/app/command"
	"github.com/kulinsky/gotestexample/internal/app/query"
	"github.com/kulinsky/gotestexample/internal/di"
	"github.com/kulinsky/gotestexample/internal/infra/httpapi"
	"github.com/kulinsky/gotestexample/internal/infra/idgenerator"
	"github.com/kulinsky/gotestexample/internal/infra/inmemory"
)

type Req struct {
	URL string `json:"url"`
}

func newSut(t *testing.T) (*httpexpect.Expect, func()) {
	repo := inmemory.New()
	idGen := idgenerator.NanoIDGenerator{}

	cmdCreateURL := command.NewCreateShortURLCommand(idGen, repo)
	queryGetFullURL := query.NewGetFullURLQuery(repo)

	container := di.New(cmdCreateURL, queryGetFullURL)

	router := httpapi.InitRouter(container)

	server := httptest.NewServer(router)

	sut := httpexpect.Default(t, server.URL)
	return sut, server.Close
}

func TestHttpApi(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Parallel()

	t.Run("create short url should retrun id", func(t *testing.T) {
		t.Parallel()

		// Given
		url := "https://google.com"

		sut, closer := newSut(t)
		defer closer()

		// When
		resp := sut.POST("/").
			WithJSON(Req{URL: url}).
			Expect().
			Status(http.StatusCreated).JSON().Object()

		// Then
		resp.Value("id").String().NotEmpty()
	})

	t.Run("get by existing key should return url", func(t *testing.T) {
		t.Parallel()

		// Given
		url := "https://google.com"

		sut, closer := newSut(t)
		defer closer()

		id := sut.POST("/").
			WithJSON(Req{URL: url}).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			Value("id").String().Raw()

		// When
		resp := sut.GET("/" + id).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("url")

			// Then
		resp.String().IsEqual(url)
	})

	t.Run("get by unknown id should return error", func(t *testing.T) {
		t.Parallel()

		// Given
		sut, closer := newSut(t)
		defer closer()

		// When
		resp := sut.GET("/unkNownID").
			Expect().
			Status(http.StatusNotFound).
			JSON().
			Object().
			Value("error")

		// Then
		resp.String().Contains("not found")
	})

	t.Run("create with invalid url should return error", func(t *testing.T) {
		t.Parallel()

		// Given
		sut, closer := newSut(t)
		defer closer()

		// When
		resp := sut.POST("/").
			WithJSON(Req{URL: "invalid-url"}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			Value("error")

		// Then
		resp.String().Contains("invalid url")
	})

	t.Run("create with invalid request body should return error", func(t *testing.T) {
		t.Parallel()

		// Given
		sut, closer := newSut(t)
		defer closer()

		// When
		resp := sut.POST("/").
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			Value("error")

		// Then
		resp.String().Contains("invalid request body")
	})
}
