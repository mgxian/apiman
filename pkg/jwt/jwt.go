package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/will835559313/apiman/pkg/setting"
)

var (
	secret string
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func JwtInint() {
	// get secret
	sec := setting.Cfg.Section("jwt")
	secret = sec.Key("secret").String()
}

func GetToken(name string, admin bool) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		name,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))

	return t, err
}

func ParseToken(tokenString string) (*jwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Name, claims.StandardClaims.ExpiresAt)
		return claims, nil
	} else {
		return nil, err
		fmt.Println(err)
	}
	return nil, err
}
