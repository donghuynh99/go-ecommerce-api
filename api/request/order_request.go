package request

type OrderRequest struct {
	Note string `gorm:"note"`
}

type OrderCancelRequest struct {
	Note string `gorm:"note"`
}
