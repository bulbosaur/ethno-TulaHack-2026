package srv

import (
	"ethno/internal/auth"
	"ethno/internal/config"
	"ethno/internal/repository"
	handler "ethno/internal/transport/http/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
    router *chi.Mux
    logger *logrus.Logger
}

func NewServer(folkRepo *repository.FolkRepository, authService *auth.AuthService, cfg *config.ServerConfig, logger *logrus.Logger) *Server {
	authHandler := handler.NewAuthHandler(authService, cfg)
    r := chi.NewRouter()
    r.Use(loggerMiddleware(logger))

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    r.Post("/api/register", authHandler.Register)
	r.Get("/api/regions", handler.GetRandomFolksHandler(folkRepo))

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "../frontend/templates/login.html")
    })
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "../frontend/templates/register.html")
    })

	fs := http.FileServer(http.Dir("../frontend/templates"))
    r.Handle("/*", http.StripPrefix("/", fs))
	
    staticFS := http.FileServer(http.Dir("../frontend/static"))
    r.Handle("/static/*", http.StripPrefix("/static/", staticFS))

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
