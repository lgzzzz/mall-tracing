package middleware

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
)

type testResponse struct {
	Status *TestResponseStatus
	Data   string
}

type TestResponseStatus struct {
	ErrorCode    int32
	ErrorMessage string
}

func TestResponseError_Success(t *testing.T) {
	mw := ResponseError(func(code int32, msg string) interface{} {
		return &TestResponseStatus{ErrorCode: code, ErrorMessage: msg}
	})

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &testResponse{Data: "hello"}, nil
	}

	resp, err := mw(handler)(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := resp.(*testResponse)
	if r.Data != "hello" {
		t.Fatalf("expected data 'hello', got %s", r.Data)
	}
	if r.Status != nil {
		t.Fatal("expected nil Status on success")
	}
}

func TestResponseError_BusinessError(t *testing.T) {
	mw := ResponseError(func(code int32, msg string) interface{} {
		return &TestResponseStatus{ErrorCode: code, ErrorMessage: msg}
	})

	bizErr := errors.NotFound("NOT_FOUND", "item not found")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &testResponse{}, bizErr
	}

	resp, err := mw(handler)(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected nil error (business error converted to status), got %v", err)
	}
	r := resp.(*testResponse)
	if r.Status == nil {
		t.Fatal("expected non-nil Status for business error")
	}
	if r.Status.ErrorMessage != "item not found" {
		t.Fatalf("expected error message 'item not found', got %s", r.Status.ErrorMessage)
	}
}

func TestResponseError_NilReply(t *testing.T) {
	mw := ResponseError(func(code int32, msg string) interface{} {
		return &TestResponseStatus{ErrorCode: code, ErrorMessage: msg}
	})

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.InternalServer("INTERNAL", "internal error")
	}

	resp, err := mw(handler)(context.Background(), nil)
	if resp != nil {
		t.Fatal("expected nil reply")
	}
	if err == nil {
		t.Fatal("expected error for nil reply")
	}
}