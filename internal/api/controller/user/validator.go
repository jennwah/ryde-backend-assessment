package user

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterUserValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validator for latitude
		v.RegisterValidation("latitude", func(fl validator.FieldLevel) bool {
			lat := fl.Field().Float()
			return lat >= -90 && lat <= 90
		})

		// Register custom validator for longitude
		v.RegisterValidation("longitude", func(fl validator.FieldLevel) bool {
			lon := fl.Field().Float()
			return lon >= -180 && lon <= 180
		})
	}
}
