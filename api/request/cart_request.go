package request

type CartRequest struct {
	Data map[string]int `gorm:"data"`
}
