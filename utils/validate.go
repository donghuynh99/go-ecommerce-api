package utils

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/donghuynh99/ecommerce_api/database"
	"github.com/donghuynh99/ecommerce_api/models"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ApiError struct {
	Param   string
	Message string
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return Translation("field_required", nil, nil)
	case "email":
		return Translation("invalid_email", nil, nil)
	case "eqfield":
		return Translation("not_match", nil, nil)
	case "uniqueEmail":
		return Translation("already_existed", nil, nil)
	case "gt":
		return Translation("must_greather_than", map[string]interface{}{
			"Number": fe.Param(),
		}, nil)
	case "arrayIn":
		alloweds := strings.Split(fe.Param(), "&")

		return Translation("should_be", map[string]interface{}{
			"Content": strings.Join(alloweds, " or "),
		}, nil)
	}
	return fe.Error()
}

func validateUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	var db *gorm.DB = database.Database
	var user models.User

	result := db.First(&user, "email = ?", email)

	if result.Error != nil {
		return true
	}

	return false
}

func arrayIn(fl validator.FieldLevel) bool {
	array, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}

	params := fl.Param()
	alloweds := strings.Split(params, "&")

	for _, allowed := range alloweds {
		if strconv.Itoa(array) == allowed {
			return true
		}
	}

	return false
}

func ValidateStruct(data interface{}) []ApiError {
	validate := validator.New()

	validate.RegisterValidation("uniqueEmail", validateUniqueEmail)
	validate.RegisterValidation("arrayIn", arrayIn)

	err := validate.Struct(data)
	if err != nil {
		output := make([]ApiError, len(err.(validator.ValidationErrors)))
		for i, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(data).FieldByName(err.StructField())
			fieldName := field.Tag.Get("gorm")
			output[i] = ApiError{fieldName, MsgForTag(err)}
		}

		return output
	}

	return nil
}
