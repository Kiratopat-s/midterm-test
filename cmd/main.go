package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kiratopat-s/workflow/internal/auth"
	"github.com/Kiratopat-s/workflow/internal/item"
	"github.com/Kiratopat-s/workflow/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// POST 	/items | req.body {"title": "ขอเบิกสินค้า", "amount": 12, "quantity": 10} => res(201) {"id": 1,"amount": 12,"quantity": 10,"status": "PENDING","owner_id": 1 }
// GET 		/items | => res(200) [{"id": 1,"amount": 12,"quantity": 10,"status": "PENDING","owner_id": 1 },... all items]
// GET 		/items/:id | req.body {"title": "ขอเบิกสินค้า", "amount": 12, "quantity": 10} => res(200) {"id": 1,"title": "ขอเบิกสินค้า", "amount": 12, "quantity": 10, "status": "PENDING", "owner_id": 1}
// PUT 		/items/:id | req.body {"title": 1,"amount": 100,"quantity": 20} => res(200) {"id": 1,"title": "ขอเบิกสินค้า", "amount": 100, "quantity": 20, "status": "PENDING", "owner_id": 1}
// PATCH	/items/:id | req.body {"status": "APPROVED"} => res(200) {"id": 1,"title": "ขอเบิกสินค้า", "amount": 12, "quantity": 10, "status": "APPROVED", "owner_id": 1}
// DELETE 	/items/:id | res(204)
// POST		/login | req.body {"username": "admin", "password": "secret"} => set Cookie token=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiYXVkIjoiU3VwZXIgQWRtaW4iLCJleHAiOjE res(200) {"message": "Login Succeed"}
// POST		/logout | res(200) {"message": "Logout Succeed"}

func main() {
	// Connect database
	db, err := gorm.Open(
		postgres.Open(
			os.Getenv("DATABASE_URL"),
			// "postgres://postgres:postgres@localhost:5432/iws",
		),
	)
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
	
	// Router
	r := gin.Default()

	config := cors.DefaultConfig()
	// frontend URL
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	r.Use(cors.New(config))

	r.GET("/version", func(c *gin.Context) {
		version, err := GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	// Register router
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	r.GET("/hello-verifytoken",verifyToken, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!",
			"uid": c.MustGet("uid"),
			"username": c.MustGet("username"),
			"firstName": c.MustGet("firstName"),
			"lastName": c.MustGet("lastName"),
			"position": c.MustGet("position"),
			"photoLink": c.MustGet("photoLink"),},
		)
	})
	r.POST("/items",verifyToken, controller.CreateItem)
	r.GET("/items",verifyToken, controller.FindAllItem)
	r.GET("/items/:id",verifyToken, controller.FindItemByID)
	r.PUT("/items/:id",verifyToken, controller.UpdateItem)
	r.PATCH("/items/:id",verifyToken, controller.UpdateItemStatus)
	r.DELETE("/items/:id",verifyToken, controller.DeleteItem)
	r.POST("/login", userController.Login)
	r.POST(("register"), userController.Register)

	r.POST("/logout",verifyToken, userController.Logout)
	
	


	// Start server
	// Get the port from the environment variable or default to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default port
    }

    // Run the server on the specified port
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

type GooseDBVersion struct {
	ID        int
	VersionID int
	IsApplied bool
	Tstamp    string
}

// TableName overrides the table name used by User to `profiles`
func (GooseDBVersion) TableName() string {
	return "goose_db_version"
}

// GetLatestDBVersion returns the latest applied version from the goose_db_version table.
func GetLatestDBVersion(db *gorm.DB) (int, error) {
	var version GooseDBVersion

	// Query to get the latest version applied
	err := db.Order("version_id desc").Where("is_applied = ?", true).First(&version).Error
	if err != nil {
		return 0, err
	}

	return version.VersionID, nil
}
