package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/Kiratopat-s/workflow/internal/auth"
	"github.com/Kiratopat-s/workflow/internal/item"
	"github.com/Kiratopat-s/workflow/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect database
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Panic(err)
	}

	// Controller
	controller := item.NewController(db)
	userController := user.NewController(db, "secret")

	// verifyToken middleware
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	verifyToken := auth.Guard(secret)
	verifyAdmin := auth.GuardAdmin(secret)

	// Router setup
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/version", func(c *gin.Context) {
		version, err := GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	// Routes
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	r.GET("/hello-verifytoken", verifyToken, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "Hello, World!",
			"uid":       c.MustGet("uid"),
			"username":  c.MustGet("username"),
			"firstName": c.MustGet("firstName"),
			"lastName":  c.MustGet("lastName"),
			"position":  c.MustGet("position"),
			"photoLink": c.MustGet("photoLink"),
		})
	})
	r.POST("/items", verifyToken, controller.CreateItem)
	r.GET("/items", verifyToken, controller.FindAllItem)
	r.GET("/items/:id", verifyToken, controller.FindItemByID)
	r.PUT("/items/:id", verifyToken, controller.UpdateItem)
	r.PATCH("/items/:id", verifyToken, controller.UpdateItemStatus)
	r.PATCH("/items/update/status/many", verifyAdmin, controller.UpdateManyItemsStatus)
	r.DELETE("/items/:id", verifyToken, controller.DeleteItem)
	r.DELETE("/items/delete/many", verifyToken, controller.DeleteManyItems)
	r.GET("/items/status/count/user", verifyToken, controller.CountItemsStatusByUser)
	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)
	// r.POST("/logout", verifyToken, userController.Logout)

	// Graceful shutdown setup
	srv := &http.Server{
		Addr:    ":" + getPort(),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for the server to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	return port
}

type GooseDBVersion struct {
	ID        int
	VersionID int
	IsApplied bool
	Tstamp    string
}

// TableName overrides the table name used by GooseDBVersion
func (GooseDBVersion) TableName() string {
	return "goose_db_version"
}

// GetLatestDBVersion returns the latest applied version from the goose_db_version table.
func GetLatestDBVersion(db *gorm.DB) (int, error) {
	var version GooseDBVersion
	err := db.Order("version_id desc").Where("is_applied = ?", true).First(&version).Error
	if err != nil {
		return 0, err
	}
	return version.VersionID, nil
}