package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/golang-jwt/jwt/v5"
)

// AuthInfo holds JWT claims extracted from the token.
type AuthInfo struct {
	UserID   int64
	Username string
	Role     string
}

type authKey struct{}

// AuthOptions configures the ServerAuth middleware behavior.
type AuthOptions struct {
	SigningMethodCheck bool
	AllowEmptyToken    bool
	ErrUnauthorized    *errors.Error
}

// AuthOption is a functional option for configuring ServerAuth.
type AuthOption func(*AuthOptions)

// WithSigningMethodCheck enables HMAC signing method validation.
func WithSigningMethodCheck() AuthOption {
	return func(o *AuthOptions) {
		o.SigningMethodCheck = true
	}
}

// WithAllowEmptyToken allows requests without a token to pass through.
func WithAllowEmptyToken() AuthOption {
	return func(o *AuthOptions) {
		o.AllowEmptyToken = true
	}
}

// WithUnauthorizedErr sets a custom error for unauthorized requests.
func WithUnauthorizedErr(err *errors.Error) AuthOption {
	return func(o *AuthOptions) {
		o.ErrUnauthorized = err
	}
}

// ServerAuth returns a server-side JWT authentication middleware.
func ServerAuth(secret string, opts ...AuthOption) middleware.Middleware {
	options := &AuthOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.ErrUnauthorized == nil {
		options.ErrUnauthorized = errors.Unauthorized("UNAUTHORIZED", "unauthorized")
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tokenStr := ExtractToken(ctx)
			if tokenStr == "" {
				if options.AllowEmptyToken {
					return handler(ctx, req)
				}
				return nil, options.ErrUnauthorized
			}

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if options.SigningMethodCheck {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, options.ErrUnauthorized
					}
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return nil, options.ErrUnauthorized
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, options.ErrUnauthorized
			}

			info := &AuthInfo{}
			if uid, ok := claims["user_id"].(float64); ok {
				info.UserID = int64(uid)
			}
			if username, ok := claims["username"].(string); ok {
				info.Username = username
			}
			if role, ok := claims["role"].(string); ok {
				info.Role = role
			}

			ctx = context.WithValue(ctx, authKey{}, info)
			return handler(ctx, req)
		}
	}
}

// GetAuthInfo retrieves AuthInfo from context.
func GetAuthInfo(ctx context.Context) (*AuthInfo, bool) {
	info, ok := ctx.Value(authKey{}).(*AuthInfo)
	return info, ok
}

// ExtractToken extracts the JWT token from gRPC metadata.
func ExtractToken(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		auth := md.Get("authorization")
		if auth != "" {
			return strings.TrimPrefix(auth, "Bearer ")
		}
	}
	return ""
}
