package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kulinsky/gotestexample/internal/common"
	"github.com/kulinsky/gotestexample/internal/di"
)

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

type GetResponse struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func InitRouter(container *di.Container) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		var req CreateRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			terminate(w, fmt.Errorf("%w: failed to read request body: %e", common.ErrValidation, err))

			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			terminate(w, fmt.Errorf("%w: invalid request body: %e", common.ErrValidation, err))

			return
		}

		id, err := container.CreateShortUrlCmd.Execute(ctx, req.URL)
		if err != nil {
			terminate(w, err)

			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateResponse{ID: id})
	})

	router.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		id := r.PathValue("id")

		url, err := container.GetFullUrlQuery.Execute(ctx, id)
		if err != nil {
			terminate(w, err)

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetResponse{URL: url})
	})

	return router
}

func matchErr(err error) (int, string) {
	if errors.Is(err, common.ErrValidation) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, common.ErrNotFound) {
		return http.StatusNotFound, err.Error()
	}

	return http.StatusInternalServerError, err.Error()
}

func terminate(w http.ResponseWriter, err error) {
	code, msg := matchErr(err)

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
