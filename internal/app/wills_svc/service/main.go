package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ctxKey int8

const (
	sessionName        = "auths"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	db     *sqlx.DB
}

func Start(cfg *Config) error {
	srv := newServer(cfg.DatabaseURL)
	return http.ListenAndServe(cfg.ListenerAddr, srv)
}

func newServer(databaseURL string) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
	s.configureStore(databaseURL)
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
}

func (s *server) configureStore(databaseURL string) {
	db, err := sqlx.Connect("postgres", databaseURL)
	defer func() {
		if err := db.Close(); err != nil {
			s.logger.Fatal(err)
		}
	}()
	if err != nil {
		s.logger.Fatal(err)
	}
	s.db = db
}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("Hello world!")
	}
}
