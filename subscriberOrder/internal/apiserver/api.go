package apiserver

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"awesomeProject5/subscriberOrder/internal/models"
	"awesomeProject5/subscriberOrder/internal/store"
)

type APIServer struct {
	Router *chi.Mux
	store  *store.Store
	redis  *store.Cache
}

func NewAPIServer(store *store.Store, redis *store.Cache) *APIServer {
	return &APIServer{
		Router: chi.NewRouter(),
		store:  store,
		redis:  redis,
	}

}

func (s *APIServer) Init() error {
	err := s.configureRouter()
	if err != nil {
		return err
	}
	err = s.configureStore()
	if err != nil {
		return err
	}
	err = s.configureCache()
	if err != nil {
		return err
	}

	return nil

}

func (s *APIServer) configureStore() error {
	if err := s.store.OpenDB(); err != nil {
		return err
	}
	return nil
}

func (s *APIServer) configureCache() error {
	if err := s.redis.InitCache(); err != nil {
		return err
	}
	return nil
}
func (s *APIServer) configureRouter() error {

	r := s.Router
	// Serve static files with appropriate MIME types
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Serve the index.html file
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		http.ServeFile(writer, request, "static/index.html")
	})

	r.Get("/getOrders", s.CacheMiddleware(s.GetOrdersByID()))
	return nil
}

func (s *APIServer) GetOrdersByID() http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		a, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		order, err := s.store.GetByID(a)
		if err != nil || order.OrderUid == "" {
			http.Error(writer, "order does not exist", http.StatusNotFound)
			return
		}

		readyToCache, err := json.Marshal(order)

		if err := s.redis.SaveToCache(readyToCache, id); err != nil {
			http.Error(writer, "error", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("static/details.html")
		if err != nil {
			http.Error(writer, "Failed to load template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(writer, order)
		if err != nil {
			http.Error(writer, "Failed to render template", http.StatusInternalServerError)
			return
		}

	}
}

func (s *APIServer) CacheMiddleware(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		orderId := request.URL.Query().Get("id")
		order, found := s.redis.GetOrderByID(orderId)

		orderRequest := &models.Order{}

		if err := json.Unmarshal(order, orderRequest); err != nil {
			http.Error(writer, "order not found", http.StatusInternalServerError)
			return
		}

		if found {
			tmpl, err := template.ParseFiles("static/details.html")
			if err != nil {
				http.Error(writer, "Failed to load template", http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(writer, orderRequest)
			if err != nil {
				http.Error(writer, "Failed to render template", http.StatusInternalServerError)
				return
			}
		} else {
			next.ServeHTTP(writer, request)
		}

	}

}
