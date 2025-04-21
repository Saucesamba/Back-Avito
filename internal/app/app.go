package app

import (
	"Backend_trainee_assigment_2025/internal/auth"
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"Backend_trainee_assigment_2025/internal/handlers"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App struct {
	Router     *mux.Router
	Config     *config.AppConfig
	DB         db.Database
	HTTPServer *http.Server
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize() error {
	cfg, err := new(config.AppConfig).LoadConfig()
	if err != nil {
		return err
	}
	app.Config = cfg
	log.Println("Waiting for database to start...")
	time.Sleep(10 * time.Second)

	avitoDB, err := new(db.AvitoDB).NewDB(&config.DBConfig{
		Host:     app.Config.Database.Host,
		Port:     app.Config.Database.Port,
		User:     app.Config.Database.User,
		Password: app.Config.Database.Password,
		Name:     app.Config.Database.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err) // Added error context
	}

	app.DB = avitoDB

	app.Router = mux.NewRouter()
	app.setRouters()
	app.HTTPServer = &http.Server{
		Addr:    app.Config.Server.Host + ":" + app.Config.Server.Port,
		Handler: app.Router,
	}
	return nil
}

func (app *App) setRouters() {
	avitoDB, ok := app.DB.(*db.AvitoDB)
	if !ok {
		log.Fatalf("expected *db.AvitoDB, got %T", app.DB)
	}

	userHandler := handlers.NewUserHandler(avitoDB, app.Config)
	pvzHandler := handlers.NewPVZHandler(avitoDB, app.Config)
	receptionHandler := handlers.NewReceptionHandler(avitoDB, app.Config)
	productHandler := handlers.NewProductHandler(avitoDB, app.Config)

	secretKey := app.Config.JWT.Secret

	// Define Routes
	app.Router.HandleFunc("/dummyLogin", userHandler.DummyLoginHandler).Methods("POST")
	app.Router.HandleFunc("/register", userHandler.RegisterHandler).Methods("POST")
	app.Router.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")

	app.Router.HandleFunc("/pvz", authMiddleware(pvzHandler.CreatePVZHandler, secretKey, "moderator")).Methods("POST")
	app.Router.HandleFunc("/pvz", authMiddleware(pvzHandler.GetPVZsHandler, secretKey, "employee", "moderator")).Methods("GET")
	app.Router.HandleFunc("/pvz/{pvzId}/close_last_reception", authMiddleware(receptionHandler.CloseLastReceptionHandler, secretKey, "employee")).Methods("POST")
	app.Router.HandleFunc("/pvz/{pvzId}/delete_last_product", authMiddleware(productHandler.DeleteLastProductHandler, secretKey, "employee")).Methods("POST")

	// Reception Routes
	app.Router.HandleFunc("/receptions", authMiddleware(receptionHandler.CreateReceptionHandler, secretKey, "employee")).Methods("POST")
	// Product Routes
	app.Router.HandleFunc("/products", authMiddleware(productHandler.AddProductHandler, secretKey, "employee")).Methods("POST")
}

func (app *App) Start() error {
	go func() {
		log.Printf("Server listening on %s", app.HTTPServer.Addr)
		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", app.HTTPServer.Addr, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	return nil
}

func authMiddleware(next http.HandlerFunc, secretKey string, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}
		tokenString = parts[1]

		claims, err := auth.ValidateToken(tokenString, secretKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		if !checkRole(claims.Role, roles) {
			respondWithError(w, http.StatusForbidden, "Unauthorized: insufficient role")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func checkRole(userRole string, roles []string) bool {

	for _, role := range roles {
		if userRole == role {
			return true
		}
	}
	return false
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
