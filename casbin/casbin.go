package casbin

import (
	"context"
	stdcasbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
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
	whiteList           map[string]struct{} // 增加白名单字段
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

// WithWhiteList 新增 WithWhiteList 选项，用于设置白名单
func WithWhiteList(whiteList []string) Option {
	return func(o *options) {
		o.whiteList = make(map[string]struct{})
		for _, path := range whiteList {
			o.whiteList[path] = struct{}{}
		}
	}
}

func Server(opts ...Option) middleware.Middleware {
	o := &options{
		securityUserCreator: nil,
		whiteList:           make(map[string]struct{}), // 初始化白名单
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

			// 获取当前操作路径
			if header, ok := transport.FromServerContext(ctx); ok {
				operation := header.Operation()
				// 如果操作路径在白名单中，跳过 Casbin 检查
				if _, ok := o.whiteList[operation]; ok {
					return handler(ctx, req) // 直接处理请求，无需 Casbin 检查
				}
			}
			// 遍历 AuthorityId 数组，逐个进行 Enforce 检查
			for _, subject := range securityUser.GetSubject() {
				allowed, err := o.enforcer.Enforce(subject, securityUser.GetObject(), securityUser.GetAction())
				if err != nil {
					return nil, err
				}
				if allowed {
					return handler(ctx, req) // 如果某个 subject 被允许，立即返回
				}
			}
			return nil, ErrUnauthorized // 如果没有任何 subject 被允许，返回 Unauthorized 错误
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
