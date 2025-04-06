package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	dionysusMiddleware "github.com/jimvid/dionysus/internal/middleware"
	"github.com/jimvid/dionysus/internal/user"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(dionysusMiddleware.ClaimsContextKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized - no claims", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "claims: %+v", claims)
}

func (s *Server) RegisterRoutes() http.Handler {
	dbInstance := s.db.GetDBInstance()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // TODO
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health
	r.Get("/health", s.healthHandler)

	// User
	userService := user.NewUserService(dbInstance)
	userHandler := user.NewUserHandler(userService)
	r.Post("/user/register", userHandler.RegisterUserHandler)
	r.Post("/user/login", userHandler.LoginUserHandler)

	// Protected routes
	r.With(dionysusMiddleware.ValidateJWTMiddleware).Get("/protected", ProtectedHandler)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
