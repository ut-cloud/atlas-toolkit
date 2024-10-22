package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserClaims struct {
	Identity string `json:"identity"`
	UserId   string `json:"userId"`
	UserName string `json:"UserName"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("YXRsYXMtYWRtaW4=")

// GenerateToken
// 生成 token
func GenerateToken(identity, userId, username string) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		UserId:         userId,
		UserName:       username,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// GetUUID
// 生成唯一码
func GetUUID() string {
	return uuid.New().String()
}

func GetLoginUserId(ctx context.Context) string {
	var u string
	if md, ok := metadata.FromServerContext(ctx); ok {
		u = md.Get("userId")
	}
	return u
}

func GetLoginIdentity(ctx context.Context) string {
	var u string
	if md, ok := metadata.FromServerContext(ctx); ok {
		u = md.Get("identity")
	}
	return u
}
