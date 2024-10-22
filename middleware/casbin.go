package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/ut-cloud/atlas-toolkit/casbin"
	"github.com/ut-cloud/atlas-toolkit/utils"
	"strings"
)

const (
	ClaimAuthorityId = "authorityId"
)

type SecurityUser struct {
	Path        string
	Method      string
	AuthorityId []string
}

func NewSecurityUser() casbin.SecurityUser {
	return &SecurityUser{}
}

func (su *SecurityUser) ParseFromContext(ctx context.Context) error {
	//if claims, ok := jwt.FromContext(ctx); ok {
	//	su.AuthorityId = claims.(jwtV5.MapClaims)[ClaimAuthorityId].(string)
	//} else {
	//	return errors.New("jwt claim missing")
	//}

	identity := utils.GetLoginIdentity(ctx)
	su.AuthorityId = strings.Split(identity, ",")

	if header, ok := transport.FromServerContext(ctx); ok {
		su.Path = header.Operation()
		su.Method = "*"
	} else {
		return errors.New("jwt claim missing")
	}

	return nil
}

func (su *SecurityUser) GetSubject() []string {
	return su.AuthorityId
}

func (su *SecurityUser) GetObject() string {
	return su.Path
}

func (su *SecurityUser) GetAction() string {
	return su.Method
}

func (su *SecurityUser) CreateAccessJwtToken(secretKey []byte) string {
	claims := jwtV5.NewWithClaims(jwtV5.SigningMethodHS256,
		jwtV5.MapClaims{
			ClaimAuthorityId: su.AuthorityId,
		})

	signedToken, err := claims.SignedString(secretKey)
	if err != nil {
		return ""
	}

	return signedToken
}

func (su *SecurityUser) ParseAccessJwtTokenFromContext(ctx context.Context) error {
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		fmt.Println("ParseAccessJwtTokenFromContext 1")
		return errors.New("no jwt token in context")
	}
	if err := su.ParseAccessJwtToken(claims); err != nil {
		fmt.Println("ParseAccessJwtTokenFromContext 2")
		return err
	}
	return nil
}

func (su *SecurityUser) ParseAccessJwtTokenFromString(token string, secretKey []byte) error {
	parseAuth, err := jwtV5.Parse(token, func(*jwtV5.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	claims, ok := parseAuth.Claims.(jwtV5.MapClaims)
	if !ok {
		return errors.New("no jwt token in context")
	}

	if err := su.ParseAccessJwtToken(claims); err != nil {
		return err
	}

	return nil
}

func (su *SecurityUser) ParseAccessJwtToken(claims jwtV5.Claims) error {
	if claims == nil {
		return errors.New("claims is nil")
	}

	mc, ok := claims.(jwtV5.MapClaims)
	if !ok {
		return errors.New("claims is not map claims")
	}

	strAuthorityId, ok := mc[ClaimAuthorityId]
	if ok {
		su.AuthorityId = strings.Split(strAuthorityId.(string), ",")
	}
	return nil
}
