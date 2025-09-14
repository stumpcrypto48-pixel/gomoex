package main

import (
	"fmt"
	"httpfromtcp/rootmod/internal/router"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func NewApp(dsn *string) (*App, error) {
	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error occured while try to connect to DB :: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Error occured while try to check DB connection :: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("Error occured while ping DB :: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	app := &App{
		DB:     db,
		Router: gin.Default(),
	}

	app.Router.Use(gin.Logger(), gin.Recovery())

	return app, nil
}

func (a *App) registerRoutes() {
	v1 := a.Router.Group("/v1/moex")
	router.MoexRouter(a.DB, v1)
	a.Router.Run()
}

func (a *App) health(c *gin.Context) {
	sqlDB, err := a.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "db internal error"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "db unreachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while reading dotenv file :: %v", err)
	}
	dsn := os.Getenv("GORM_POSTGRES")

	app, err := NewApp(&dsn)
	if err != nil {
		log.Fatalf("Error while init app :: %v", err)
	}
	app.registerRoutes()

}
