package app

import (
	"database/sql"
	"log"
	"net/http"
	"pocketanalyst/internal/controllers"
	"pocketanalyst/internal/repositories"
	"pocketanalyst/internal/services"
	"pocketanalyst/pkg/clients"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// App represents the main application with all its dependencies
type App struct {
	DB     *sql.DB
	Router *http.ServeMux
	Config *Config
}

// Config holds all application configuration
type Config struct {
	DatabaseURL           string
	FMPAPIKey             string
	FMPBaseURL            string
	Port                  string
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifetime time.Duration
}

// NewApp creates a new app instance.
func NewApp(config *Config) (*App, error) {
	app := &App{
		Config: config,
		Router: http.NewServeMux(),
	}

	// Initialize DB connection
	if err := app.initDatabase(); err != nil {
		return nil, err
	}

	// Setup routes and dependencies
	if err := app.setupRoutes(); err != nil {
		return nil, err
	}

	return app, nil
}

// initDatabase establishse a db connection
func (app *App) initDatabase() error {
	db, err := sql.Open("postgres", app.Config.DatabaseURL)
	if err != nil {
		return err
	}

	// Configure connection pool for prod use
	db.SetMaxIdleConns(app.Config.MaxIdleConnections)
	db.SetMaxOpenConns(app.Config.MaxOpenConnections)
	db.SetConnMaxLifetime(app.Config.ConnectionMaxLifetime)

	// Verify db connection
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	app.DB = db
	log.Println("Database connection established successfully")
	return nil
}

// setupRoutes initalizes all application dependencies and sets up HTTP routes
func (app *App) setupRoutes() error {
	// Initalize repositories
	stockRepo := repositories.NewStockRepository(app.DB)

	// Initialize client factory and register providers
	factory := clients.NewClientFactory()
	factory.RegisterProvider("fmp", app.Config.FMPBaseURL, app.Config.FMPAPIKey)

	// Create FMP Client using factory
	client, err := factory.CreateClient("fmp")
	if err != nil {
		return err
	}

	// Initialize services
	stockService := services.NewStockService(stockRepo, client)

	// Initialize controllers
	stockController := controllers.NewStockController(stockService)

	// Register route endpoints with middleware
	app.Router.HandleFunc("/api/stocks/fetch", app.withMiddleware(stockController.HandleStockFetchRequest))
	app.Router.HandleFunc("/api/stocks/get-stock", app.withMiddleware(stockController.HandleStockHistoryRequest))
	app.Router.HandleFunc("/api/stocks/health", app.withMiddleware(stockController.HandleHealthCheckRequest))
	app.Router.HandleFunc("/api/stocks/get-distinct-symbols", app.withMiddleware(stockController.HandleGetDistinctSymbolRequest))

	log.Println("Routes configured successfully")
	return nil
}

// withMiddleware applies common middleware to all routes
func (app *App) withMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply logging middleware
		app.loggingMiddleware(
			app.corsMiddleware(
				app.recoveryMiddleware(handler),
			),
		)(w, r)
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture status codes
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs all incoming requests
func (app *App) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		log.Printf("%s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	}
}

// corsMiddleWare handles Cross-Origin resource sharing
func (app *App) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// recoveryMiddleware handles recovery from panics and returns a 500 error.
func (app *App) recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

// Start the HTTP server with the configured timeouts
func (app *App) Start() error {
	server := &http.Server{
		Addr:         ":" + app.Config.Port,
		Handler:      app.Router,
		ReadTimeout:  app.Config.ReadTimeout,
		WriteTimeout: app.Config.WriteTimeout,
	}

	log.Printf("Starting server on port %s", app.Config.Port)
	return server.ListenAndServe()
}

// Gracefully shuts down the application
func (app *App) Close() error {
	if app.DB != nil {
		log.Println("Closing database connection")
		return app.DB.Close()
	}
	return nil
}
