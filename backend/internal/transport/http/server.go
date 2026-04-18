package srv

import (
	"ethno/internal/auth"
	"ethno/internal/config"
	"ethno/internal/repository"
	handler "ethno/internal/transport/http/handlers"
	"ethno/internal/transport/http/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *chi.Mux
	logger *logrus.Logger
}

func NewServer(
	folkRepo *repository.FolkRepository,
	authService *auth.AuthService,
	cfg *config.ServerConfig,
	logger *logrus.Logger,
	questRepo *repository.QuestRepository,
) *Server {
	authHandler := handler.NewAuthHandler(authService, cfg)
	questHandler := handler.NewQuestHandler(questRepo, logger)

	r := chi.NewRouter()
	r.Use(loggerMiddleware(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)
	r.Get("/api/regions", handler.GetRandomFolksHandler(folkRepo))
	
	questHandler.RegisterRoutes(r) 

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthFromCookie(authService.AuthProv, logger, "auth_token"))
		r.Get("/api/me", authHandler.GetMe)
	})

	staticFS := http.FileServer(http.Dir("frontend/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", staticFS))

	r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/profile.html")
	})
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/login.html")
	})
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/register.html")
	})
	r.Get("/quests/{slug}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/quest.html")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/index.html")
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