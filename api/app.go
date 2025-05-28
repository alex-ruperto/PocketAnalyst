package api

import (
	"PocketAnalyst/clients"
	"PocketAnalyst/controllers"
	"PocketAnalyst/repositories"
	"PocketAnalyst/services"
	"database/sql"
	"log"
	"net/http"
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
	AlphaVantageAPIKey    string
	AlphaVantageBaseURL   string
	Port                  string
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifeTime time.Duration
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
	db.SetConnMaxLifetime(app.Config.ConnectionMaxLifeTime)

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
	datasourceRepo := repositories.NewDataSourceRepository(app.DB)

	// Initialize external clients
	alphaClient := clients.NewAlphaVantageClient(
		app.Config.AlphaVantageBaseURL,
		app.Config.AlphaVantageAPIKey,
	)

	// Initialize services
	stockService := services.NewStockService(stockRepo, alphaClient)

	// Initialize controllers
	stockController := controllers.NewStockController(stockService)

	// Register routes with middleware
	app.Router.HandleFunc("/api/stocks/fetch", app.withMiddleware(stockController.HandleStockFetchRequest))
	app.Router.HandleFunc("/api/stocks/get", app.withMiddleware(stockController.HandleStockHistoryRequest))

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
