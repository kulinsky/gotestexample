package main

import (
	"log"
	"net/http"

	"github.com/kulinsky/gotestexample/internal/app/command"
	"github.com/kulinsky/gotestexample/internal/app/query"
	"github.com/kulinsky/gotestexample/internal/di"
	"github.com/kulinsky/gotestexample/internal/infra/httpapi"
	"github.com/kulinsky/gotestexample/internal/infra/idprovider"
	"github.com/kulinsky/gotestexample/internal/infra/inmemory"
)

func main() {
	idp := idprovider.NanoIDProvider{}
	repo := inmemory.New()

	cmdCreateURL := command.NewCreateShortURLCmd(idp, repo)
	queryGetURL := query.NewGetLongURLQuery(repo)

	container := di.New(cmdCreateURL, queryGetURL)

	router := httpapi.InitRouter(container)

	//nolint:gosec // it's ok
	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}
