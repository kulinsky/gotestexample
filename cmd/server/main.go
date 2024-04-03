package main

import (
	"log"
	"net/http"

	"github.com/kulinsky/gotestexample/internal/app/command"
	"github.com/kulinsky/gotestexample/internal/app/query"
	"github.com/kulinsky/gotestexample/internal/di"
	"github.com/kulinsky/gotestexample/internal/infra/httpapi"
	"github.com/kulinsky/gotestexample/internal/infra/idgenerator"
	"github.com/kulinsky/gotestexample/internal/infra/inmemory"
)

func main() {
	idGen := idgenerator.NanoIDGenerator{}
	repo := inmemory.New()

	cmdCreateURL := command.NewCreateShortURLCommand(idGen, repo)
	queryGetURL := query.NewGetFullURLQuery(repo)

	container := di.New(cmdCreateURL, queryGetURL)

	router := httpapi.InitRouter(container)

	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}
