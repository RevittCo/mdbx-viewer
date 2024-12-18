package buckets

import (
	"github.com/RevittConsulting/mdbx-viewer/config"
	"github.com/RevittConsulting/mdbx-viewer/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Deps interface {
	GetConfig() *config.Config
	GetRouterApiV1() chi.Router
	GetBucketsService() *Service
}

type Api struct {
	deps Deps
}

func NewApi(deps Deps) *Api {
	h := &Api{
		deps: deps,
	}
	h.SetupApi(deps.GetRouterApiV1())

	return h
}

func (api *Api) SetupApi(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/buckets/open", api.Open)
		r.Route("/buckets/{name}", func(r chi.Router) {
			r.Get("/read/{num}/{len}", api.Read)
		})
	})
}

func (api *Api) Open(w http.ResponseWriter, r *http.Request) {
	req := OpenReq{}
	err := utils.ReadJSON(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	list, err := api.deps.GetBucketsService().Open(req.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.WriteJSON(w, list)
}

func (api *Api) Read(w http.ResponseWriter, r *http.Request) {
	var err error
	bucketName := chi.URLParam(r, "name")
	pageNum, err := strconv.Atoi(chi.URLParam(r, "num"))
	if err != nil {
		http.Error(w, "invalid page number", http.StatusBadRequest)
		return
	}
	pageLen, err := strconv.Atoi(chi.URLParam(r, "len"))
	if err != nil {
		http.Error(w, "invalid page length", http.StatusBadRequest)
		return
	}

	if pageLen > MaxPageLen {
		pageLen = MaxPageLen
	}

	list, err := api.deps.GetBucketsService().Read(bucketName, uint64(pageNum), uint64(pageLen))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, list)
}
