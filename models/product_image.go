package models

type ProductImage struct {
	ID        string `gorm:"size:100;not null;auto_increment,uniqueIndex;primary_key"`
	Product   Product
	ProductID string `gorm:"size:36;index"`
	Path      string `gorm:"type:text"`
	Alt       string `gorm:"type:text"`
}
