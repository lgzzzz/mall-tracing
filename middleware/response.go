package middleware

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// ResponseStatusFactory creates a ResponseStatus proto message.
type ResponseStatusFactory func(code int32, msg string) interface{}

// ResponseError returns a middleware that automatically fills ResponseStatus
// in the response message. Business errors are converted to ResponseStatus
// with gRPC OK status. Transport/infrastructure errors are passed through.
func ResponseError(newStatus ResponseStatusFactory) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if reply == nil {
				return reply, err
			}

			if setStatusByReflection(reply, err, newStatus) {
				return reply, nil
			}

			return reply, err
		}
	}
}

func setStatusByReflection(reply interface{}, err error, newStatus ResponseStatusFactory) bool {
	v := reflect.ValueOf(reply)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return false
	}

	statusField := v.Elem().FieldByName("Status")
	if !statusField.IsValid() || !statusField.CanSet() {
		return false
	}

	if err == nil {
		statusField.Set(reflect.Zero(statusField.Type()))
		return true
	}

	if kratosErr := errors.FromError(err); kratosErr != nil {
		status := newStatus(kratosErr.Code, kratosErr.Message)
		statusField.Set(reflect.ValueOf(status))
		return true
	}

	return false
}
