package request

type AddressRequest struct {
	Name      string `gorm:"name" json:"name" validate:"required"`
	IsPrimary bool   `gorm:"is_primary" json:"is_primary" validate:"boolean"`
	PostCode  string `gorm:"post_code" json:"post_code" validate:"required"`
}
