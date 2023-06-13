package admin

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/donghuynh99/ecommerce_api/api/request"
	"github.com/donghuynh99/ecommerce_api/config"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/donghuynh99/ecommerce_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gosimple/slug"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func (controller *AdminController) ListProduct(c *gin.Context) {
	var products []models.Product

	page, limit, offset, err := controller.GetDataPagination(c)

	var totalCount int64
	controller.db.Preload(clause.Associations).Offset(offset).Limit(limit).Find(&products)
	controller.db.Count(&totalCount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	result := config.PaginationStruct{
		Data:       products,
		TotalCount: int(totalCount),
		Page:       page,
	}

	c.JSON(http.StatusOK, result)
}

func (controller *AdminController) ShowProduct(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product
	err := controller.CheckExisted(&product, map[string]interface{}{
		"id": productID,
	}, &[]string{"Images", "Categories"})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})

	return
}

func (controller *AdminController) StoreProduct(c *gin.Context) {
	var request request.ProductRequest

	if err := c.MustBindWith(&request, binding.Form); err != nil {
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

	form, errForm := c.MultipartForm()
	if errForm != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
			"error":   errForm.Error(),
		})

		return
	}

	images := form.File["images"]

	if len(images) > config.GetConfig().GeneralConfig.ImageLimit {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param: "images",
			Message: utils.Translation("image_limit", map[string]interface{}{
				"Number": 5,
			}, nil),
		}})
		return
	}

	var category models.Category
	err := controller.CheckExisted(&category, map[string]interface{}{
		"id": request.CategoryID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "category_id",
			Message: utils.Translation("not_found", nil, nil),
		}})
		return
	}

	var checkNameExisted models.Product
	isNameExisted := controller.CheckExisted(&checkNameExisted, map[string]interface{}{
		"name": request.Name,
	}, nil)

	if isNameExisted == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "name",
			Message: utils.Translation("already_existed", nil, nil),
		}})
		return
	}

	parents := category.GetParents()

	imagesModels := []models.Image{}

	for _, file := range images {
		destination := config.GetConfig().GeneralConfig.DestinationStoreProductImage

		imagesName := utils.GenerateUUID() + file.Filename

		os.Chmod(destination, 0755)

		c.SaveUploadedFile(file, destination+imagesName)

		imagesModels = append(imagesModels, models.Image{
			Path: destination + imagesName,
			Name: imagesName,
			Alt:  slug.Make(request.Description),
		})
	}

	product := models.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       decimal.NewFromFloat32(request.Price),
		Status:      request.Status,
		Categories:  parents,
		Slug:        slug.Make(request.Name),
		Images:      imagesModels,
	}

	errSave := controller.db.Save(&product)

	if errSave.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": utils.Translation("create_fail", nil, nil),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.Translation("create_success", nil, nil),
		"data":    product,
	})

	return
}

func (controller *AdminController) UpdateProduct(c *gin.Context) {
	// abc := controller.db.Delete(&[]models.Image{models.Image{ID: 8}, models.Image{ID: 9}}).Error
	var request request.ProductUpdateRequest

	var productID = c.Param("id")

	if err := c.MustBindWith(&request, binding.Form); err != nil {
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

	var product models.Product
	err := controller.CheckExisted(&product, map[string]interface{}{
		"id": productID,
	}, &[]string{"Images"})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	var productCheck models.Product
	err = controller.CheckExisted(&productCheck, map[string]interface{}{
		"name":    request.Name,
		"id != ?": product.ID,
	}, nil)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "name",
			Message: utils.Translation("already_existed", nil, nil),
		}})
		return
	}

	var category models.Category
	err = controller.CheckExisted(&category, map[string]interface{}{
		"id": request.CategoryID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param:   "category_id",
			Message: utils.Translation("not_found", nil, nil),
		}})
		return
	}

	form, errForm := c.MultipartForm()
	if errForm != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.Translation("bad_request", nil, nil),
			"error":   errForm.Error(),
		})

		return
	}

	images := form.File["images"]

	product.Name = request.Name
	product.Price = decimal.NewFromFloat32(request.Price)
	product.Slug = slug.Make(request.Name)
	product.Description = request.Description
	product.Status = request.Status

	parents := category.GetParents()

	controller.db.Model(&product).Association("Categories").Replace(&parents)

	imageRemoves := make([]models.Image, len(request.ImageRemoves))

	if len(imageRemoves) > 0 {
		controller.db.Model(&product).Where("name IN ?", request.ImageRemoves).Association("Images").Find(&imageRemoves)
	}

	if len(product.Images)+len(images)-len(imageRemoves) > config.GetConfig().GeneralConfig.ImageLimit {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
			Param: "images",
			Message: utils.Translation("image_limit", map[string]interface{}{
				"Number": strconv.Itoa(5 + len(request.ImageRemoves) - len(product.Images)),
			}, nil),
		}})
		return
	}

	imagesModels := []models.Image{}
	for _, file := range images {
		if !utils.IsImageFile(file) {
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.ApiError{
				Param:   "images",
				Message: utils.Translation("image_type_invalid", nil, nil),
			}})
			return
		}

		destination := "assets/products/images/"

		imagesName := utils.GenerateUUID() + file.Filename

		os.Chmod(destination, 0755)

		c.SaveUploadedFile(file, destination+imagesName)

		imagesModels = append(imagesModels, models.Image{
			Path: destination + imagesName,
			Name: imagesName,
			Alt:  slug.Make(request.Description),
		})
	}

	controller.db.Model(&product).Association("Images").Append(&imagesModels)

	result := controller.db.Save(&product)

	if len(imageRemoves) > 0 {
		controller.db.Delete(&imageRemoves)
	}

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

func (controller *AdminController) UpdateStatusProduct(c *gin.Context) {
	var request request.ProductUpdateStatusRequest
	productID := c.Param("id")

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

	var product models.Product
	err := controller.CheckExisted(&product, map[string]interface{}{
		"id": productID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	product.Status = request.Status
	result := controller.db.Save(&product)

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

func (controller *AdminController) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product
	err := controller.CheckExisted(&product, map[string]interface{}{
		"id": productID,
	}, nil)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": utils.Translation("not_found", nil, nil),
		})

		return
	}

	var images []models.Image
	controller.db.Model(&product).Association("Images").Find(&images)

	errImageDelete := controller.db.Delete(&images)

	result := controller.db.Select("Categories").Delete(&product)

	if result.Error != nil || errImageDelete.Error != nil {
		log.Panicln(result.Error, errImageDelete.Error)
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
