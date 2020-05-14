package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	jwt "github.com/dgrijalva/jwt-go"
	"moshopserver/utils"
)

var (
	key = []byte("adfadf!@#2")
	expireTime = 20
	LoginUserId = ""
)

type CustomClaims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

func GetUserID(tokenString string) string {
	token := Parse(tokenString)
	if token == nil {
		return ""
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.UserID
	}
	return ""
}

func Parse(tokenString string) *jwt.Token {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if token.Valid {
		return token
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("The token is expired or not valid.")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	return nil
}

func Create(userId string) string {
	claims := CustomClaims{
		userId, jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(expireTime)).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)

	if err == nil {
		return tokenString
	}
	return ""
}

func Verify(tokenString string) bool {
	token := Parse(tokenString)
	return token != nil
}

func getControllerAndAction(rawValue string) (controller, action string) {
	valueArray := strings.Split(rawValue, "/")
	return valueArray[2], valueArray[2] + "/" + valueArray[3]
}

func FilterFunc(ctx *context.Context) {

	controller, action := getControllerAndAction(ctx.Request.RequestURI)
	token := ctx.Input.Header("x-nideshop-token")

	if action == "auth/loginByWeixin" {
		return
	}

	if token == "" {
		data := utils.GetHTTPRtnJsonData(401, "need relogin")
		ctx.Output.JSON(data, true, false)
		ctx.Redirect(200, "/")
		return
	}
	LoginUserId = GetUserID(token)

	publicControllerList := beego.AppConfig.String("controller::publicController")
	publicActionList := beego.AppConfig.String("action::publicAction")

	if !strings.Contains(publicControllerList, controller) && !strings.Contains(publicActionList, action) {
		if LoginUserId == "" {
			data := utils.GetHTTPRtnJsonData(401, "need relogin")
			ctx.Output.JSON(data, true, false)
			ctx.Redirect(200, "/")
			//http.Redirect(ctx.ResponseWriter, ctx.Request, "/", http.StatusMovedPermanently)
		}
	}
}
