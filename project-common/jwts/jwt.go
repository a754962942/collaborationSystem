package jwts

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtToken struct {
	AccessToken  string
	RefreshToken string
	AccessExp    int64
	RefreshExp   int64
}

func CreateToken(val string, exp time.Duration, secret string, refreshSecret string, refreshExp time.Duration) *JwtToken {
	aExp := time.Now().Add(exp).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   aExp,
	})
	aToken, _ := accessToken.SignedString(secret)
	bExp := time.Now().Add(refreshExp).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   bExp,
	})
	rToken, _ := refreshToken.SignedString(refreshSecret)
	return &JwtToken{
		AccessToken:  aToken,
		RefreshToken: rToken,
		AccessExp:    aExp,
		RefreshExp:   bExp,
	}
}
