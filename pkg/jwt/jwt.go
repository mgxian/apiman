package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/will835559313/apiman/pkg/setting"
)

var (
	secret string
	expire int
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
	expire, _ = strconv.Atoi(sec.Key("expire").String())
	fmt.Println(expire)
}

func GetToken(name string, admin bool) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		name,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expire)).Unix(),
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

	if err != nil {
		//fmt.Println("in 100")
		fmt.Println(err)
		//fmt.Println("in 101")
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Name, claims.StandardClaims.ExpiresAt)
		return claims, nil
	} else {
		return nil, err
		fmt.Println(err)
	}
	return nil, err
}

func GetClaims(c echo.Context) (*jwtCustomClaims, error) {
	auth := c.Request().Header.Get("Authorization")
	if len(auth) < 8 {
		return nil, errors.New("请添加token请求头")
	}
	token := auth[7:]
	fmt.Println(token)
	claims, err := ParseToken(token)
	if err != nil {
		if strings.Contains(err.Error(), "expire") {
			return nil, errors.New("token已过期")
		}
		return nil, errors.New("token错误")
	}
	return claims, err
}
