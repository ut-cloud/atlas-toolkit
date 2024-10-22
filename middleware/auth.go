package middleware

import (
	"context"
	"errors"
	"fmt"
	errors2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/ut-cloud/atlas-toolkit/utils"
	"strings"

	middleware2 "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Auth() middleware2.Middleware {
	return func(handler middleware2.Handler) middleware2.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				authHeader := tr.RequestHeader().Get("Authorization")
				if authHeader == "" {
					return nil, errors.New("no Auth")
				}
				// Check if the header starts with "Bearer "
				if !strings.HasPrefix(authHeader, "Bearer ") {
					return nil, fmt.Errorf("authorization header is not a bearer token")
				}
				// Extract the token part by trimming "Bearer "
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if token == "" {
					return nil, fmt.Errorf("token is missing")
				}
				userClaims, err := utils.AnalyseToken(token)
				if err != nil {
					return nil, errors2.New(401, err.Error(), "invalid token")
				}
				if userClaims.Identity == "" {
					return nil, errors.New("no Auth")
				}
				ctx = metadata.NewServerContext(ctx, metadata.New(map[string][]string{
					"username": []string{userClaims.UserName},
					"userId":   []string{userClaims.UserId},
					"identity": []string{userClaims.Identity},
				}))
			}
			return handler(ctx, req)
		}
	}
}
