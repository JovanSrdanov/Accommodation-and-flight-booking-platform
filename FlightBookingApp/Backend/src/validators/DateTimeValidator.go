package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func NotBeforeCurrentDate(field validator.FieldLevel) bool {
	//Ovako se type castuje u runtime: interface().(assertion)
	timeToCheck, _ := field.Field().Interface().(time.Time)
	return time.Now().Before(timeToCheck)
}
