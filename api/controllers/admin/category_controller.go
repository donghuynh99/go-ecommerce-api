package admin

import (
	"net/http"
	"strconv"

	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"
)

func (controller *AdminController) ListCategories(c *gin.Context) {
	var categories []models.Category

	page, limit, offset, err := controller.GetDataPagination(c)

	var totalCount int64
	controller.db.Preload(clause.Associations).Offset(offset).Limit(limit).Find(&categories)
	controller.db.Count(&totalCount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	result := config.PaginationStruct{
		Data:       categories,
		TotalCount: int(totalCount),
		Page:       page,
	}

	c.JSON(http.StatusOK, result)
}

func (controller *AdminController) StoreCategories(c *gin.Context) {
	var request request.CategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
			"error":   err.Error(),
		})

		return
	}

	errors := utils.ValidateStruct(request)

	if errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}

	if request.ParentID != nil {
		var parentCategory models.Category
		checkParent := controller.db.First(&parentCategory, "id = ?", request.ParentID)
		if checkParent.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
				Param:   "parent_id",
				Message: "Not match with any category!",
			}})
			return
		}
	}

	var categoryCheck models.Category
	err := controller.CheckExisted(&categoryCheck, map[string]interface{}{
		"name": request.Name,
	}, nil)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "name",
			Message: "Already existed!",
		}})
		return
	}

	category := models.Category{
		ParentID: request.ParentID,
		Name:     request.Name,
		Slug:     slug.Make(request.Name),
	}

	errStore := controller.db.Create(&category)

	if errStore.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Create failure!",
			"error":   errStore.Error,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Create category successful!",
		"category": category,
	})

	return
}

func (controller *AdminController) UpdateCategory(c *gin.Context) {
	var request request.CategoryRequest
	var categoryID = c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
			"error":   err.Error(),
		})

		return
	}

	errors := utils.ValidateStruct(request)

	if errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}

	var category models.Category
	err := controller.CheckExisted(&category, map[string]interface{}{
		"id": categoryID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found!",
		})

		return
	}

	var categoryCheck models.Category
	err = controller.CheckExisted(&categoryCheck, map[string]interface{}{
		"name":    request.Name,
		"id != ?": category.ID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "name",
			Message: "Already existed!",
		}})
		return
	}

	if request.ParentID != nil {
		var categoryCheck models.Category
		err := controller.CheckExisted(&categoryCheck, map[string]interface{}{
			"id": request.ParentID,
		}, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
				Param:   "parent_id",
				Message: "Not match with any category!",
			}})
			return
		}
	}

	category.Name = request.Name
	category.Slug = slug.Make(request.Name)
	category.ParentID = request.ParentID

	result := controller.db.Save(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Update fail!",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update successful!",
	})

	return
}

func (controller *AdminController) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category
	err := controller.CheckExisted(&category, map[string]interface{}{
		"id": categoryID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found!",
		})

		return
	}

	childs, err := category.GetChildren()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error!",
		})

		return
	}

	if len(childs) == 0 {
		result := controller.db.Select("Products").Delete(&category)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Delete fail!",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Delete successful!",
		})

		return
	}

	url := c.Request.Host + "/api/admin/categories/" + strconv.FormatUint(uint64(category.ID), 10) + "/force-delete"
	keySignature := "force-delete-" + strconv.FormatUint(uint64(category.ID), 10)
	signatureDelete := utils.CreateSignature(keySignature)

	c.JSON(http.StatusOK, gin.H{
		"message": "This category have children. Please confirm to remove this category!",
		"url":     url + "?signature=" + signatureDelete,
	})

	return
}

func (controller *AdminController) ForceDeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	signature := c.Query("signature")
	keyString := "force-delete-" + categoryID

	var category models.Category

	err := controller.CheckExisted(&category, map[string]interface{}{
		"id": categoryID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Category not found!",
		})

		return
	}

	if !utils.VerifySignature(keyString, signature) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid signature",
		})

		return
	}

	childs, err := category.GetChildren()

	controller.db.Select("Products").Delete(&childs)

	c.JSON(http.StatusOK, gin.H{
		"message": "Force delete successful!",
	})

	return
}
