package srv

import (
	"ethno/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
    router *chi.Mux
    logger *logrus.Logger
}

func NewServer(cfg *config.ServerConfig, logger *logrus.Logger) *Server {
    r := chi.NewRouter()
    r.Use(loggerMiddleware(logger))

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })
	
    return &Server{
        router: r,
        logger: logger,
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.router.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
    s.logger.Infof("HTTP server starting on %s", addr)
    return http.ListenAndServe(addr, s.router)
}

func loggerMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Infof("%s %s", r.Method, r.URL.Path)
            next.ServeHTTP(w, r)
        })
    }
}
