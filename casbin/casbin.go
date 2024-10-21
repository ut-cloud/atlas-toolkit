package casbin

import (
	"context"
	stdcasbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

type contextKey string

const (
	ModelContextKey        contextKey = "CasbinModel"
	PolicyContextKey       contextKey = "CasbinPolicy"
	EnforcerContextKey     contextKey = "CasbinEnforcer"
	SecurityUserContextKey contextKey = "CasbinSecurityUser"

	reason string = "FORBIDDEN"
)

var (
	ErrSecurityUserCreatorMissing = errors.Forbidden(reason, "SecurityUserCreator is required")
	ErrEnforcerMissing            = errors.Forbidden(reason, "Enforcer is missing")
	ErrSecurityParseFailed        = errors.Forbidden(reason, "Security Info fault")
	ErrUnauthorized               = errors.Forbidden(reason, "Unauthorized Access")
)

type Option func(*options)

type options struct {
	securityUserCreator SecurityUserCreator
	model               model.Model
	policy              persist.Adapter
	enforcer            *stdcasbin.SyncedEnforcer
}

func WithSecurityUserCreator(securityUserCreator SecurityUserCreator) Option {
	return func(o *options) {
		o.securityUserCreator = securityUserCreator
	}
}

func WithCasbinModel(model model.Model) Option {
	return func(o *options) {
		o.model = model
	}
}

func WithCasbinPolicy(policy persist.Adapter) Option {
	return func(o *options) {
		o.policy = policy
	}
}

func Server(opts ...Option) middleware.Middleware {
	o := &options{
		securityUserCreator: nil,
	}
	for _, opt := range opts {
		opt(o)
	}

	o.enforcer, _ = stdcasbin.NewSyncedEnforcer(o.model, o.policy)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			if o.enforcer == nil {
				return nil, ErrEnforcerMissing
			}
			if o.securityUserCreator == nil {
				return nil, ErrSecurityUserCreatorMissing
			}

			securityUser := o.securityUserCreator()
			if err := securityUser.ParseFromContext(ctx); err != nil {
				return nil, ErrSecurityParseFailed
			}

			ctx = context.WithValue(ctx, SecurityUserContextKey, securityUser)

			allowed, err := o.enforcer.Enforce(securityUser.GetSubject(), securityUser.GetObject(), securityUser.GetAction())
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, ErrUnauthorized
			}
			return handler(ctx, req)
		}
	}
}

func Client(opts ...Option) middleware.Middleware {
	o := &options{
		securityUserCreator: nil,
	}
	for _, opt := range opts {
		opt(o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
	}
}