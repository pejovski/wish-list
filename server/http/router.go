package http

import (
	"github.com/gorilla/mux"
	_ "github.com/pejovski/wish-list/app/statik"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Router struct {
	router  *mux.Router
	handler *Handler
}

func NewRouter(h *Handler) *Router {
	s := &Router{
		router:  mux.NewRouter(),
		handler: h,
	}
	s.routes()
	s.swagger()
	s.health()

	return s
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.router.ServeHTTP(w, r)
}

func (rtr Router) routes() {
	rtr.router.HandleFunc("/wish-list/{user_id}", rtr.handler.GetList()).Methods("GET")
	rtr.router.HandleFunc("/wish-list/{user_id}", rtr.handler.AddItem()).Methods("POST")
	rtr.router.HandleFunc("/wish-list/{user_id}/{product_id}", rtr.handler.RemoveItem()).Methods("DELETE")
}

func (rtr Router) swagger() {
	// swagger handler
	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to find statik", err)
	}
	sh := http.FileServer(statikFS)
	rtr.router.Handle("/", sh).Methods("GET")
	rtr.router.PathPrefix("/swagger").Handler(sh)
}

func (rtr Router) health() {
	rtr.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Up"))
	}).Methods("GET")
}
