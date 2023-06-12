package request

type ProductRequest struct {
	CategoryID  uint    `gorm:"category_id" form:"category_id" json:"category_id" validate:"required"`
	Name        string  `gorm:"name" form:"name" validate:"required"`
	Price       float32 `gorm:"price" form:"price" validate:"required,gt=0"`
	Description string  `gorm:"description" form:"description" validate:"required"`
	Status      int     `gorm:"status" form:"status" json:"status" validate:"required,arrayIn=1&2"`
}

type ProductUpdateRequest struct {
	CategoryID   uint     `gorm:"category_id" form:"category_id" json:"category_id" validate:"required"`
	Name         string   `gorm:"name" form:"name" validate:"required"`
	Price        float32  `gorm:"price" form:"price" validate:"required,gt=0"`
	Description  string   `gorm:"description" form:"description" validate:"required"`
	Status       int      `gorm:"status" form:"status" json:"status" validate:"required,arrayIn=1&2"`
	ImageRemoves []string `gorm:"image_removes" form:"image_removes"`
}

type ProductUpdateStatusRequest struct {
	Status int `gorm:"status" form:"status" json:"status" validate:"required,arrayIn=1&2"`
}
