package lib
import (
	"fmt"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/astaxie/beego"
)

var (
	key []byte = []byte("-jwt-gree@gree.com")
)

// 产生json web token
func GenToken() string {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		Issuer:    "gree132121",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		beego.Error("this is error")
		return ""
	}
	return ss
}

// 校验token是否有效
func CheckToken(token string) bool {
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return false
	}
	return true
}
