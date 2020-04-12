package validator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			return name
		})
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func HandleErrors(c *gin.Context, err error) {
	for _, fieldErr := range err.(validator.ValidationErrors) {
		c.JSON(http.StatusBadRequest, FieldError{fieldErr}.String())
	}
}

type FieldError struct {
	Err validator.FieldError
}

func (q FieldError) String() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + q.Err.Field() + "'")
	sb.WriteString(", condition: " + q.Err.ActualTag())

	// Print condition parameters, e.g. oneof=red blue -> { red blue }
	if q.Err.Param() != "" {
		sb.WriteString(" { " + q.Err.Param() + " }")
	}

	if q.Err.Value() != nil && q.Err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.Err.Value()))
	}

	return sb.String()
}
