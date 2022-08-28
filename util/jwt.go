package util

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
	"trojan/module/constant"
	"trojan/module/vo"
)

// 过期时间默认2小时
const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("4eb01fa4acef754ad4fa94f4467fd343")

type MyClaims struct {
	AccountVo vo.AccountVo `json:"accountVo"`
	jwt.StandardClaims
}

// 生成Token
func GenToken(accountVo vo.AccountVo) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		// 自定义字段
		AccountVo: accountVo,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			// 签发人
			Issuer: "trojan-panel",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, errors.New(constant.IllegalTokenError)
	}
	// 校验Token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New(constant.TokenExpiredError)
}
