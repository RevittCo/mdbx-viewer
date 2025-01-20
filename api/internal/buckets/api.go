package buckets

import (
	"encoding/json"
	"github.com/RevittConsulting/mdbx-viewer/config"
	"github.com/RevittConsulting/mdbx-viewer/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
		r.Get("/buckets/data-source", api.GetDataSource)
		r.Post("/buckets/open", api.Open)
		r.Route("/buckets/{name}", func(r chi.Router) {
			r.Get("/read/{num}/{len}", api.Read)
			r.Get("/ws/read/{num}/{len}", api.StartStreamRead)
		})
	})
}

func (api *Api) GetDataSource(w http.ResponseWriter, r *http.Request) {
	var err error
	dataSource, err := api.deps.GetBucketsService().GetDataSource()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, dataSource)
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

	w.WriteHeader(http.StatusOK)
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

func (api *Api) StartStreamRead(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

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

	for {
		list, err := api.deps.GetBucketsService().Read(bucketName, uint64(pageNum), uint64(pageLen))
		if err != nil {
			log.Println("error polling chain data:", err)
			continue
		}

		bytes, err := json.Marshal(list)
		if err != nil {
			log.Println("error marshalling chain data:", err)
			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: client disconnected unexpectedly: %v", err)
			} else {
				log.Println("error writing message:", err)
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}
