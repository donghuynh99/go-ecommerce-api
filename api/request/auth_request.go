package request

type LoginRequest struct {
	Email    string `gorm:"email" validate:"required,email"`
	Password string `gorm:"password" validate:"required"`
}

type RegisterRequest struct {
	Email           string `gorm:"email" validate:"required,email,uniqueEmail"`
	Password        string `gorm:"password" validate:"required"`
	ConfirmPassword string `gorm:"confirm_password" json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `gorm:"first_name" json:"first_name" validate:"required"`
	LastName        string `gorm:"last_name" json:"last_name" validate:"required"`
}
