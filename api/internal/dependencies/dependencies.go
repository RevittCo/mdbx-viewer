package dependencies

import (
	"github.com/RevittConsulting/mdbx-viewer/config"
	"github.com/RevittConsulting/mdbx-viewer/internal/buckets"
	"github.com/RevittConsulting/mdbx-viewer/internal/db_mdbx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Dependencies struct {
	cfg          *config.Config
	globalRouter chi.Router
	routerApiV1  chi.Router
	dbMdbx       *db_mdbx.MDBX

	bucketsApi     *buckets.Api
	bucketsService *buckets.Service
}

func NewDependencies(cfg *config.Config) *Dependencies {
	deps := &Dependencies{}

	// CONFIG
	deps.cfg = cfg

	// ROUTER
	deps.globalRouter = chi.NewRouter()
	deps.routerApiV1 = chi.NewRouter()

	// CORS
	deps.globalRouter.Use(cors.AllowAll().Handler)

	// MDBX
	deps.dbMdbx = db_mdbx.New()

	// BUCKETS
	deps.bucketsApi = buckets.NewApi(deps)
	deps.bucketsService = buckets.NewService(deps.dbMdbx)

	// MOUNT ROUTERS TO GLOBAL
	deps.globalRouter.Mount("/api/v1", deps.routerApiV1)

	return deps
}

func (d Dependencies) GetConfig() *config.Config {
	return d.cfg
}

func (d Dependencies) GetRouter() chi.Router {
	return d.globalRouter
}

func (d Dependencies) GetRouterApiV1() chi.Router {
	return d.routerApiV1
}

func (d Dependencies) GetDbMdbx() *db_mdbx.MDBX {
	return d.dbMdbx
}

func (d Dependencies) GetBucketsApi() *buckets.Api {
	return d.bucketsApi
}

func (d Dependencies) GetBucketsService() *buckets.Service {
	return d.bucketsService
}
