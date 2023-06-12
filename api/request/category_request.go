package request

type CategoryRequest struct {
	ParentID *uint  `gorm:"parent_id" json:"parent_id"`
	Name     string `gorm:"name" validate:"required"`
}
