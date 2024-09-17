package user

import (
	"fmt"
	"net/http"

	"github.com/Kiratopat-s/workflow/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB, secret string) Controller {
	return Controller{
		Service: NewService(db, secret),
	}
}

func (controller Controller) Login(ctx *gin.Context) {
	var (
		request model.RequestLogin
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := controller.Service.Login(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login succeed",
		"token": "Bearer "+token,
	})
}

// func (controller Controller) Logout(ctx *gin.Context) {
// 	token := ctx.GetHeader("Authorization")
// 	if token == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"message": "token is required",
// 		})
// 		return
// 	}

// 	err := controller.Service.Logout(token)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "logout succeed",
// 	})
// }

func (controller Controller) Register(ctx *gin.Context) {
	var (
		request model.RequestRegister
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Println("CTR || ",request)

	err := controller.Service.Register(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "register succeed",
	})
}