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
				Message: utils.Translation("not_found", nil, nil),
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
			Message: utils.Translation("already_existed", nil, nil),
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
			"message": utils.Translation("create_fail", nil, nil),
			"error":   errStore.Error,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  utils.Translation("create_success", nil, nil),
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
			"message": utils.Translation("not_found", nil, nil),
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
			Message: utils.Translation("already_existed", nil, nil),
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
				Message: utils.Translation("not_found", nil, nil),
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
			"message": utils.Translation("update_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("update_success", nil, nil),
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
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	childs, err := category.GetChildren()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})

		return
	}

	if len(childs) == 0 {
		result := controller.db.Select("Products").Delete(&category)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": utils.Translation("delete_fail", nil, nil),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": utils.Translation("delete_success", nil, nil),
		})

		return
	}

	url := c.Request.Host + "/api/admin/categories/" + strconv.FormatUint(uint64(category.ID), 10) + "/force-delete"
	keySignature := "force-delete-" + strconv.FormatUint(uint64(category.ID), 10)
	signatureDelete := utils.CreateSignature(keySignature)

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("already_had_children", nil, nil),
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
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	if !utils.VerifySignature(keyString, signature) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("invalid_signature", nil, nil),
		})

		return
	}

	childs, err := category.GetChildren()

	controller.db.Select("Products").Delete(&childs)

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("delete_success", nil, nil),
	})

	return
}
