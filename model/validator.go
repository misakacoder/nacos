package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"nacos/router"
	"nacos/util/collection"
	"reflect"
)

func BindHeader[T any](context *gin.Context, object T) T {
	return BindAny(context.ShouldBindHeader, object)
}

func BindUri[T any](context *gin.Context, object T) T {
	return BindAny(context.ShouldBindUri, object)
}

func BindQuery[T any](context *gin.Context, object T) T {
	return BindAny(context.ShouldBindQuery, object)
}

func BindJSON[T any](context *gin.Context, object T) T {
	return BindAny(context.ShouldBindJSON, object)
}

func Bind[T any](context *gin.Context, object T) T {
	return BindAny(context.ShouldBind, object)
}

func BindAny[T any](bindFn func(v any) error, object T) T {
	if err := bindFn(object); err != nil {
		validationError(object, err)
	}
	return object
}

func validationError(object any, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		joiner := collection.NewJoiner(", ", "", "")
		objectElement := reflect.TypeOf(object).Elem()
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			if field, ok := objectElement.FieldByName(fieldError.Field()); ok {
				msg := field.Tag.Get("msg")
				if msg == "" {
					msg = fieldName + " is required"
				}
				joiner.Append(msg)
			}
		}
		if joiner.Size() > 0 {
			panic(router.ParameterMissing.With(joiner.String()))
		}
	}
	panic(err)
}
