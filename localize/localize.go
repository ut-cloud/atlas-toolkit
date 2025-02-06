package localize

import (
	"context"
	"gopkg.in/yaml.v3"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type localizerKey struct{}

func I18N() middleware.Middleware {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	//bundle.MustLoadMessageFile("/data/language/active.zh.yaml")
	//bundle.MustLoadMessageFile("/data/language/active.en.yaml")
	bundle.MustLoadMessageFile("i18n/active.zh.yaml")
	bundle.MustLoadMessageFile("i18n/active.en.yaml")

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				accept := tr.RequestHeader().Get("accept-language")
				localizer := i18n.NewLocalizer(bundle, accept)
				ctx = context.WithValue(ctx, localizerKey{}, localizer)
			}
			return handler(ctx, req)
		}
	}
}

func FromContext(ctx context.Context) *i18n.Localizer {
	return ctx.Value(localizerKey{}).(*i18n.Localizer)
}

func GetLocalizeMsg(ctx context.Context, messageId string) string {
	localizer := FromContext(ctx)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageId,
	})
	if err != nil {
		return ""
	}
	return msg
}
