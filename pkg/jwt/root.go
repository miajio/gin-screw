package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Token 生成token的结构体
type Token struct {
	Params map[string]string
	jwt.RegisteredClaims
}

/*
EncryptionToken 加密生成token
@param params 以map方式存储的key、value数据
@param secret 秘钥
@param timeout token过期时间
@return string 加密成功的token, error 加密失败的错误信息
*/
func EncryptionToken(params map[string]string, secret string, timeout time.Duration) (string, error) {
	t := Token{
		Params: params,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeout)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),              // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),              // 生效时间
		},
	}
	to := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	return to.SignedString([]byte(secret))
}

// DecryptionToken 解密token
func DecryptionToken(token string, secret string) (map[string]string, error) {
	t, err := jwt.ParseWithClaims(token, &Token{}, Secret(secret))
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		} else {
			return nil, err
		}
	}
	if params, ok := t.Claims.(*Token); ok && t.Valid {
		return params.Params, nil
	}
	return nil, errors.New("couldn't handle this token")
}

// Secret 安全认证
func Secret(secret string) jwt.Keyfunc {
	return func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
