package item

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kiratopat-s/workflow/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB) Controller {
	return Controller{
		Service: NewService(db),
	}
}

type ApiError struct {
	Field  string
	Reason string
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "จำเป็นต้องกรอกข้อมูลนี้"
	case "email":
		return "Invalid email"
	case "gt":
		return fmt.Sprintf("Number must greater than %v", param)
	case "gte":
		return fmt.Sprintf("Number must greater than or equal %v", param)
	}
	return ""
}

func getValidationErrors(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Param())}
		}
		return out
	}
	return nil
}

func (controller Controller) CreateItem(ctx *gin.Context) {
	// Bind
	var request model.RequestCreateItem

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": getValidationErrors(err),
		})
		return
	}

	// Create item
	// get owner_id from context
	ownerIdFloat := ctx.MustGet("uid").(float64)
	ownerId := int(ownerIdFloat)
	item, err := controller.Service.Create(request, ownerId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"data": item,
	})
}


func (controller Controller) FindAllItem(ctx *gin.Context) {

	
	items, err := controller.Service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (controller Controller) FindItemByID(ctx *gin.Context) {
	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Find item
	item, err := controller.Service.FindByID(uint(id))
	if err != nil {
		// Check if the error is because the record was not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Item not found",
			})
			return
		}

		// Other types of errors (e.g., database issues)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error finding item: " + err.Error(),
		})
		return
	}

	// Successfully found item
	ctx.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}


func (controller Controller) UpdateItem(ctx *gin.Context) {
	// Bind
	var (
		request model.RequestUpdateItem
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Update item
	item, err := controller.Service.UpdateItem(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func (controller Controller) UpdateItemStatus(ctx *gin.Context) {
		// Bind
		var (
			request model.RequestPatchItemStatus
		)
	
		if err := ctx.Bind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}
	
		// Path param
		id, _ := strconv.Atoi(ctx.Param("id"))
	
		// Update status
		item, err := controller.Service.UpdateStatus(uint(id), request.Status)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{
			"data": item,
		})
	}

func (controller Controller) DeleteItem(ctx *gin.Context) {
	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Delete
	if err := controller.Service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}



func (controller Controller) UpdateManyItemsStatus(ctx *gin.Context) {
	// get ids : int[], status: string from body to do next
	var request model.RequestPatchManyItemStatus
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Update status
	err := controller.Service.UpdateManyStatus(request.IDs, request.Status )
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": "Updated",
	})
}



func (controller Controller) DeleteManyItems(ctx *gin.Context) {
	// get ids : int[] from body to do next
	var request model.RequestDeleteManyItems
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// get position and uid from context
	ownerIdFloat := ctx.MustGet("uid").(float64)
	ownerId := int(ownerIdFloat)
	userPostion := ctx.MustGet("position").(string)


	// Delete
	err := controller.Service.DeleteMany(request.IDs,ownerId,userPostion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}


func (controller Controller) CountItemsStatusByUser(ctx *gin.Context) {
	// get owner_id from context
	ownerIdFloat := ctx.MustGet("uid").(float64)
	ownerId := int(ownerIdFloat)

	// Count
	counts, err := controller.Service.CountItemsStatusByUser(ownerId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": counts,
	})
}
